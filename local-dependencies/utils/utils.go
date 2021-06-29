package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
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
	log.Println("HEREBRUV:" + os.Getenv(DATABASE_URL))
	if requestHandler.Db == nil {
		var err error

		dsn := os.Getenv(DATABASE_URL)
		requestHandler.Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Printf("Could not connect to DB, got error: %s", err)
			return structs.ErrorResponse{Message: InternalServerError, StatusCode: http.StatusInternalServerError}

		}

		// defer db.Close()
	}
	return structs.ErrorResponse{}
}

func (requestHandler *RequestHandler) GetRandomQuoteFromDb(requestBody *structs.Request) (structs.TopicViewAPIModel, error) {
	const NR_OF_QUOTES = 639028
	const NR_OF_ENGLISH_QUOTES = 634841
	var dbPointer *gorm.DB
	var topicResult structs.TopicViewDBModel

	var shouldDoQuick = true

	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	m1 := regexp.MustCompile(` `)
	phrasesearch := m1.ReplaceAllString(requestBody.SearchString, " <-> ")

	//Random quote from a particular topic
	if requestBody.TopicId > 0 {
		dbPointer = requestHandler.Db.Table("topicsview, plainto_tsquery(?) as plainq, to_tsquery(?) as phraseq", requestBody.SearchString, phrasesearch).Where("topic_id = ?", requestBody.TopicId)
		shouldDoQuick = false
	} else {
		dbPointer = requestHandler.Db.Table("searchview, plainto_tsquery(?) as plainq, to_tsquery(?) as phraseq", requestBody.SearchString, phrasesearch)
	}

	//Random quote from a particular author
	if requestBody.AuthorId > 0 {
		dbPointer = dbPointer.Where("author_id = ?", requestBody.AuthorId)
		shouldDoQuick = false
	}

	//Random quote from a particular language
	dbPointer = QuoteLanguageSQL(requestBody.Language, dbPointer)

	if strings.ToLower(requestBody.Language) == "icelandic" {
		shouldDoQuick = false
	}

	if requestBody.SearchString != "" {
		dbPointer = dbPointer.Where("( quote_tsv @@ plainq OR quote_tsv @@ phraseq)")
		shouldDoQuick = false
	}

	//Order by used to get random quote if there are "few" rows returned
	if !shouldDoQuick {
		dbPointer = dbPointer.Order("random()") //Randomized, O( n*log(n) )
	} else {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		if strings.ToLower(requestBody.Language) == "english" {
			dbPointer = dbPointer.Offset(r.Intn(NR_OF_ENGLISH_QUOTES))
		} else {
			dbPointer = dbPointer.Offset(r.Intn(NR_OF_QUOTES))
		}

	}

	//** ---------- Paramater configuratino for DB query ends ---------- **//
	err := dbPointer.Limit(1).Find(&topicResult).Error
	if err != nil {
		return structs.TopicViewAPIModel{}, err
	}
	return topicResult.ConvertToAPIModel(), nil
}

//setQOD inserts a new row into qod/qodice table
func (requestHandler *RequestHandler) SetQOD(language string, date string, quoteId int) error {
	switch strings.ToLower(language) {
	case "icelandic":
		return requestHandler.Db.Exec("insert into qodice (quote_id, date) values((select id from quotes where id = ? and is_icelandic), ?) on conflict (date) do update set quote_id = ?", quoteId, date, quoteId).Error
	default:
		return requestHandler.Db.Exec("insert into qod (quote_id, date) values((select id from quotes where id = ? and not is_icelandic), ?) on conflict (date) do update set quote_id = ?", quoteId, date, quoteId).Error
	}
}

