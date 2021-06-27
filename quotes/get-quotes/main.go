package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/Skjaldbaka17/quotel-sls-api/structs"
)

// swagger:route POST /quotes QUOTES GetQuotes
// Get quotes by their ids
//
// responses:
//	200: searchViewsResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	data := structs.Request{}
	err := json.Unmarshal([]byte(request.Body), &data)
	if err != nil {
		fmt.Println(err.Error())
		//invalid character '\'' looking for beginning of object key string
	}
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("GetQuotes123, %v", data),
		StatusCode: 200,
	}, nil
}

// GetQuotes handles POST requests to get the quotes, and their authors, that have the given ids
// func GetQuotes(rw http.ResponseWriter, r *http.Request) {
// 	var requestBody structs.Request
// 	if err := handlers.GetRequestBody(rw, r, &requestBody); err != nil {
// 		return
// 	}
// 	var quotes []structs.SearchViewDBModel
// 	//** ---------- Paramatere configuratino for DB query begins ---------- **//

// 	dbPointer := handlers.Db.Table("searchview").Order("quote_id ASC")
// 	if requestBody.AuthorId > 0 {
// 		dbPointer = dbPointer.
// 			Where("author_id = ?", requestBody.AuthorId)
// 		dbPointer = pagination(requestBody, dbPointer)
// 	} else {
// 		dbPointer = dbPointer.Where("quote_id in ?", requestBody.Ids)
// 	}
// 	//** ---------- Paramatere configuratino for DB query ends ---------- **//

// 	err := dbPointer.Find(&quotes).Error

// 	if err != nil {
// 		rw.WriteHeader(http.StatusInternalServerError)
// 		log.Printf("Got error when querying DB in GetQuotes: %s", err)
// 		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: handlers.InternalServerError})
// 		return
// 	}

// 	//Update popularity in background!
// 	go handlers.DirectFetchQuotesCountIncrement(requestBody.Ids)

// 	searchViewsAPI := structs.ConvertToSearchViewsAPIModel(quotes)
// 	json.NewEncoder(rw).Encode(searchViewsAPI)
// }

func main() {
	lambda.Start(handler)
}
