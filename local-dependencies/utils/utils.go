package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
	"github.com/aws/aws-lambda-go/events"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type RequestHandler struct {
	Db *gorm.DB
}

func (requestHandler *RequestHandler) InitializeDB() structs.ErrorResponse {
	if os.Getenv("DATABASE_URL") == "" {
		godotenv.Load("../../../.env")
	}

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

func (requestHandler *RequestHandler) GetRandomQuoteFromDb(requestBody *structs.Request) (structs.QuoteDBModel, error) {
	var dbPointer *gorm.DB
	var topicResult []structs.QuoteDBModel

	var shouldDoQuick = true

	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	m1 := regexp.MustCompile(` `)
	phrasesearch := m1.ReplaceAllString(requestBody.SearchString, " <-> ")

	//Random quote from a particular topic
	if len(requestBody.TopicIds) > 0 {
		dbPointer = requestHandler.Db.Table("topicsview, plainto_tsquery(?) as plainq, to_tsquery(?) as phraseq", requestBody.SearchString, phrasesearch).Where("topic_id in ?", requestBody.TopicIds)
		shouldDoQuick = false
	} else {
		dbPointer = requestHandler.Db.Table("quotes, plainto_tsquery(?) as plainq, to_tsquery(?) as phraseq", requestBody.SearchString, phrasesearch)
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
		dbPointer = dbPointer.Raw("select * from quotes tablesample system(0.1)")
	}

	//** ---------- Paramater configuratino for DB query ends ---------- **//
	err := dbPointer.Limit(100).Find(&topicResult).Error

	if err != nil {
		return structs.QuoteDBModel{}, err
	}
	if len(topicResult) == 0 {
		return structs.QuoteDBModel{}, nil
	}
	toReturn := topicResult[0]
	//Sometimes if request is for an english quote we are returned an icelandic quote (because of tablesample) therefore go through the 100 quotes returned and return the first
	//non-icelandic one you find
	if toReturn.IsIcelandic && strings.ToLower(requestBody.Language) != "icelandic" {
		for _, quote := range topicResult {
			if !quote.IsIcelandic {
				toReturn = quote
				break
			}
		}
	}
	return toReturn, nil
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
	dbPointer := requestHandler.Db.Table("qods")
	switch strings.ToLower(language) {
	case "icelandic":
		return dbPointer.Where("is_icelandic")
	default:
		return dbPointer.Not("is_icelandic")
	}
}

//authorLanguageSQL adds to the sql query for the authors db a condition of whether the authors to be fetched have quotes in a particular language
func AuthorLanguageSQL(language string, dbPointer *gorm.DB) *gorm.DB {
	if language != "" {
		switch strings.ToLower(language) {
		case "english":
			dbPointer = dbPointer.Where("nr_of_icelandic_quotes = 0")
		case "icelandic":
			dbPointer = dbPointer.Where("nr_of_icelandic_quotes > 0")
		}
	}
	return dbPointer
}

//aodLanguageSQL adds to the sql query for the authors db a condition of whether the authors to be fetched have quotes in a particular language
func AodLanguageSQL(language string, dbPointer *gorm.DB) *gorm.DB {
	dbPointer = dbPointer.Table("aods")
	switch strings.ToLower(language) {
	case "icelandic":
		return dbPointer.Where("is_icelandic")
	default:
		return dbPointer.Not("is_icelandic")
	}
}

// CheckForSpellingErrorsInSearchString takes the searchstring and partitions it into its separate words (max 20)
// Then it runs a similarity search on the unique_lexeme table (where all the words from quotes.quote and authors.name are stored)
// this is to check if the user made some spelling errors. Then the searchstring is put into the firstSearch and the
// same search as before run again.
func (requestHandler *RequestHandler) CheckForSpellingErrorsInSearchString(searchString string, table string) string {
	newSearchString := ""
	for idx, word := range strings.Fields(searchString) {
		if idx >= 20 {
			break
		}
		var theWord []string
		err := requestHandler.Db.Table(table).
			Select("word").
			Where("similarity(word, ?) > 0.4", word).
			Clauses(clause.OrderBy{
				Expression: clause.Expr{SQL: "similarity(word,?) desc, length(word)", Vars: []interface{}{word}, WithoutParentheses: true},
			}).
			Limit(1).
			Find(&theWord).Error

		if err != nil {
			log.Printf("Got error when querying DB in SearchByString: %s", err)
		}
		if len(theWord) > 0 && theWord[0] != "" {
			newSearchString = strings.Join([]string{newSearchString, theWord[0]}, " ")
		}
	}
	return newSearchString
}
