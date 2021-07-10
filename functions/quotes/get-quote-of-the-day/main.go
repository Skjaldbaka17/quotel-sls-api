package main

import (
	"encoding/json"
	"fmt"
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

// swagger:route POST /quotes/qod QUOTES GetQuoteOfTheDay
// gets the quote of the day
// responses:
//	200: qodResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

//GetQuoteOfTheyDay gets the quote of the day
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

	if requestBody.Language == "" {
		requestBody.Language = "English"
	}

	var quote structs.QodViewDBModel
	var err error
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	dbPointer := requestHandler.QodLanguageSQL(requestBody.Language).Where("date = current_date")
	//** ---------- Paramatere configuratino for DB query ends ---------- **//
	err = dbPointer.Scan(&quote).Error

	if err != nil {
		log.Printf("Got error when querying DB in GetQODs: %s", err)
		errResponse := structs.ErrorResponse{
			Message: utils.InternalServerError,
		}
		return events.APIGatewayProxyResponse{
			Body:       errResponse.ToString(),
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	if (structs.QodViewDBModel{}) == quote {
		fmt.Println("Setting a brand new QOD for today")
		err = requestHandler.SetNewRandomQOD(&requestBody)
		if err != nil {
			log.Printf("Got error when setting new random qod: %s", err)
			errResponse := structs.ErrorResponse{
				Message: utils.InternalServerError,
			}
			return events.APIGatewayProxyResponse{
				Body:       errResponse.ToString(),
				StatusCode: http.StatusInternalServerError,
			}, nil
		}

		return requestHandler.handler(request)
	}

	out, _ := json.Marshal(quote.ConvertToAPIModel())
	return events.APIGatewayProxyResponse{
		Body:       string(out),
		StatusCode: http.StatusOK,
	}, nil

}

func main() {
	lambda.Start(theReqHandler.handler)
}
