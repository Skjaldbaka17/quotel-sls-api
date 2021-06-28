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

// swagger:route POST /quotes/random QUOTES GetRandomQuote
// Get a random quote according to the given parameters
// responses:
//  200: topicViewResponse
//  400: incorrectBodyStructureResponse
//  404: notFoundResponse
//  500: internalServerErrorResponse

// GetRandomQuote handles POST requests for getting a random quote
func (requestHandler *RequestHandler) handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	result, err := requestHandler.GetRandomQuoteFromDb(&requestBody)
	if err != nil {
		log.Printf("Got error when querying DB in GetRandomQuote: %s", err)
		return events.APIGatewayProxyResponse{
			Body:       utils.InternalServerError,
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	if result == (structs.TopicViewAPIModel{}) {
		log.Printf("Got error when querying DB in GetRandomQuote: %s", err)
		return events.APIGatewayProxyResponse{
			Body:       "No quote exists that matches the given parameters",
			StatusCode: http.StatusNotFound,
		}, nil
	}

	out, _ := json.Marshal(result)
	return events.APIGatewayProxyResponse{
		Body:       string(out),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(theReqHandler.handler)
}
