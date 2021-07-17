package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
	"github.com/aws/aws-lambda-go/events"
)

//Returns AODs and AODICEs, in that order, put into the DB
func Setup(handler *RequestHandler, t *testing.T) (structs.QodDBModel, structs.QodDBModel) {
	handler.InitializeDB()

	var englishQuote structs.QuoteDBModel
	var icelandicQuote structs.QuoteDBModel
	err := handler.Db.Table("quotes").Not("is_icelandic").Limit(1).Find(&englishQuote).Error
	if err != nil {
		t.Fatalf("Setup error: %s", err)
	}
	err = handler.Db.Table("quotes").Where("is_icelandic").Limit(3).Find(&icelandicQuote).Error
	if err != nil {
		t.Fatalf("Setup error: %s", err)
	}

	year, month, day := time.Now().Date()
	today := fmt.Sprintf("%d-%d-%d", year, month, day)
	QOD := englishQuote.ConvertToQODDBModel(today)
	QODICE := icelandicQuote.ConvertToQODDBModel(today)
	err = handler.Db.Table("qods").Create(&QOD).Error
	if err != nil {
		t.Fatalf("Setup error 2: %s", err)
	}
	err = handler.Db.Table("qodices").Create(&QODICE).Error
	if err != nil {
		t.Fatalf("Setup error 2: %s", err)
	}

	//CleanUp
	t.Cleanup(func() {
		handler.Db.Exec("delete from qods")
		handler.Db.Exec("delete from qodices")
	})

	return QOD, QODICE
}
func TestHandler(t *testing.T) {
	var testingHandler = RequestHandler{}
	QOD, QODICE := Setup(&testingHandler, t)

	t.Run("Time Test for getting qod", func(t *testing.T) {
		maxTime := 50
		t.Run("Time: Should get QOD", func(t *testing.T) {
			start := time.Now()
			//Get History:
			jsonStr := []byte(fmt.Sprintf(`{"language":"%s"}`, "english"))
			_, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected the QOD but got an error: %+v", err)
			}
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(maxTime) {
				t.Fatalf("Expected getting history of quotes to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}
		})
	})

	t.Run("Get quotes", func(t *testing.T) {
		t.Run("Should get Quote of the day", func(t *testing.T) {
			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: "{}"})
			if err != nil {
				t.Fatalf("Expected the QOD but got an error: %+v", err)
			}

			var quote structs.QodAPIModel
			json.Unmarshal([]byte(response.Body), &quote)

			if quote.QuoteId != QOD.QuoteId {
				t.Fatalf("Expected the quote for today that setup just inserted, i.e. id %d but got quote with id %d", QOD.QuoteId, quote.QuoteId)
			}

		})

		t.Run("Should get English Quote of the day", func(t *testing.T) {
			var jsonStr = []byte(fmt.Sprintf(`{"language":"%s"}`, "english"))
			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected the QOD but got an error: %+v", err)
			}

			var quote structs.QodAPIModel
			json.Unmarshal([]byte(response.Body), &quote)

			if quote.QuoteId != QOD.QuoteId {
				t.Fatalf("Expected the quote for today that setup just inserted, i.e. id %d but got quote with id %d", QOD.QuoteId, quote.QuoteId)
			}

		})

		t.Run("Should get Icelandic Quote of the day", func(t *testing.T) {
			var jsonStr = []byte(fmt.Sprintf(`{"language":"%s"}`, "icelandic"))
			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected the QODICE but got an error: %+v", err)
			}

			var quote structs.QodAPIModel
			json.Unmarshal([]byte(response.Body), &quote)

			if quote.QuoteId != QODICE.QuoteId {
				t.Fatalf("Expected the icelandic quote for today that setup just inserted, i.e. id %d but got quote with id %d", QODICE.QuoteId, quote.QuoteId)
			}
		})

	})

}
