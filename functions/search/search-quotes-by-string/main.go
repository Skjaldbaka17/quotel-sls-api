package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gorm.io/gorm/clause"
)

type RequestHandler struct {
	utils.RequestHandler
}

var theReqHandler = RequestHandler{}

// swagger:route POST /search/quotes SEARCH SearchQuotesByString
// Quotes search. Searching quotes by a given search string
// responses:
//  200: topicViewsResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// SearchQuotesByString handles POST requests to search for quotes by a search-string
func (requestHandler *RequestHandler) handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//Initialize DB if requestHandler.Db = nil
	if errResponse := requestHandler.InitializeDB(); errResponse != (structs.ErrorResponse{}) {
		return events.APIGatewayProxyResponse{
			Body:       errResponse.Message,
			StatusCode: errResponse.StatusCode,
		}, nil
	}

	requestBody, errResponse := requestHandler.ValidateRequest(request)

	if errResponse != (structs.ErrorResponse{}) {
		return events.APIGatewayProxyResponse{
			Body:       errResponse.Message,
			StatusCode: errResponse.StatusCode,
		}, nil
	}

	var topicResults []structs.TopicViewDBModel
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	dbPointer := requestHandler.GetBasePointer(requestBody)
	dbPointer = dbPointer.Where("( quote_tsv @@ plainq OR quote_tsv @@ phraseq OR quote_tsv @@ generalq)")

	if requestBody.AuthorId > 0 {
		dbPointer = dbPointer.Where("author_id = ?", requestBody.AuthorId)
	}

	//Order by quote_id to have definitive order (when for examplke some quotes rank the same for plain, phrase and general)
	dbPointer = dbPointer.
		Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "plainrank DESC, phraserank DESC, generalrank DESC, quote_id DESC", Vars: []interface{}{}, WithoutParentheses: true},
		})

	//Particular language search
	dbPointer = utils.QuoteLanguageSQL(requestBody.Language, dbPointer)
	//** ---------- Paramatere configuratino for DB query ends ---------- **//
	err := utils.Pagination(requestBody, dbPointer).
		Find(&topicResults).Error

	if err != nil {
		log.Printf("Got error when querying DB in SearchQuotesByString: %s", err)
		return events.APIGatewayProxyResponse{
			Body:       utils.InternalServerError,
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	//Update popularity in background! TODO: add as its own lambda function
	// go handlers.TopicViewAppearInSearchCountIncrement(topicResults)
	apiResults := structs.ConvertToTopicViewsAPIModel(topicResults)
	out, _ := json.Marshal(apiResults)
	return events.APIGatewayProxyResponse{
		Body:       string(out),
		StatusCode: http.StatusOK,
	}, nil

}

func main() {
	lambda.Start(theReqHandler.handler)
}
