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

// swagger:route POST /search search SearchByString
//
// Search general
//
// Use this route to search for quotes / authors by a general full test search that searches both in the names of the authors and the quotes themselves
//
// responses:
//  200: quotesApiResponse
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
			//Make nicer: --but what this is doing:
			/*
				First do authorsearch: do a general similarity search through the authors table (only 30.000 rows)...
				If you find a matchin author: Fetch his/hers quotes from "quotes" in quoteId order and return.
				If no match just break the loop and return an empty array signaling nothing was found to match the searchstring
			*/
			var results []structs.AuthorDBModel
			dbPointer = authorSearch(&requestBody, requestHandler.Db, requestBody.SearchString)
			err := utils.Pagination(requestBody, dbPointer).
				Find(&results).Error
			if err != nil {
				log.Printf("Got error when querying authorDB in SearchByString: %s", err)
				errResponse := structs.ErrorResponse{
					Message: utils.InternalServerError,
				}
				return events.APIGatewayProxyResponse{
					Body:       errResponse.ToString(),
					StatusCode: http.StatusInternalServerError,
				}, nil
			}
			if len(results) > 0 {
				dbPointer = requestHandler.Db.Table("quotes").Where("author_id = ?", results[0].ID).Order("id")
				err = utils.Pagination(requestBody, dbPointer).
					Find(&topicResults).Error
				if err != nil {
					log.Printf("Got error when querying authorDB in SearchByString: %s", err)
					errResponse := structs.ErrorResponse{
						Message: utils.InternalServerError,
					}
					return events.APIGatewayProxyResponse{
						Body:       errResponse.ToString(),
						StatusCode: http.StatusInternalServerError,
					}, nil
				}
				break
			}
			break
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

	// if len(topicResults) == 0 {
	// 	errResponse := structs.ErrorResponse{
	// 		Message:    "There are no quotes matching your search. Please check your string for spelling errors etc.",
	// 		StatusCode: http.StatusOK,
	// 	}
	// 	return events.APIGatewayProxyResponse{
	// 		Body:       errResponse.ToString(),
	// 		StatusCode: http.StatusOK,
	// 	}, nil

	// }
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