// Check whether user has GOD-tier permissions
func (requestHandler *RequestHandler) AuthorizeGODApiKey(request events.APIGatewayProxyRequest) structs.ErrorResponse {
	requestBody := structs.Request{}
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		log.Printf("Got err: %s", err)
		return structs.ErrorResponse{Message: InternalServerError, StatusCode: http.StatusInternalServerError}
	}

	var user structs.UserDBModel
	if err := requestHandler.Db.Table("users").Where("api_key = ?", requestBody.ApiKey).First(&user).Error; err != nil {
		log.Printf("error when searching for user with the given api key in AuthorIzeGOD (api key validation): %s", err)
		return structs.ErrorResponse{
			Message:    "You need special privileges to access this route.",
			StatusCode: http.StatusUnauthorized,
		}
	}

	if user.Tier != TIERS[len(TIERS)-1] {
		return structs.ErrorResponse{
			Message:    "you do not have the authorization to perform this action. Is your name Bassi Maraj? This is not meant for you... Sorry for the inconvenience",
			StatusCode: http.StatusUnauthorized,
		}
	}

	return structs.ErrorResponse{}
}

//ValidateRequestBody takes in the request and validates all the input fields, returns an error with reason for validation-failure
//if validation fails.
//TODO: Make validation better! i.e. make it "real"
func (requestHandler *RequestHandler) ValidateRequest(request events.APIGatewayProxyRequest) (structs.Request, structs.ErrorResponse) {
	requestBody := structs.Request{}
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		log.Printf("Got err: %s", err)
		return structs.Request{}, structs.ErrorResponse{
			Message:    "request body is not structured correctly. Please refer to the /docs page for information on how to structure the request body",
			StatusCode: http.StatusBadRequest}
	}

	if err := requestHandler.ValidateRequestApiKey(request); err != (structs.ErrorResponse{}) {
		return structs.Request{}, err
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
					return structs.Request{}, structs.ErrorResponse{
						Message:    fmt.Sprintf("the date is not structured correctly, should be in %s format", layout),
						StatusCode: http.StatusBadRequest}
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
					return structs.Request{}, structs.ErrorResponse{
						Message:    fmt.Sprintf("the date is not structured correctly, should be in %s format", layout),
						StatusCode: http.StatusBadRequest}
				}

				requestBody.Aods[idx].Date = parsedDate.UTC().Format(layout)
			}
		}
	}

	if requestBody.Minimum != "" {

		_, err := time.Parse(layout, requestBody.Minimum)
		if err != nil {
			log.Printf("Got error when decoding: %s", err)
			return structs.Request{}, structs.ErrorResponse{
				Message:    fmt.Sprintf("the minimum date is not structured correctly, should be in %s format", layout),
				StatusCode: http.StatusBadRequest}
		}
	}

	if requestBody.Maximum != "" {

		parseDate, err := time.Parse(layout, requestBody.Maximum)
		if err != nil {
			log.Printf("Got error when decoding: %s", err)
			return structs.Request{}, structs.ErrorResponse{
				Message:    fmt.Sprintf("the maximum date is not structured correctly, should be in %s format", layout),
				StatusCode: http.StatusBadRequest}
		}
		requestBody.Minimum = parseDate.Format("01-02-2006")
	}

	return requestBody, structs.ErrorResponse{}
}

// ValidateRequestApiKey checks if the ApiKey supplied exists and wether the user has finished his allowed request in the past
// hour. Also adds to the requestHistory... Maybe move that to the end of a request?
func (requestHandler *RequestHandler) ValidateRequestApiKey(request events.APIGatewayProxyRequest) structs.ErrorResponse {
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

	return structs.ErrorResponse{}
}

func Pagination(requestBody structs.Request, dbPointer *gorm.DB) *gorm.DB {
	return dbPointer.Limit(requestBody.PageSize).
		Offset(requestBody.Page * requestBody.PageSize)
}

//quoteLanguageSQL adds to the sql query for the quotes db a condition of whether the quotes to be fetched are in a particular language
func QuoteLanguageSQL(language string, dbPointer *gorm.DB) *gorm.DB {
	if language != "" {
		switch strings.ToLower(language) {
		case "english":
			dbPointer = dbPointer.Not("is_icelandic")
		case "icelandic":
			dbPointer = dbPointer.Where("is_icelandic")
		}
	}
	return dbPointer
}

