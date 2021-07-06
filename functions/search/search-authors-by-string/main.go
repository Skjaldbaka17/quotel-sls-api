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

func generalSearch(requestBody *structs.Request, dbPointer *gorm.DB) *gorm.DB {
	//Order by authorid to have definitive order (when for examplke some names rank the same for similarity), same for why quote_id
	dbPointer = dbPointer.Table("authors").
		Where("( similarity(name, ?) > 0.4)", requestBody.SearchString).
		Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "similarity(name,?) desc", Vars: []interface{}{requestBody.SearchString}, WithoutParentheses: true},
		})

	//Particular language search
	dbPointer = utils.AuthorLanguageSQL(requestBody.Language, dbPointer)
	return dbPointer
}

func search(requestBody *structs.Request, dbPointer *gorm.DB) *gorm.DB {
	//Order by authorid to have definitive order (when for examplke some names rank the same for similarity), same for why quote_id
	dbPointer = dbPointer.Table("authors, plainto_tsquery('english', ?) as plainq", requestBody.SearchString).Select("*, ts_rank(tsv, plainq) as plainrank").
		Where("( tsv @@ plainq )").Order("plainrank desc, id desc")

	//Particular language search
	dbPointer = utils.AuthorLanguageSQL(requestBody.Language, dbPointer)
	return dbPointer
}

// swagger:route POST /search/authors SEARCH SearchAuthorsByString
//
// Authors search. Searching authors by a given search string
//
// responses:
//	200: authorsResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// SearchAuthorsByString handles POST requests to search for authors by a search-string
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

	var results []structs.AuthorDBModel

	var dbPointer *gorm.DB
	for i := 0; i < 3; i++ {
		if i == 1 {
			requestBody.SearchString = requestHandler.CheckForSpellingErrorsInSearchString(requestBody.SearchString, "unique_lexeme_authors")
		} else if i == 2 {
			dbPointer = generalSearch(&requestBody, requestHandler.Db)
		}

		if i != 2 {
			dbPointer = search(&requestBody, requestHandler.Db)
		}

		err := utils.Pagination(requestBody, dbPointer).
			Find(&results).Error
		//  500: internalServerErrorResponse
		if err != nil {
			log.Printf("Got error when querying DB in SearchAuthorsByString: %s", err)
			errResponse := structs.ErrorResponse{
				Message: utils.InternalServerError,
			}
			return events.APIGatewayProxyResponse{
				Body:       errResponse.ToString(),
				StatusCode: http.StatusInternalServerError,
			}, nil
		}

		if len(results) > 0 {
			break
		}
	}

	//Update popularity in background! TODO: add as its own lambda function
	go requestHandler.AuthorsAppearInSearchCountIncrement(results)

	authorsAPI := structs.ConvertToAuthorsAPIModel(results)
	out, _ := json.Marshal(authorsAPI)
	return events.APIGatewayProxyResponse{
		Body:       string(out),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(theReqHandler.handler)
}
