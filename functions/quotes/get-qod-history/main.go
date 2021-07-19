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

// swagger:route POST /quotes/qod/history QUOTES GetQODHistory
// Gets the history for the quotes of the day
// responses:
//	200: qodHistoryResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

//GetQODHistory gets Qod history starting from some point
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

	var quotes []structs.QodDBModel
	var err error
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	if requestBody.TopicId > 0 {
		dbPointer := requestHandler.Db.Table("qods").Where("topic_id = ?", requestBody.TopicId)

		//Not maximum because then possibility of endless cycle with the if statement below!
		if requestBody.Minimum != "" {
			dbPointer = dbPointer.Where("date >= ?", requestBody.Minimum)
		}
		dbPointer = dbPointer.Where("date <= current_date").Order("date DESC")
		//** ---------- Paramatere configuratino for DB query ends ---------- **//
		err = dbPointer.Find(&quotes).Error
	} else {
		dbPointer := requestHandler.QodLanguageSQL(requestBody.Language)

		//Not maximum because then possibility of endless cycle with the if statement below!
		if requestBody.Minimum != "" {
			dbPointer = dbPointer.Where("date >= ?", requestBody.Minimum)
		}
		dbPointer = dbPointer.Where("date <= current_date").Where("topic_id = 0").Order("date DESC")
		//** ---------- Paramatere configuratino for DB query ends ---------- **//
		err = dbPointer.Find(&quotes).Error
	}

	if err != nil {
		log.Printf("Got error when querying DB in GetQODHistory: %s", err)
		errResponse := structs.ErrorResponse{
			Message: utils.InternalServerError,
		}
		return events.APIGatewayProxyResponse{
			Body:       errResponse.ToString(),
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	qodHistoryAPI := structs.ConvertToQodAPIModel(quotes)
	out, _ := json.Marshal(qodHistoryAPI)
	return events.APIGatewayProxyResponse{
		Body:       string(out),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(theReqHandler.handler)
}