//setMaxMinNumber sets the condition for which authors to return
func SetMaxMinNumber(orderConfig structs.OrderConfig, column string, orderDirection string, dbPointer *gorm.DB) *gorm.DB {
	if nr, err := strconv.Atoi(orderConfig.Maximum); err == nil {
		dbPointer = dbPointer.Where(column+" <= ?", nr)
	}
	if nr, err := strconv.Atoi(orderConfig.Minimum); err == nil {
		dbPointer = dbPointer.Where(column+" >= ?", nr)
	}
	return dbPointer.Order(column + " " + orderDirection)
}

//qodLanguageSQL adds to the sql query for the quotes db a condition of whether the quotes to be fetched are quotes in a particular language
func (requestHandler *RequestHandler) QodLanguageSQL(language string) *gorm.DB {
	switch strings.ToLower(language) {
	case "icelandic":
		return requestHandler.Db.Table("qodiceview")
	default:
		return requestHandler.Db.Table("qodview")
	}
}

//SetNewRandomQOD sets a random quote as the qod for today (if language=icelandic is supplied then it adds the random qod to the icelandic qod table)
func (requestHandler *RequestHandler) SetNewRandomQOD(language string) error {
	var quoteItem structs.QuoteDBModel
	var dbPointer *gorm.DB
	dbPointer = requestHandler.Db.Table("quotes")
	dbPointer = QuoteLanguageSQL(language, dbPointer)
	if strings.ToLower(language) != "icelandic" {
		dbPointer = dbPointer.Where("Random() < 0.005")
	}

	err := dbPointer.Order("random()").Limit(1).Scan(&quoteItem).Error
	if err != nil {
		return err
	}

	return requestHandler.setQOD(language, time.Now().Format("2006-01-02"), quoteItem.Id)
}

//setQOD inserts a new row into qod/qodice table
func (requestHandler *RequestHandler) setQOD(language string, date string, quoteId int) error {
	switch strings.ToLower(language) {
	case "icelandic":
		return requestHandler.Db.Exec("insert into qodice (quote_id, date) values((select id from quotes where id = ? and is_icelandic), ?) on conflict (date) do update set quote_id = ?", quoteId, date, quoteId).Error
	default:
		return requestHandler.Db.Exec("insert into qod (quote_id, date) values((select id from quotes where id = ? and not is_icelandic), ?) on conflict (date) do update set quote_id = ?", quoteId, date, quoteId).Error
	}
}

//authorLanguageSQL adds to the sql query for the authors db a condition of whether the authors to be fetched have quotes in a particular language
func AuthorLanguageSQL(language string, dbPointer *gorm.DB) *gorm.DB {
	if language != "" {
		switch strings.ToLower(language) {
		case "english":
			dbPointer = dbPointer.Not("has_icelandic_quotes")
		case "icelandic":
			dbPointer = dbPointer.Where("has_icelandic_quotes")
		}
	}
	return dbPointer
}

//aodLanguageSQL adds to the sql query for the authors db a condition of whether the authors to be fetched have quotes in a particular language
func AodLanguageSQL(language string, dbPointer *gorm.DB) *gorm.DB {
	switch strings.ToLower(language) {
	case "icelandic":
		return dbPointer.Table("aodiceview")
	default:
		return dbPointer.Table("aodview")
	}
}

//setAOD inserts a new row into the aod/aodice table
func (requestHandler *RequestHandler) SetAOD(language string, date string, authorId int) error {
	switch strings.ToLower(language) {
	case "icelandic":
		return requestHandler.Db.Exec("insert into aodice (author_id, date) values((select id from authors where id = ? and has_icelandic_quotes), ?) on conflict (date) do update set author_id = ?", authorId, date, authorId).Error
	default:
		return requestHandler.Db.Exec("insert into aod (author_id, date) values((select id from authors where id = ? and not has_icelandic_quotes), ?) on conflict (date) do update set author_id = ?", authorId, date, authorId).Error
	}
}

