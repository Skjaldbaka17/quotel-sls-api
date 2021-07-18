package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RequestHandler struct {
	utils.RequestHandler
}

var theReqHandler = RequestHandler{}

func search(requestBody *structs.Request, dbPointer *gorm.DB, searchString string) *gorm.DB {
	table := "quotes"
	//TODO: Validate that this topicId exists
	if len(requestBody.TopicIds) > 0 {
		table = "topicsview"
	}
	dbPointer = dbPointer.Table(table+", plainto_tsquery('english', ?) as plainq ",
		searchString).Select("*, ts_rank(tsv, plainq) as plainrank")

	if len(requestBody.TopicIds) > 0 {
		dbPointer = dbPointer.Where("topic_id in ?", requestBody.TopicIds)
	}

	//Order by authorid to have definitive order (when for examplke some quotes rank the same for plain, phrase, general and similarity)
	dbPointer = dbPointer.Where("( tsv @@ plainq )").Order("plainrank desc, id desc")

	//Particular language search
	dbPointer = utils.QuoteLanguageSQL(requestBody.Language, dbPointer)
	return dbPointer
}

func authorSearch(requestBody *structs.Request, dbPointer *gorm.DB, searchString string) *gorm.DB {
	table := "quotes"
	//TODO: Validate that this topicId exists
	if len(requestBody.TopicIds) > 0 {
		table = "topicsview"
	}
	dbPointer = dbPointer.Table(table).
		Where("( similarity(name, ?) > 0.4)", requestBody.SearchString)

	dbPointer = utils.AuthorLanguageSQL(requestBody.Language, dbPointer)

	if len(requestBody.TopicIds) > 0 {
		dbPointer = dbPointer.Where("topic_id in ?", requestBody.TopicIds)
	}

	//Particular language search
	dbPointer = utils.QuoteLanguageSQL(requestBody.Language, dbPointer)

	dbPointer = dbPointer.
		Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "similarity(name,?) desc, id desc", Vars: []interface{}{requestBody.SearchString}, WithoutParentheses: true},
		})
	return dbPointer
}

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

	var topicResults []structs.QuoteDBModel

	var dbPointer *gorm.DB
	for i := 0; i < 3; i++ {
		if i == 0 {
			dbPointer = search(&requestBody, requestHandler.Db, requestBody.SearchString)
		} else if i == 1 {
			searchString := requestHandler.CheckForSpellingErrorsInSearchString(requestBody.SearchString, "unique_lexeme")
			dbPointer = search(&requestBody, requestHandler.Db, searchString)
		} else if i == 2 {
			dbPointer = authorSearch(&requestBody, requestHandler.Db, requestBody.SearchString)
		}

		err := utils.Pagination(requestBody, dbPointer).
			Find(&topicResults).Error

		if err != nil {
			log.Printf("Got error when querying DB in SearchByString: %s", err)
			errResponse := structs.ErrorResponse{
				Message: utils.InternalServerError,
			}
			return events.APIGatewayProxyResponse{
				Body:       errResponse.ToString(),
				StatusCode: http.StatusInternalServerError,
			}, nil
		}
		if len(topicResults) > 0 {
			break
		}
	}

	if len(topicResults) == 0 {
		errResponse := structs.ErrorResponse{
			Message:    "There are no quotes matching your search. Please check your string for spelling errors etc.",
			StatusCode: http.StatusOK,
		}
		return events.APIGatewayProxyResponse{
			Body:       errResponse.ToString(),
			StatusCode: http.StatusOK,
		}, nil

	}
	//Update popularity in background! TODO: Add as its own lambda function
	go requestHandler.TopicViewAppearInSearchCountIncrement(topicResults)
	apiResults := structs.ConvertToQuotesAPIModel(topicResults)
	out, _ := json.Marshal(apiResults)
	return events.APIGatewayProxyResponse{
		Body:       string(out),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(theReqHandler.handler)
}
