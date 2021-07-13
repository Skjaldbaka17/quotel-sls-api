package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RequestHandler struct {
	utils.RequestHandler
}

var theReqHandler = RequestHandler{}

// swagger:route POST /authors/aod/history AUTHORS GetAODHistory
// Gets the history for the authors of the day
// responses:
//	200: aodHistoryResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

//GetAODHistory gets Aod history starting from some point
func (requestHandler *RequestHandler) handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	//Initialize DB if requestHandler.Db = nil
	if errResponse := requestHandler.InitializeDB(); errResponse != (structs.ErrorResponse{}) {
		return events.APIGatewayProxyResponse{
			Body:       errResponse.ToString(),
			StatusCode: errResponse.StatusCode,
		}, nil
	}

	requestBody, errResponse := requestHandler.ValidateRequest(request)

	if errResponse != (structs.ErrorResponse{}) {
		return events.APIGatewayProxyResponse{
			Body:       errResponse.ToString(),
			StatusCode: errResponse.StatusCode,
		}, nil
	}

	if requestBody.Language == "" {
		requestBody.Language = "English"
	}
	var authors []structs.AodDBModel
	var err error
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	dbPointer := utils.AodLanguageSQL(requestBody.Language, requestHandler.Db)

	if requestBody.Minimum == "" {
		requestBody.Minimum = "1900-12-21"
	}
	now := time.Now()
	minDate, err := time.Parse("2006-01-02", requestBody.Minimum)

	if err != nil {
		log.Printf("Got error when parsing mindate in GetAODHistory: %s", err)
		errResponse := structs.ErrorResponse{
			Message: "Please supply date in '2020-12-21' format",
		}
		return events.APIGatewayProxyResponse{
			Body:       errResponse.ToString(),
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	if !now.After(minDate) {
		log.Printf("Got error when comparing mindate to today in GetAodHistory: %s", err)
		errResponse := structs.ErrorResponse{
			Message: "Please send a minimum date that is before today",
		}
		return events.APIGatewayProxyResponse{
			Body:       errResponse.ToString(),
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	//Not maximum because then possibility of endless cycle with the if statement below!
	if requestBody.Minimum != "" {
		dbPointer = dbPointer.Where("date >= ?", requestBody.Minimum)
	}
	dbPointer = dbPointer.Where("date <= current_date").Order("date DESC")
	//** ---------- Paramatere configuratino for DB query ends ---------- **//
	err = dbPointer.Find(&authors).Error

	if err != nil {
		log.Printf("Got error when querying DB in GetAodHistory: %s", err)
		errResponse := structs.ErrorResponse{
			Message: utils.InternalServerError,
		}
		return events.APIGatewayProxyResponse{
			Body:       errResponse.ToString(),
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	reg := regexp.MustCompile(time.Now().Format("2006-01-02"))

	if (len(authors) == 0 || !reg.Match([]byte(authors[0].Date))) && !requestBody.StopRecursion {
		err = requestHandler.SetNewRandomAOD(requestBody.Language)
		if err != nil {
			log.Printf("Got error when setting newRandomAOD in getAODHistory: %s", err)
			errResponse := structs.ErrorResponse{
				Message: utils.InternalServerError,
			}
			return events.APIGatewayProxyResponse{
				Body:       errResponse.ToString(),
				StatusCode: http.StatusInternalServerError,
			}, nil
		}
		requestBody.StopRecursion = true
		bod, _ := json.Marshal(requestBody)
		request.Body = string(bod)
		return requestHandler.handler(request)
	}

	authorsAPI := structs.ConvertToAodAPIModel(authors)
	out, _ := json.Marshal(authorsAPI)
	return events.APIGatewayProxyResponse{
		Body:       string(out),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(theReqHandler.handler)
}