//SetNewRandomQOD sets a random quote as the qod for today (if language=icelandic is supplied then it adds the random qod to the icelandic qod table)
func (requestHandler *RequestHandler) SetNewRandomAOD(language string) error {
	var authorItem structs.AuthorDBModel

	if language == "" {
		language = "english"
	}
	dbPointer := requestHandler.Db.Table("authors")
	dbPointer = AuthorLanguageSQL(language, dbPointer)

	err := dbPointer.Order("random()").Limit(1).Scan(&authorItem).Error
	if err != nil {
		return err
	}

	return requestHandler.SetAOD(language, time.Now().Format("2006-01-02"), authorItem.Id)
}

//getBasePointer returns a base DB pointer for a table for a thorough full text search
func (requestHandler *RequestHandler) GetBasePointer(requestBody structs.Request) *gorm.DB {
	table := "searchview"
	//TODO: Validate that this topicId exists
	if requestBody.TopicId > 0 {
		table = "topicsview"
	}
	m1 := regexp.MustCompile(` `)
	phrasesearch := m1.ReplaceAllString(requestBody.SearchString, " <-> ")
	generalsearch := m1.ReplaceAllString(requestBody.SearchString, " | ")
	dbPointer := requestHandler.Db.Table(table+", plainto_tsquery(?) as plainq, to_tsquery(?) as phraseq,to_tsquery(?) as generalq ",
		requestBody.SearchString, phrasesearch, generalsearch).Select("*, ts_rank(quote_tsv, plainq) as plainrank, ts_rank(quote_tsv, phraseq) as phraserank, ts_rank(quote_tsv, generalq) as generalrank")

	if requestBody.TopicId > 0 {
		dbPointer = dbPointer.Where("topic_id = ?", requestBody.TopicId)
	}
	return dbPointer
}

//ValidateUserRequestBody takes in the request and validates all the input fields, returns an error with reason for validation-failure
//if validation fails.
//TODO: Make validation better! i.e. make it "real"

func GetUserRequestBody(request events.APIGatewayProxyRequest) (structs.UserApiModel, structs.ErrorResponse) {
	//Save the state back into the body for later use (Especially useful for getting the AOD/QOD because if the AOD has not been set a random AOD is set and the function called again)
	requestBody := structs.UserApiModel{}
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		log.Printf("Got err: %s", err)
		return structs.UserApiModel{}, structs.ErrorResponse{
			Message:    "request body is not structured correctly. Please refer to the /docs page for information on how to structure the request body",
			StatusCode: http.StatusBadRequest}
	}

	return requestBody, structs.ErrorResponse{}
}

func ValidateUserInformation(requestBody *structs.UserApiModel) structs.ErrorResponse {
	//TODO: Add email validation
	if requestBody.Email == "" {
		return structs.ErrorResponse{
			Message:    "email should not be empty",
			StatusCode: http.StatusBadRequest}
	}

	if requestBody.Name == "" {
		return structs.ErrorResponse{
			Message:    "name should not be empty",
			StatusCode: http.StatusBadRequest}
	}

	if requestBody.Password == "" {
		return structs.ErrorResponse{
			Message:    "password should not be empty",
			StatusCode: http.StatusBadRequest}
	}

	if len(requestBody.Password) < 8 {
		return structs.ErrorResponse{
			Message:    "password should be at least 8 characters long",
			StatusCode: http.StatusBadRequest}
	}

	if requestBody.PasswordConfirmation == "" {
		return structs.ErrorResponse{
			Message:    "password confirmation should not be empty",
			StatusCode: http.StatusBadRequest}
	}

	if requestBody.PasswordConfirmation != requestBody.Password {
		return structs.ErrorResponse{
			Message:    "passwords do not match",
			StatusCode: http.StatusBadRequest}
	}
	return structs.ErrorResponse{}
}
