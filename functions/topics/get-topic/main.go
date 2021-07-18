package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gorm.io/gorm/clause"
)

type RequestHandler struct {
	utils.RequestHandler
}

var theReqHandler = RequestHandler{}

// swagger:route POST /topic TOPICS GetTopic
// Get quotes from a particular topic
// responses:
//	200: topicViewsResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// GetTopic handles POST requests for getting quotes from a particular topic
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

	var results []structs.QuoteDBModel
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	//Order by quoteid to have definitive order (when for examplke some quotes rank the same for plain, phrase and general)
	dbPoint := requestHandler.Db.Table("topicsview").Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "id DESC", Vars: []interface{}{}, WithoutParentheses: true},
	})

	if requestBody.Topic != "" {
		dbPoint = dbPoint.Where("lower(topic_name) = lower(?)", requestBody.Topic)
	} else {
		dbPoint = dbPoint.Where("topic_id = ?", requestBody.TopicId)
	}

	//** ---------- Paramatere configuratino for DB query ends ---------- **//
	err := utils.Pagination(requestBody, dbPoint).Find(&results).Error

	if err != nil {
		log.Printf("Got error when querying DB in GetTopic: %s", err)
		errResponse := structs.ErrorResponse{
			Message: utils.InternalServerError,
		}
		return events.APIGatewayProxyResponse{
			Body:       errResponse.ToString(),
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	//Update popularity in background! TODO: Add as its own lambda function
	go requestHandler.DirectFetchTopicCountIncrement(requestBody.Id, requestBody.Topic)
	quotesAPI := structs.ConvertToQuotesAPIModel(results)
	out, _ := json.Marshal(quotesAPI)
	return events.APIGatewayProxyResponse{
		Body:       string(out),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(theReqHandler.handler)
}
