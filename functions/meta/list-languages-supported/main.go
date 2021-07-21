package main

import (
	"encoding/json"
	"net/http"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RequestHandler struct {
	utils.RequestHandler
}

var theReqHandler = RequestHandler{}
var languages = []string{"English", "Icelandic"}

// swagger:route POST /languages meta GetLanguages
//
// Get languages
//
// Use this route to get the languages supported by the api at any moment
//
// responses:
//	200: listOfStrings

// ListLanguages handles POST requests for getting the languages supported by the api
func (requestHandler *RequestHandler) handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	type response = struct {
		Languages []string `json:"languages"`
	}

	out, _ := json.Marshal(&response{
		Languages: languages,
	})
	return events.APIGatewayProxyResponse{
		Body:       string(out),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(theReqHandler.handler)
}
