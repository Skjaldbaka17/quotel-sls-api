package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
	"github.com/aws/aws-lambda-go/events"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type RequestHandler struct {
	Db *gorm.DB
}

func (requestHandler *RequestHandler) InitializeDB() structs.ErrorResponse {
	if requestHandler.Db == nil {
		var err error

		dsn := os.Getenv(DATABASE_URL)
		requestHandler.Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			log.Printf("Could not connect to DB, got error: %s", err)
			return structs.ErrorResponse{Message: InternalServerError, StatusCode: http.StatusInternalServerError}

		}

		// defer db.Close()
	}
	return structs.ErrorResponse{}

}

//ValidateRequestBody takes in the request and validates all the input fields, returns an error with reason for validation-failure
//if validation fails.
//TODO: Make validation better! i.e. make it "real"
func (requestHandler *RequestHandler) ValidateRequest(request events.APIGatewayProxyRequest) (interface{}, error) {
	requestBody := structs.Request{}
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		log.Printf("Got err:", err)
		return structs.ErrorResponse{
			Message:    "request body is not structured correctly. Please refer to the /docs page for information on how to structure the request body",
			StatusCode: http.StatusBadRequest}, err
	}

	if err := requestHandler.ValidateRequestApiKey(request); err != nil {
		return err, errors.New("Some Error from validating apikey")
	}

	if requestBody.PageSize < 1 || requestBody.PageSize > maxPageSize {
		requestBody.PageSize = defaultPageSize
	}

	if requestBody.Page < 0 {
		requestBody.Page = 0
	}

	if requestBody.MaxQuotes < 0 || requestBody.MaxQuotes > maxQuotes {
		requestBody.MaxQuotes = maxQuotes
	}

	if requestBody.MaxQuotes <= 0 {
		requestBody.MaxQuotes = defaultMaxQuotes
	}

	const layout = "2006-01-02"
	//Set date into correct format, if supplied, otherwise input today's date in the correct format for all qods
	if len(requestBody.Qods) != 0 {
		for idx, _ := range requestBody.Qods {
			if requestBody.Qods[idx].Date == "" {
				requestBody.Qods[idx].Date = time.Now().UTC().Format(layout)
			} else {
				var parsedDate time.Time
				parsedDate, err := time.Parse(layout, requestBody.Qods[idx].Date)
				if err != nil {
					log.Printf("Got error when decoding: %s", err)
					return structs.ErrorResponse{
						Message:    fmt.Sprintf("the date is not structured correctly, should be in %s format", layout),
						StatusCode: http.StatusBadRequest}, err
				}

				requestBody.Qods[idx].Date = parsedDate.UTC().Format(layout)
			}
		}
	}

	//Set date into correct format, if supplied, otherwise input today's date in the correct format for all qods
	if len(requestBody.Aods) != 0 {
		for idx, _ := range requestBody.Aods {
			if requestBody.Aods[idx].Date == "" {
				requestBody.Aods[idx].Date = time.Now().UTC().Format(layout)
			} else {
				var parsedDate time.Time
				parsedDate, err := time.Parse(layout, requestBody.Aods[idx].Date)
				if err != nil {
					log.Printf("Got error when decoding: %s", err)
					return structs.ErrorResponse{
						Message:    fmt.Sprintf("the date is not structured correctly, should be in %s format", layout),
						StatusCode: http.StatusBadRequest}, err
				}

				requestBody.Aods[idx].Date = parsedDate.UTC().Format(layout)
			}
		}
	}

	if requestBody.Minimum != "" {

		_, err := time.Parse(layout, requestBody.Minimum)
		if err != nil {
			log.Printf("Got error when decoding: %s", err)
			return structs.ErrorResponse{
				Message:    fmt.Sprintf("the minimum date is not structured correctly, should be in %s format", layout),
				StatusCode: http.StatusBadRequest}, err
		}
	}

	if requestBody.Maximum != "" {

		parseDate, err := time.Parse(layout, requestBody.Maximum)
		if err != nil {
			log.Printf("Got error when decoding: %s", err)
			return structs.ErrorResponse{
				Message:    fmt.Sprintf("the maximum date is not structured correctly, should be in %s format", layout),
				StatusCode: http.StatusBadRequest}, err
		}
		requestBody.Minimum = parseDate.Format("01-02-2006")
	}

	return requestBody, nil
}

