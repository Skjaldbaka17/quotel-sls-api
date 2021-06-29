package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RequestHandler struct {
	utils.RequestHandler
}

var theReqHandler = RequestHandler{}

// swagger:route POST /authors/random AUTHORS GetRandomAuthor
// Get a random Author, and some of his quotes, according to the given parameters
// responses:
//	200: searchViewsResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// GetRandomAuthor handles POST requests for getting a random author
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

	var result []structs.SearchViewDBModel
	var author structs.AuthorDBModel
	//** ---------- Paramatere configuratino for DB query begins ---------- **//

	//Get Random author
	dbPointer := requestHandler.Db.Table("authors").Order("random()")

	//author from a particular language
	dbPointer = utils.AuthorLanguageSQL(requestBody.Language, dbPointer)
	//** ---------- Paramatere configuratino for DB query ends ---------- **//

	err := dbPointer.First(&author).Error

	if err != nil {
		log.Printf("Got error when querying DB, first one, in GetRandomAuthor: %s", err)
		return events.APIGatewayProxyResponse{
			Body:       utils.InternalServerError,
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	dbPointer = requestHandler.Db.Table("searchview").Where("author_id = ?", author.Id)

	//An icelandic quote from the particular/random author
	dbPointer = utils.QuoteLanguageSQL(requestBody.Language, dbPointer)

	err = dbPointer.Limit(requestBody.MaxQuotes).Find(&result).Error

	if err != nil {
		log.Printf("Got error when querying DB, second one, in GetAuthors: %s", err)
		return events.APIGatewayProxyResponse{
			Body:       utils.InternalServerError,
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	searchViewAPI := structs.ConvertToSearchViewsAPIModel(result)
	out, _ := json.Marshal(searchViewAPI)
	return events.APIGatewayProxyResponse{
		Body:       string(out),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(theReqHandler.handler)
}
