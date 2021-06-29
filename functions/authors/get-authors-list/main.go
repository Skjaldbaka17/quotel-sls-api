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

// swagger:route POST /authors/list AUTHORS ListAuthors
//
// Get a list of authors according to some ordering / parameters
//
// responses:
//	200: authorsResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// GetAuthorsList handles POST requests to get the authors that fit the parameters
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

	var authors []structs.AuthorDBModel
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	dbPointer := requestHandler.Db.Table("authors")

	dbPointer = utils.AuthorLanguageSQL(requestBody.Language, dbPointer)

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
		dbPointer = dbPointer.Order("count " + orderDirection)
	case "nrofquotes":
		switch strings.ToLower(requestBody.Language) {
		case "english":
			dbPointer = utils.SetMaxMinNumber(requestBody.OrderConfig, "nr_of_english_quotes", orderDirection, dbPointer)
		case "icelandic":
			dbPointer = utils.SetMaxMinNumber(requestBody.OrderConfig, "nr_of_icelandic_quotes", orderDirection, dbPointer)
		default:
			dbPointer = utils.SetMaxMinNumber(requestBody.OrderConfig, "nr_of_icelandic_quotes + nr_of_english_quotes", orderDirection, dbPointer)
		}

	default:
		//Minimum letter to start with (i.e. start from given minimum letter of the alphabet)
		if requestBody.OrderConfig.Minimum != "" {
			dbPointer = dbPointer.Where("initcap(name) >= ?", strings.ToUpper(requestBody.OrderConfig.Minimum))
		}
		//Maximum letter to start with (i.e. end at the given maximum letter of the alphabet)
		if requestBody.OrderConfig.Maximum != "" {
			dbPointer = dbPointer.Where("initcap(name) <= ?", strings.ToUpper(requestBody.OrderConfig.Maximum))
		}
		dbPointer = dbPointer.Order("initcap(name) " + orderDirection)
	}

	//** ---------- Paramatere configuratino for DB query ends---------- **//
	err := utils.Pagination(requestBody, dbPointer).Order("id").
		Find(&authors).
		Error

	if err != nil {
		log.Printf("Got error when querying DB in GetAuthors: %s", err)
		return events.APIGatewayProxyResponse{
			Body:       utils.InternalServerError,
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	//Update popularity in background! TODO: Put into its own Lambda function
	// go handlers.AuthorsAppearInSearchCountIncrement(authors)

	authorsAPI := structs.ConvertToAuthorsAPIModel(authors)
	out, _ := json.Marshal(authorsAPI)
	return events.APIGatewayProxyResponse{
		Body:       string(out),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(theReqHandler.handler)
}
