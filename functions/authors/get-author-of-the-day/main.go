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

// swagger:route POST /authors/aod authors GetAuthorOfTheDay
//
// Get the author of the day (AOD)
//
// Use this route to get the AOD for today for "English" or "icelandic" authors
//
// responses:
//	200: aodResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

//GetAuthorOfTheDay gets the author of the day
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

	var author structs.AodDBModel

	//** ---------- Paramatere configuratino for DB query begins ---------- **//

	//Which table to look for quotes (ice table has icelandic quotes)
	err := utils.AodLanguageSQL(requestBody.Language, requestHandler.Db).
		Where("date = current_date").Limit(1).Scan(&author).Error
	//** ---------- Paramatere configuratino for DB query ends ---------- **//

	if author == (structs.AodDBModel{}) {
		err = utils.AodLanguageSQL(requestBody.Language, requestHandler.Db).
			Order("date desc").Limit(1).Scan(&author).Error
	}

	if err != nil {
		log.Printf("Got error when querying DB in GetAuthorOfTheDay: %s", err)
		errResponse := structs.ErrorResponse{
			Message: utils.InternalServerError,
		}
		return events.APIGatewayProxyResponse{
			Body:       errResponse.ToString(),
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	authorApi := author.ConvertToAPIModel()
	out, _ := json.Marshal(authorApi)
	return events.APIGatewayProxyResponse{
		Body:       string(out),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(theReqHandler.handler)
}
