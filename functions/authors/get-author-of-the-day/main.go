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

// swagger:route POST /authors/aod AUTHORS GetAuthorOfTheDay
// Gets the author of the day
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
	var err error

	//** ---------- Paramatere configuratino for DB query begins ---------- **//

	//Which table to look for quotes (ice table has icelandic quotes)
	dbPointer := utils.AodLanguageSQL(requestBody.Language, requestHandler.Db).
		Where("date = current_date")
	//** ---------- Paramatere configuratino for DB query ends ---------- **//

	err = dbPointer.Scan(&author).Error

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

	if (structs.AodDBModel{}) == author {
		err = requestHandler.SetNewRandomAOD(requestBody.Language)
		if err != nil {
			log.Printf("Got error when setting new random AOD in GetAuthorOfTheDay: %s", err)
			errResponse := structs.ErrorResponse{
				Message: utils.InternalServerError,
			}
			return events.APIGatewayProxyResponse{
				Body:       errResponse.ToString(),
				StatusCode: http.StatusInternalServerError,
			}, nil
		}

		return requestHandler.handler(request) //Dangerous? possibility of endless cycle? Only iff the setNewRandomAOD fails in some way. Or the date is not saved correctly into the DB?
	}

	out, _ := json.Marshal(author)
	return events.APIGatewayProxyResponse{
		Body:       string(out),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(theReqHandler.handler)
}
