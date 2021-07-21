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

// swagger:route POST /authors authors GetAuthors
//
// Get authors by ids
//
// Use this route to get the authors that have the given ids.
//
// responses:
//	200: authorsResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// Get Authors handles POST requests to get the authors, and their quotes, that have the given ids
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

	var authors []structs.AuthorDBModel
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	err := requestHandler.Db.Table("authors").
		Where("id in (?)", requestBody.Ids).
		Scan(&authors).
		Error
	//** ---------- Paramatere configuratino for DB query ends ---------- **//

	if err != nil {
		log.Printf("Got error when querying DB in GetAuthorsById: %s", err)
		errResponse := structs.ErrorResponse{
			Message: utils.InternalServerError,
		}
		return events.APIGatewayProxyResponse{
			Body:       errResponse.ToString(),
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	//Update popularity in background! TODO: PUT IN ITS OWN LAMBDA FUNCTION!
	go requestHandler.DirectFetchAuthorsCountIncrement(requestBody.Ids)
	authorsAPI := structs.ConvertToAuthorsAPIModel(authors)
	out, _ := json.Marshal(authorsAPI)
	return events.APIGatewayProxyResponse{
		Body:       string(out),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(theReqHandler.handler)
}
