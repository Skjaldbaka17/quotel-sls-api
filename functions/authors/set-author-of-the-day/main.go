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

// swagger:route POST /authors/aod/new AUTHORS SetAuthorOfTheDay
//
// sets the author of the day for the given dates
//
// responses:
//	200: successResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

//SetAuthorOfTheDay sets the author of the day.
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

	if len(requestBody.Aods) == 0 {
		log.Println("No author supplied")
		return events.APIGatewayProxyResponse{
			Body:       "Please supply some authors",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	for _, aod := range requestBody.Aods {
		err := requestHandler.SetAOD(requestBody.Language, aod.Date, aod.Id)
		if err != nil {
			log.Printf("Got err when setting the AOD for %+v: %s", aod, err)
			return events.APIGatewayProxyResponse{
				Body:       "Some of the authors (ids) you supplied do not have " + requestBody.Language + " quotes",
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
