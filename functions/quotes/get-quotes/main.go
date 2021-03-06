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

// swagger:route POST /quotes quotes GetQuotes
//
// Get quotes by ids
//
// Use this route to either get quotes straight from their ids or use it to get all the quotes from a particular author (by supplying the author's id)
//
// responses:
//	200: quotesApiResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// GetQuotes handles POST requests to get the quotes, and their authors, that have the given ids
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

	var quotes []structs.QuoteDBModel
	//** ---------- Paramatere configuratino for DB query begins ---------- **//

	dbPointer := requestHandler.Db.Table("quotes").Order("id ASC")
	if requestBody.AuthorId > 0 {
		dbPointer = dbPointer.
			Where("author_id = ?", requestBody.AuthorId)
		dbPointer = utils.Pagination(requestBody, dbPointer)
	} else {
		dbPointer = dbPointer.Where("id in ?", requestBody.Ids)
	}
	//** ---------- Paramatere configuration for DB query ends ---------- **//

	err := dbPointer.Find(&quotes).Error

	if err != nil {
		log.Printf("Got error when querying DB in GetQuotes: %s", err)
		errResponse := structs.ErrorResponse{
			Message: utils.InternalServerError,
		}
		return events.APIGatewayProxyResponse{
			Body:       errResponse.ToString(),
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	//Update popularity in background! TODO: PUT IN ITS OWN LAMBDA FUNCTION!
	go requestHandler.DirectFetchQuotesCountIncrement(requestBody.Ids)

	quotesApi := structs.ConvertToQuotesAPIModel(quotes)

	out, _ := json.Marshal(quotesApi)
	return events.APIGatewayProxyResponse{
		Body:       string(out),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(theReqHandler.handler)
}
