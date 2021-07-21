package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gorm.io/gorm"
)

type RequestHandler struct {
	utils.RequestHandler
}

var theReqHandler = RequestHandler{}

func search(requestBody *structs.Request, dbPointer *gorm.DB) *gorm.DB {
	table := "quotes"
	//TODO: Validate that this topicId exists
	if len(requestBody.TopicIds) > 0 {
		table = "topicsview"
	}
	dbPointer = dbPointer.Table(table+", plainto_tsquery('english', ?) as plainq ",
		requestBody.SearchString).Select("*, ts_rank(tsv, plainq) as plainrank")

	if requestBody.TopicId > 0 {
		dbPointer = dbPointer.Where("topic_id in ?", requestBody.TopicId)
	}

	//Order by authorid to have definitive order (when for examplke some quotes rank the same for plain, phrase, general and similarity)
	dbPointer = dbPointer.Where("( quote_tsv @@ plainq )").Order("plainrank desc, author_id desc")

	//Particular language search
	dbPointer = utils.QuoteLanguageSQL(requestBody.Language, dbPointer)
	return dbPointer
}

// swagger:route POST /search/quotes search SearchQuotesByString
//
// Search quotes
//
// Use this route to search for quotes by a general full test search that searches for words and phrases in the quotes.
//
// responses:
//  200: quotesApiResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// SearchQuotesByString handles POST requests to search for quotes by a search-string
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

	var topicResults []structs.QuoteDBModel

	var dbPointer *gorm.DB
	for i := 0; i < 2; i++ {
		if i == 1 {
			requestBody.SearchString = requestHandler.CheckForSpellingErrorsInSearchString(requestBody.SearchString, "unique_lexeme_quotes")
		}

		dbPointer = search(&requestBody, requestHandler.Db)
		err := utils.Pagination(requestBody, dbPointer).
			Find(&topicResults).Error

		if err != nil {
			log.Printf("Got error when querying DB in SearchQuotesByString: %s", err)
			errResponse := structs.ErrorResponse{
				Message: utils.InternalServerError,
			}
			return events.APIGatewayProxyResponse{
				Body:       errResponse.ToString(),
				StatusCode: http.StatusInternalServerError,
			}, nil
		}

		if len(topicResults) > 0 {
			break
		}
	}

	//Update popularity in background! TODO: add as its own lambda function
	go requestHandler.TopicViewAppearInSearchCountIncrement(topicResults)
	apiResults := structs.ConvertToQuotesAPIModel(topicResults)
	out, _ := json.Marshal(apiResults)
	return events.APIGatewayProxyResponse{
		Body:       string(out),
		StatusCode: http.StatusOK,
	}, nil

}

func main() {
	lambda.Start(theReqHandler.handler)
}
