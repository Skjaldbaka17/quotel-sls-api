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

	var quote structs.QodDBModel
	var err error
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	if requestBody.TopicId > 0 {
		err = requestHandler.Db.Table("qods").Where("topic_id = ?", requestBody.TopicId).Where("date = current_date").Limit(1).Scan(&quote).Error
	} else {
		err = requestHandler.QodLanguageSQL(requestBody.Language).Where("topic_id = 0").Where("date = current_date").Limit(1).Scan(&quote).Error
	}
	//** ---------- Paramatere configuratino for DB query ends ---------- **//

	//If date for today has not been set then just fetch the newest qod
	if quote == (structs.QodDBModel{}) {
		if requestBody.TopicId > 0 {
			err = requestHandler.Db.Table("qods").Where("topic_id = ?", requestBody.TopicId).Order("date desc").Limit(1).Scan(&quote).Error
		} else {
			err = requestHandler.QodLanguageSQL(requestBody.Language).Where("topic_id = 0").Order("date desc").Limit(1).Scan(&quote).Error
		}
	}

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

	out, _ := json.Marshal(quote.ConvertToAPIModel())
	return events.APIGatewayProxyResponse{
		Body:       string(out),
		StatusCode: http.StatusOK,
	}, nil

}

func main() {
	lambda.Start(theReqHandler.handler)
}
