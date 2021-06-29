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

// swagger:route POST /search SEARCH SearchByString
// Search for quotes / authors by a general string-search that searches both in the names of the authors and the quotes themselves
//
// responses:
//  200: topicViewsResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// SearchByString handles POST requests to search for quotes / authors by a search-string
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

	var topicResults []structs.TopicViewDBModel
	//** ---------- Paramatere configuratino for DB query begins ---------- **//
	dbPointer := requestHandler.GetBasePointer(requestBody)
	//Order by authorid to have definitive order (when for examplke some quotes rank the same for plain, phrase, general and similarity)
	dbPointer = dbPointer.
		Where("( tsv @@ plainq OR tsv @@ phraseq OR ? % ANY(STRING_TO_ARRAY(name,' ')) OR tsv @@ generalq)", requestBody.SearchString).
		Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "phraserank DESC,similarity(name, ?) DESC, plainrank DESC, generalrank DESC, author_id DESC", Vars: []interface{}{requestBody.SearchString}, WithoutParentheses: true},
		})

	//Particular language search
	dbPointer = utils.QuoteLanguageSQL(requestBody.Language, dbPointer)
	//** ---------- Paramatere configuratino for DB query ends ---------- **//
	err := utils.Pagination(requestBody, dbPointer).
		Find(&topicResults).Error

	if err != nil {
		log.Printf("Got error when querying DB in SearchByString: %s", err)
		return events.APIGatewayProxyResponse{
			Body:       utils.InternalServerError,
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	//Update popularity in background! TODO: Add as its own lambda function
	// go handlers.TopicViewAppearInSearchCountIncrement(topicResults)
	apiResults := structs.ConvertToTopicViewsAPIModel(topicResults)
	out, _ := json.Marshal(apiResults)
	return events.APIGatewayProxyResponse{
		Body:       string(out),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(theReqHandler.handler)
}
