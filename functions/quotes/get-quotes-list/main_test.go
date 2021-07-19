package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
	"github.com/aws/aws-lambda-go/events"
)

func Setup(handler *RequestHandler, t *testing.T) structs.QuoteDBModel {
	handler.InitializeDB()
	var quote structs.QuoteDBModel
	err := handler.Db.Table("quotes").First(&quote).Error
	if err != nil {
		t.Fatalf("got error in setup: %s", err)
	}

	err = handler.Db.Table("quotes").Model(&quote).Update("count", 10000).Error
	if err != nil {
		t.Fatalf("got error in setup: %s", err)
	}

	t.Cleanup(func() {
		handler.Db.Table("quotes").Model(&quote).Update("count", 0)
	})
	return quote
}
func TestHandler(t *testing.T) {
	var testingHandler = RequestHandler{}
	quote := Setup(&testingHandler, t)
	t.Run("Time Test for getting quotes", func(t *testing.T) {
		maxTime := 25
		longTime := 230
		fkkingTooLong := 1000
		t.Run("Should return first 50 quotes (by quoteId)", func(t *testing.T) {
			start := time.Now()
			pageSize := 50
			var jsonStr = fmt.Sprintf(`{"pageSize": %d}`, pageSize)

			testingHandler.handler(events.APIGatewayProxyRequest{Body: jsonStr})

			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(maxTime) {
				t.Fatalf("Expected getting history of quotes to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}
		})
		t.Run("Should return first quotes in reverse quoteId order (i.e. first quote has id larger than 639.028)", func(t *testing.T) {
			start := time.Now()
			var jsonStr = fmt.Sprintf(`{"orderConfig":{"reverse":%s}}`, "true")

			testingHandler.handler(events.APIGatewayProxyRequest{Body: jsonStr})

			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(maxTime) {
				t.Fatalf("Expected getting history of quotes to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}
		})
		t.Run("Should return first quotes starting from id 300.000  (i.e. greater than or equal to 300.000)", func(t *testing.T) {
			start := time.Now()
			minimum := 300000
			orderBy := "id"
			var jsonStr = fmt.Sprintf(`{"orderConfig":{"orderBy":"%s","minimum":"%d"}}`, orderBy, minimum)

			testingHandler.handler(events.APIGatewayProxyRequest{Body: jsonStr})
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(maxTime) {
				t.Fatalf("Expected getting history of quotes to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}

		})
		t.Run("Should return first quotes with quote-length at least 10 an most 11", func(t *testing.T) {
			start := time.Now()

			minimum := 10
			maximum := 11
			var jsonStr = fmt.Sprintf(`{"orderConfig":{"orderBy":"length","maximum":"%d", "minimum":"%d"}}`, maximum, minimum)

			testingHandler.handler(events.APIGatewayProxyRequest{Body: jsonStr})
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(fkkingTooLong) {
				t.Fatalf("Expected getting history of quotes to take less than %dms but it took %dms", fkkingTooLong, duration.Milliseconds())
			}

		})
		t.Run("Should return first 50 quotes (ordered by most popular, i.e. DESC count)", func(t *testing.T) {
			start := time.Now()
			var jsonStr = fmt.Sprintf(`{"orderConfig":{"orderBy":"%s"}}`, "popularity")
			testingHandler.handler(events.APIGatewayProxyRequest{Body: jsonStr})
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(longTime) {
				t.Fatalf("Expected getting history of quotes to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}
		})
		t.Run("Should return first 50 quotes in reverse popularity order (i.e. least popular first i.e. ASC count)", func(t *testing.T) {
			start := time.Now()
			var jsonStr = fmt.Sprintf(`{"orderConfig":{"orderBy":"%s","reverse":true}}`, "popularity")

			testingHandler.handler(events.APIGatewayProxyRequest{Body: jsonStr})
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(longTime) {
				t.Fatalf("Expected getting history of quotes to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}

		})

	})

	t.Run("Get quotes", func(t *testing.T) {

		t.Run("Quoteslist Test", func(t *testing.T) {

			t.Run("Should return first 50 quotes (by quoteId)", func(t *testing.T) {
				pageSize := 50
				var jsonStr = fmt.Sprintf(`{"pageSize": %d}`, pageSize)

				response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: jsonStr})
				if err != nil {
					t.Fatalf("Expected 3 quotes but got an error: %+v", err)
				}

				var respQuotes []structs.QuoteAPIModel
				json.Unmarshal([]byte(response.Body), &respQuotes)

				if len(respQuotes) != 50 {
					t.Fatalf("got list of length %d, but expected list of length %d", len(respQuotes), pageSize)
				}

				firstQuote := respQuotes[0]
				if firstQuote.QuoteId != 1 {
					t.Fatalf("got %d, want quote with id 1. Resp: %+v", firstQuote.QuoteId, firstQuote)
				}

			})
			t.Run("Should return first quotes, in Icelandic", func(t *testing.T) {
				language := "icelandic"
				var jsonStr = fmt.Sprintf(`{"language": "%s"}`, language)

				response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: jsonStr})
				if err != nil {
					t.Fatalf("Expected 3 quotes but got an error: %+v", err)
				}

				var respQuotes []structs.QuoteAPIModel
				json.Unmarshal([]byte(response.Body), &respQuotes)

				firstQuote := respQuotes[0]

				if !firstQuote.IsIcelandic {
					t.Fatalf("got %+v, but expected a quote in Icelandic.", firstQuote)
				}

			})
			t.Run("Should return first quotes in reverse quoteId order (i.e. first quote has id larger than 639.028)", func(t *testing.T) {

				var jsonStr = fmt.Sprintf(`{"orderConfig":{"reverse":%s}}`, "true")

				response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: jsonStr})
				if err != nil {
					t.Fatalf("Expected 3 quotes but got an error: %+v", err)
				}

				var respQuotes []structs.QuoteAPIModel
				json.Unmarshal([]byte(response.Body), &respQuotes)

				firstQuote := respQuotes[0]

				if firstQuote.QuoteId < 639028 {
					t.Fatalf("got %+v, but want quote with larger quoteid i.e. want last quote in db", firstQuote)
				}

			})
			t.Run("Should return first quotes starting from id 300.000  (i.e. greater than or equal to 300.000)", func(t *testing.T) {
				minimum := 300000
				orderBy := "id"
				var jsonStr = fmt.Sprintf(`{"orderConfig":{"orderBy":"%s","minimum":"%d"}}`, orderBy, minimum)

				response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: jsonStr})
				if err != nil {
					t.Fatalf("Expected 3 quotes but got an error: %+v", err)
				}

				var respQuotes []structs.QuoteAPIModel
				json.Unmarshal([]byte(response.Body), &respQuotes)

				firstQuote := respQuotes[0]

				if int(firstQuote.QuoteId) < minimum {
					t.Fatalf("got %+v, want quote that has id larger or equal to 300.000", firstQuote)
				}

			})
			t.Run("Should return quotes with less than or equal to 5 letters in the quote", func(t *testing.T) {

				maximum := 5
				var jsonStr = fmt.Sprintf(`{"orderConfig":{"orderBy":"length","maximum":"%d"}}`, maximum)

				response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: jsonStr})
				if err != nil {
					t.Fatalf("Expected 3 quotes but got an error: %+v", err)
				}

				var respQuotes []structs.QuoteAPIModel
				json.Unmarshal([]byte(response.Body), &respQuotes)

				firstQuote := respQuotes[0]

				if len(firstQuote.Quote) > 5 {
					t.Fatalf("got %+v, but expected a quote that has no more than 5 letters", firstQuote)
				}

			})
			t.Run("Should return first quotes with quote-length at least 10 an most 11", func(t *testing.T) {

				minimum := 10
				maximum := 11
				var jsonStr = fmt.Sprintf(`{"orderConfig":{"orderBy":"length","maximum":"%d", "minimum":"%d"}}`, maximum, minimum)

				response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: jsonStr})
				if err != nil {
					t.Fatalf("Expected 3 quotes but got an error: %+v", err)
				}

				var respQuotes []structs.QuoteAPIModel
				json.Unmarshal([]byte(response.Body), &respQuotes)

				firstQuote := respQuotes[0]

				if len(firstQuote.Quote) != 10 {
					t.Fatalf("got %+v, but expected a quote that has no fewer than 10 letters", firstQuote)
				}

			})
			t.Run("Should return first Quotes with less than letters in the quote in total in reversed order (start with those quotes of length 10)", func(t *testing.T) {

				maximum := 10
				var jsonStr = fmt.Sprintf(`{"orderConfig":{"orderBy":"length","maximum":"%d","reverse":true}}`, maximum)

				response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: jsonStr})
				if err != nil {
					t.Fatalf("Expected 3 quotes but got an error: %+v", err)
				}

				var respQuotes []structs.QuoteAPIModel
				json.Unmarshal([]byte(response.Body), &respQuotes)

				firstQuote := respQuotes[0]

				if len(firstQuote.Quote) != 10 {
					t.Fatalf("got %+v, but expected a quote that has 10 letters", firstQuote)
				}

			})
			t.Run("Should return first 50 quotes (ordered by most popular, i.e. DESC count)", func(t *testing.T) {
				var jsonStr = fmt.Sprintf(`{"orderConfig":{"orderBy":"%s"}}`, "popularity")
				response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: jsonStr})
				if err != nil {
					t.Fatalf("Expected 3 quotes but got an error: %+v", err)
				}

				var respQuotes []structs.QuoteAPIModel
				json.Unmarshal([]byte(response.Body), &respQuotes)

				firstQuote := respQuotes[0]

				if firstQuote.QuoteId != quote.Id {
					t.Fatalf("got %+v, but expected a quote that has more than 0 popularity count", firstQuote)
				}
			})
			t.Run("Should return first 50 quotes in reverse popularity order (i.e. least popular first i.e. ASC count)", func(t *testing.T) {

				var jsonStr = fmt.Sprintf(`{"orderConfig":{"orderBy":"%s","reverse":true}}`, "popularity")

				response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: jsonStr})
				if err != nil {
					t.Fatalf("Expected 3 quotes but got an error: %+v", err)
				}

				var respQuotes []structs.QuoteAPIModel
				json.Unmarshal([]byte(response.Body), &respQuotes)

				firstQuote := respQuotes[0]

				//Useless test, this field is always zero for the api. -- maybe change that?
				if firstQuote.Count != 0 {
					t.Fatalf("got %+v, but expected an author that has 0 popularity count", firstQuote)
				}

			})
			t.Run("Should return first 100 Quotes", func(t *testing.T) {
				pageSize := 100
				var jsonStr = fmt.Sprintf(`{"pageSize":%d}`, pageSize)

				response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: jsonStr})
				if err != nil {
					t.Fatalf("Expected 3 quotes but got an error: %+v", err)
				}

				var respQuotes []structs.QuoteAPIModel
				json.Unmarshal([]byte(response.Body), &respQuotes)

				if len(respQuotes) != 100 {
					t.Fatalf("got %d nr of quotes, but expected %d quotes", len(respQuotes), pageSize)
				}
			})
			t.Run("Should return the next 50 quotes starting from quoteId 250.000 (i.e. pagination, page 1, quoteId order)", func(t *testing.T) {

				pageSize := 100
				minimum := 250000
				var jsonStr = fmt.Sprintf(`{"pageSize":%d, "orderConfig":{"minimum":"%d"}}`, pageSize, minimum)

				response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: jsonStr})
				if err != nil {
					t.Fatalf("Expected quotes but got an error: %+v", err)
				}

				var respQuotes []structs.QuoteAPIModel
				json.Unmarshal([]byte(response.Body), &respQuotes)

				objToFetch := respQuotes[50]

				if int(respQuotes[0].QuoteId) < minimum {
					t.Fatalf("got %+v, but expected quote with a higher quoteid than %d", respQuotes[0], minimum)
				}

				pageSize = 50
				page := 1
				jsonStr = fmt.Sprintf(`{"pageSize":%d, "page":%d, "orderConfig":{"minimum":"%d"}}`, pageSize, page, minimum)

				response, err = testingHandler.handler(events.APIGatewayProxyRequest{Body: jsonStr})
				if err != nil {
					t.Fatalf("Expected quotes but got an error: %+v", err)
				}
				json.Unmarshal([]byte(response.Body), &respQuotes)

				if objToFetch.QuoteId != respQuotes[0].QuoteId {
					t.Fatalf("got %+v, but expected %+v", respQuotes[0], objToFetch)
				}

			})

		})

	})

}