// ValidateRequestApiKey checks if the ApiKey supplied exists and wether the user has finished his allowed request in the past
// hour. Also adds to the requestHistory... Maybe move that to the end of a request?
func (requestHandler *RequestHandler) ValidateRequestApiKey(request events.APIGatewayProxyRequest) interface{} {
	requestBody := structs.Request{}
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		return structs.ErrorResponse{
			Message:    "request body is not structured correctly. Please refer to the /docs page for information on how to structure the request body",
			StatusCode: http.StatusBadRequest}
	}

	if requestBody.ApiKey == "" {
		log.Printf("no ApiKey given when accessing resource")
		return structs.ErrorResponse{
			Message:    fmt.Sprintf("you need to supply an apiKey to access this resource. Create a user and get a free-tier apiKey here: %s", os.Getenv("WEBSITE_URL")),
			StatusCode: http.StatusForbidden}
	}

	var user structs.UserDBModel

	err = requestHandler.Db.Table("users").Where("api_key = ?", requestBody.ApiKey).First(&user).Error
	// Err==nil if user with given api_key does not exist or internal server error
	if err != nil {
		m1 := regexp.MustCompile(`record not found`)
		if m1.Match([]byte(err.Error())) {
			log.Printf("the api-key that the requester supplied does not exist")
			return structs.ErrorResponse{
				Message:    fmt.Sprintf("you need to supply an apiKey to access this resource. Create a user and get a free-tier apiKey here: %s", os.Getenv("WEBSITE_URL")),
				StatusCode: http.StatusForbidden}
		}
		log.Printf("error when searching for user with the given api key (api key validation): %s", err)
		return structs.ErrorResponse{
			Message:    InternalServerError,
			StatusCode: http.StatusInternalServerError}
	}

	//Check if requests from this api-key the past hour are less than allowed for the users-tier (i.e. if this next request is
	// allowed then save the request to request-history)
	type countStruct struct {
		Count int `json:"count"`
	}
	var count countStruct
	if err := requestHandler.Db.Table("requesthistory").Select("count(*)").
		Where("created_at >= (NOW() - INTERVAL '1 hour')").
		Where("user_id = ?", user.Id).
		First(&count).Error; err != nil {
		log.Printf("error when counting request history: %s", err)
		return structs.ErrorResponse{
			Message:    InternalServerError,
			StatusCode: http.StatusInternalServerError}
	}

	if float64(count.Count) >= REQUESTS_PER_HOUR[user.Tier] {
		return structs.ErrorResponse{
			Message:    fmt.Sprintf("you have used all the requests per hour that your tier %s allows for, i.e. %f request per hour. See %s for more info and pricing plans to upgrade your tier if necessary", user.Tier, REQUESTS_PER_HOUR[user.Tier], os.Getenv("WEBSITE_URL")),
			StatusCode: http.StatusUnauthorized}
	}

	//TODO: Put the following in its own golang function and run as a separate process!
	requestAsString, _ := json.Marshal(request)
	requestEvent := structs.RequestEvent{
		UserId:      user.Id,
		RequestBody: request.Body,
		Route:       request.Path,
		ApiKey:      user.ApiKey,
		Request:     string(requestAsString),
	}
	result := requestHandler.Db.Table("requesthistory").Create(&requestEvent)
	if result.Error != nil {
		log.Printf("error when inserting into requestHistory: %s", result.Error)
		return structs.ErrorResponse{
			Message:    InternalServerError,
			StatusCode: http.StatusInternalServerError}
	}

	return nil
}
