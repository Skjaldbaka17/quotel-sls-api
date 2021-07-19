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
func Setup(handler *RequestHandler, t *testing.T) (structs.QodDBModel, structs.QodDBModel, structs.QodDBModel) {
	handler.InitializeDB()

	var englishQuote structs.QuoteDBModel
	var icelandicQuote structs.QuoteDBModel
	var topicQuote structs.QuoteDBModel
	err := handler.Db.Table("quotes").Not("is_icelandic").Limit(1).Find(&englishQuote).Error
	if err != nil {
		t.Fatalf("Setup error: %s", err)
	}
	err = handler.Db.Table("quotes").Where("is_icelandic").Limit(1).Find(&icelandicQuote).Error
	if err != nil {
		t.Fatalf("Setup error: %s", err)
	}
	err = handler.Db.Table("topicsview").Where("random() < 0.001").Limit(1).Find(&topicQuote).Error
	if err != nil {
		t.Fatalf("Setup error: %s", err)
	}

	year, month, day := time.Now().Date()
	today := fmt.Sprintf("%d-%d-%d", year, month, day)
	QOD := englishQuote.ConvertToQODDBModel(today)
	QODICE := icelandicQuote.ConvertToQODDBModel(today)
	QODtopic := topicQuote.ConvertToQODDBModel(today)
	createQODs := []structs.QodDBModel{
		QOD,
		QODICE,
		QODtopic,
	}
	err = handler.Db.Table("qods").Create(&createQODs).Error
	if err != nil {
		t.Fatalf("Setup error 2: %s", err)
	}

	//CleanUp
	t.Cleanup(func() {
		handler.Db.Exec("delete from qods")
	})

	return QOD, QODICE, QODtopic
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
	QOD, QODICE, QODtopic := Setup(&testingHandler, t)

	t.Run("Time Test for getting qod", func(t *testing.T) {
		maxTime := 50
		t.Run("Time: Should get QOD", func(t *testing.T) {
			start := time.Now()
			//Get History:
			jsonStr := fmt.Sprintf(`{"language":"%s"}`, "english")
			GetRequest(jsonStr, nil, t)
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(maxTime) {
				t.Fatalf("Expected getting history of quotes to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}
		})
	})

	t.Run("Get quotes", func(t *testing.T) {
		t.Run("Should get Quote of the day", func(t *testing.T) {
			var quote structs.QodAPIModel
			GetRequest("{}", &quote, t)

			if quote.QuoteId != QOD.QuoteId {
				t.Fatalf("Expected the quote for today that setup just inserted, i.e. id %d but got quote with id %d", QOD.QuoteId, quote.QuoteId)
			}

		})

		t.Run("Should get English Quote of the day", func(t *testing.T) {
			var jsonStr = fmt.Sprintf(`{"language":"%s"}`, "english")
			var quote structs.QodAPIModel
			GetRequest(jsonStr, &quote, t)

			if quote.QuoteId != QOD.QuoteId {
				t.Fatalf("Expected the quote for today that setup just inserted, i.e. id %d but got quote with id %d", QOD.QuoteId, quote.QuoteId)
			}

		})

		t.Run("Should get Icelandic Quote of the day", func(t *testing.T) {
			var jsonStr = fmt.Sprintf(`{"language":"%s"}`, "icelandic")
			var quote structs.QodAPIModel
			GetRequest(jsonStr, &quote, t)

			if quote.QuoteId != QODICE.QuoteId {
				t.Fatalf("Expected the icelandic quote for today that setup just inserted, i.e. id %d but got quote with id %d", QODICE.QuoteId, quote.QuoteId)
			}
		})

		t.Run("Should get quote of the day for topic "+QODtopic.TopicName, func(t *testing.T) {
			var jsonStr = fmt.Sprintf(`{"topicId":%d}`, QODtopic.TopicId)
			var quote structs.QodAPIModel
			GetRequest(jsonStr, &quote, t)
			if quote.QuoteId != QODtopic.QuoteId {
				t.Fatalf("Expected the quote for today for topic %s that setup just inserted, i.e. id %d but got quote with id %d", QODtopic.TopicName, QOD.QuoteId, quote.QuoteId)
			}

			if quote.TopicId != QODtopic.TopicId {
				t.Fatalf("Expected the topicId for topic %s that setup just inserted, i.e. id %d but got quote with topicid %d", QODtopic.TopicName, QOD.QuoteId, quote.QuoteId)
			}

		})

	})

}
