package main

import (
	"encoding/json"
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

// swagger:route POST /meta/nationalities meta GetNationalities
//
// Get nationalities
//
// Use this route to get all authors' nationalities available in the api at any given moment
//
// responses:
//	200: listNationalities
//  500: internalServerErrorResponse

// ListNationalities handles POST requests for getting the disctint nationalities in the database
func (requestHandler *RequestHandler) handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	//Initialize DB if requestHandler.Db = nil
	if errResponse := requestHandler.InitializeDB(); errResponse != (structs.ErrorResponse{}) {
		return events.APIGatewayProxyResponse{
			Body:       errResponse.ToString(),
			StatusCode: errResponse.StatusCode,
		}, nil
	}

	nationalities := []string{}
	err := requestHandler.Db.Table("authors").Select("distinct nationality").Find(&nationalities).Error
	if err != nil {
		errResponse := structs.ErrorResponse{
			Message: utils.InternalServerError,
		}
		return events.APIGatewayProxyResponse{
			Body:       errResponse.ToString(),
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	type Response struct {
		Nationalities []string
	}
	out, _ := json.Marshal(&Response{
		Nationalities: nationalities,
	})
	return events.APIGatewayProxyResponse{
		Body:       string(out),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(theReqHandler.handler)
}
