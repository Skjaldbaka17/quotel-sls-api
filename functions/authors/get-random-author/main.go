package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RequestHandler struct {
	utils.RequestHandler
}

var theReqHandler = RequestHandler{}

// swagger:route POST /authors/random authors GetRandomAuthor
//
// Get a random Author
//
// Use this route to get a random author, and some of his quotes.
//
// responses:
//	200: quotesApiResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// GetRandomAuthor handles POST requests for getting a random author
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

	var result []structs.QuoteDBModel
	var author structs.AuthorDBModel
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	var shouldDoQuick = false
	//Get Random author
	dbPointer := requestHandler.Db.Table("authors")

	if strings.ToLower(requestBody.Language) != "icelandic" && strings.ToLower(requestBody.Language) != "english" {
		shouldDoQuick = true
	}
	//author from a particular language
	dbPointer = utils.AuthorLanguageSQL(requestBody.Language, dbPointer)
	//** ---------- Paramatere configuratino for DB query ends ---------- **//

	var err error
	if !shouldDoQuick {
		err = dbPointer.Order("random()").First(&author).Error
	} else {
		err = dbPointer.Raw("select * from authors tablesample system(0.25)").First(&author).Error
		//If no author in row
		if err != nil {
			err = dbPointer.Raw("select * from authors tablesample system(0.25)").First(&author).Error

			if err != nil {
				err = dbPointer.Raw("select * from authors tablesample system(0.25)").First(&author).Error
			}
		}
	}

	if err != nil {
		log.Printf("Got error when querying DB, first one, in GetRandomAuthor: %s", err)
		errResponse := structs.ErrorResponse{
			Message: utils.InternalServerError,
		}
		return events.APIGatewayProxyResponse{
			Body:       errResponse.ToString(),
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	dbPointer = requestHandler.Db.Table("quotes").Where("author_id = ?", author.ID)

	//An icelandic quote from the particular/random author
	dbPointer = utils.QuoteLanguageSQL(requestBody.Language, dbPointer)

	err = dbPointer.Limit(requestBody.MaxQuotes).Find(&result).Error

	if err != nil {
		log.Printf("Got error when querying DB, second one, in GetAuthors: %s", err)
		errResponse := structs.ErrorResponse{
			Message: utils.InternalServerError,
		}
		return events.APIGatewayProxyResponse{
			Body:       errResponse.ToString(),
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	quotesAPI := structs.ConvertToQuotesAPIModel(result)
	out, _ := json.Marshal(quotesAPI)
	return events.APIGatewayProxyResponse{
		Body:       string(out),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(theReqHandler.handler)
}
