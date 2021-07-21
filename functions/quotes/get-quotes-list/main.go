package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RequestHandler struct {
	utils.RequestHandler
}

var theReqHandler = RequestHandler{}

// swagger:route POST /quotes/list quotes GetQuotesList
//
// List quotes
//
// Use this route to get a list of quotes according to some ordering / parameters -- list them based on length of quote, popularity or by their id
//
// responses:
//	200: quotesApiResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// GetQuotesList handles POST requests to get the quotes that fit the parameters
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

	var quotes []structs.QuoteDBModel
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	dbPointer := requestHandler.Db.Table("quotes")
	dbPointer = utils.QuoteLanguageSQL(requestBody.Language, dbPointer)

	orderDirection := "ASC"
	if requestBody.OrderConfig.Reverse {
		orderDirection = "DESC"
	}

	switch strings.ToLower(requestBody.OrderConfig.OrderBy) {
	case "popularity": //TODO: add popularity ordering
		orderDirection = "DESC"
		if requestBody.OrderConfig.Reverse {
			orderDirection = "ASC"
		}

		//where > 0 just because otherwise most quotes have 0 count and therefore it takes a very long time to fetch even thought indexed, because most have same value in count
		//note: this is only problem in the beginning
		dbPointer = dbPointer.Where("count > 0").Order("count " + orderDirection)
	case "length":
		dbPointer = utils.SetMaxMinNumber(requestBody.OrderConfig, "quote_length", orderDirection, dbPointer)
	default:
		dbPointer = utils.SetMaxMinNumber(requestBody.OrderConfig, "id", orderDirection, dbPointer)
	}

	//** ---------- Paramatere configuratino for DB query ends ---------- **//

	err := utils.Pagination(requestBody, dbPointer).Order("id").
		Find(&quotes).
		Error

	if err != nil {
		log.Printf("Got error when querying DB in GetQuotesList: %s", err)
		errResponse := structs.ErrorResponse{
			Message: utils.InternalServerError,
		}
		return events.APIGatewayProxyResponse{
			Body:       errResponse.ToString(),
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	//Update popularity in background! TODO: PUT IN ITS OWN LAMBDA FUNCTION!
	go requestHandler.QuotesAppearInSearchCountIncrement(quotes)
	quotesAPI := structs.ConvertToQuotesAPIModel(quotes)
	out, _ := json.Marshal(quotesAPI)
	return events.APIGatewayProxyResponse{
		Body:       string(out),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(theReqHandler.handler)
}
