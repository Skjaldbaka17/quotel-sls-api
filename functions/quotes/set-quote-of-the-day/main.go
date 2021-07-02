package main

import (
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

// swagger:route POST /quotes/qod/new QUOTES SetQuoteOfTheDay
// Sets the quote of the day for the given dates
// responses:
//	200: successResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

//SetQuoteOfTheyDay sets the quote of the day (is password protected)
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

	if requestBody.Language == "" {
		requestBody.Language = "English"
	}

	if len(requestBody.Qods) == 0 {
		log.Println("Not QODS supplied when setting quote of the day")
		return events.APIGatewayProxyResponse{
			Body:       "Please supply some quotes",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	for _, qod := range requestBody.Qods {
		err := requestHandler.SetQOD(requestBody.Language, qod.Date, qod.Id)
		if err != nil {
			log.Printf("Got error when settin the qod %+v as QOD: %s", qod, err)
			return events.APIGatewayProxyResponse{
				Body:       "Some of the quotes (ids) you supplied are not in " + requestBody.Language,
				StatusCode: http.StatusBadRequest,
			}, nil
		}
	}

	return events.APIGatewayProxyResponse{
		Body:       "Successfully inserted quote of the day!",
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(theReqHandler.handler)
}
