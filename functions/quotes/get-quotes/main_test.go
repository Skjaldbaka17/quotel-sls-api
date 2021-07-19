package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
	"github.com/aws/aws-lambda-go/events"
)

var testingHandler = RequestHandler{}

//Returns AODs and AODICEs, in that order, put into the DB
func Setup(handler *RequestHandler, t *testing.T) ([]structs.AuthorDBModel, []structs.QuoteDBModel) {
	handler.InitializeDB()
	var authors []structs.AuthorDBModel
	var quotes []structs.QuoteDBModel

	err := handler.Db.Table("authors").Limit(3).Find(&authors).Error
	if err != nil {
		t.Fatalf("Setup error: %s", err)
	}
	err = handler.Db.Table("quotes").Limit(3).Find(&quotes).Error
	if err != nil {
		t.Fatalf("Setup error: %s", err)
	}

	//CleanUp
	t.Cleanup(func() {
		handler.Db.Table("authors").Model(&authors).Update("count", 0)
		handler.Db.Table("quotes").Model(&quotes).Update("count", 0)
	})

	return authors, quotes
}

func GetRequest(jsonStr string, obj interface{}, t *testing.T) string {
	response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: jsonStr})
	if err != nil {
		t.Fatalf("Expected 3 quotes but got an error: %+v", err)
	}
	json.Unmarshal([]byte(response.Body), &obj)
	return response.Body
}
func TestHandler(t *testing.T) {

	authors, quotes := Setup(&testingHandler, t)

	t.Run("Time Test for getting quotes", func(t *testing.T) {
		maxTime := 25
		t.Run("should return Quotes by ids", func(t *testing.T) {
			start := time.Now()
			var jsonStr = fmt.Sprintf(`{"ids":  [%d,%d,%d]}`, quotes[0].Id, quotes[1].Id, quotes[2].Id)
			GetRequest(jsonStr, nil, t)

			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(maxTime) {
				t.Fatalf("Expected getting history of quotes to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}
		})

		t.Run("should get Quotes for author by his id", func(t *testing.T) {
			start := time.Now()
			var jsonStr = fmt.Sprintf(`{"authorId":  %d}`, authors[0].ID)
			GetRequest(jsonStr, nil, t)
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(maxTime) {
				t.Fatalf("Expected getting history of quotes to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}

		})
	})

	t.Run("Get quotes", func(t *testing.T) {

		t.Run("should return Quotes by ids", func(t *testing.T) {
			var jsonStr = fmt.Sprintf(`{"ids":  [%d,%d,%d]}`, quotes[0].Id, quotes[1].Id, quotes[2].Id)
			var respQuotes []structs.QuoteAPIModel
			GetRequest(jsonStr, &respQuotes, t)

			if len(respQuotes) != len(quotes) {
				t.Fatalf("got list of length %d but expected list of length %d", len(respQuotes), len(quotes))
			}

			for _, quote := range quotes {
				for j, testQuote := range respQuotes {
					if testQuote.QuoteId == quote.Id {
						break
					}
					if j == len(quotes)-1 {
						t.Fatalf("expected quote with id %d to be amongst the returned quotes but got %+v", quote.Id, respQuotes)
					}
				}
			}
		})

		t.Run("should get Quotes for author by his id", func(t *testing.T) {
			var jsonStr = fmt.Sprintf(`{"authorId":  %d}`, authors[0].ID)
			var respQuotes []structs.QuoteAPIModel
			responseBod := GetRequest(jsonStr, &respQuotes, t)

			if len(respQuotes) == 0 {
				t.Fatalf("got list of length 0 but expected some quotes, response : %s", responseBod)
			}

			if respQuotes[0].Id != authors[0].ID {
				t.Fatalf("got quotes for author with id %d but expected quotes for the author with id %d, respObj: %s", respQuotes[0].Id, authors[0].ID, responseBod)
			}
		})

	})

}
