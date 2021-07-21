package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
	"github.com/aws/aws-lambda-go/events"
)

var testingHandler = RequestHandler{}

func Setup(t *testing.T) []structs.TopicDBModel {
	testingHandler.InitializeDB()
	var topics []structs.TopicDBModel
	var topic structs.TopicDBModel
	err := testingHandler.Db.Table("topics").Where("name = ?", "Motivational").First(&topic).Error
	if err != nil {
		t.Fatalf("got error in setup motivational: %s", err)
	}
	topics = append(topics, topic)

	topic = structs.TopicDBModel{}
	err = testingHandler.Db.Table("topics").Where("name = ?", "Inspirational").First(&topic).Error
	if err != nil {
		t.Fatalf("got error in setup happiness: %s", err)
	}

	topics = append(topics, topic)

	t.Cleanup(func() {
		testingHandler.Db.Table("authors").Update("count", 0)
		testingHandler.Db.Exec("update topics set count = 0 where count > 0")
	})
	return topics
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

	topics := Setup(t)

	t.Run("Time Test for searching by string", func(t *testing.T) {
		maxTime := 25
		longTime := 200

		t.Run("Should return the first 25 quotes from a topic 'nameOfTopic'", func(t *testing.T) {
			start := time.Now()
			var jsonStr = fmt.Sprintf(`{"topic": "%s", "pageSize":25}`, strings.ToLower(topics[1].Name))
			var respQuotes []structs.QuoteAPIModel
			GetRequest(jsonStr, &respQuotes, t)

			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(longTime) {
				t.Fatalf("Expected search for author to take less than %dms but it took %dms", longTime, duration.Milliseconds())
			}

		})

		t.Run("Should return the first 25 quotes from a topic with id", func(t *testing.T) {
			start := time.Now()
			var jsonStr = fmt.Sprintf(`{"topicId": %d, "pageSize":25}`, topics[0].Id)
			var respQuotes []structs.QuoteAPIModel
			GetRequest(jsonStr, &respQuotes, t)
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(maxTime) {
				t.Fatalf("Expected search for author to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}

		})

	})

	t.Run("Search by quotes string", func(t *testing.T) {
		t.Run("Should return the first 25 quotes from a topic 'nameOfTopic'", func(t *testing.T) {
			inspirational := topics[1]
			pageSize := 25
			var jsonStr = fmt.Sprintf(`{"topic": "%s", "pageSize":%d}`, strings.ToLower(inspirational.Name), pageSize)
			var respQuotes []structs.QuoteAPIModel
			GetRequest(jsonStr, &respQuotes, t)

			if len(respQuotes) != pageSize {
				t.Fatalf("got %d number of quotes, expected %d as the pagesize", len(respQuotes), pageSize)
			}

			for _, obj := range respQuotes {
				if respQuotes[0].TopicName != inspirational.Name {
					t.Fatalf("got %+v but expected a quote with topic %s", obj, inspirational.Name)
				}
			}

		})

		t.Run("Should return the first 25 quotes from a topic with id", func(t *testing.T) {

			motivational := topics[0]
			pageSize := 25
			var jsonStr = fmt.Sprintf(`{"topicId": %d, "pageSize":%d}`, motivational.Id, pageSize)
			var respQuotes []structs.QuoteAPIModel
			GetRequest(jsonStr, &respQuotes, t)

			if len(respQuotes) != pageSize {
				t.Fatalf("got %d number of quotes, expected %d as the pagesize", len(respQuotes), pageSize)
			}

			for _, obj := range respQuotes {
				if respQuotes[0].TopicId != motivational.Id {
					t.Fatalf("got %+v but expected a quote with topicId %d", obj, motivational.Id)
				}
			}

		})

		t.Run("Should test pagination for a specific topic, by id", func(t *testing.T) {

			motivational := topics[0]
			pageSize := 25
			page := 1
			var jsonStr = fmt.Sprintf(`{"topicId": %d, "pageSize":%d, "page":%d}`, motivational.Id, pageSize, page)
			var respQuotes []structs.QuoteAPIModel
			GetRequest(jsonStr, &respQuotes, t)

			obj26 := respQuotes[0]

			// Then get the first 100 quotes, i.e. first page with pagesize 100
			pageSize = 100
			page = 0
			jsonStr = fmt.Sprintf(`{"topicId": %d, "pageSize":%d, "page":%d}`, motivational.Id, pageSize, page)
			GetRequest(jsonStr, &respQuotes, t)

			if len(respQuotes) != pageSize {
				t.Fatalf("got %d number of quotes, expected %d as the pagesize", len(respQuotes), pageSize)
			}

			//Compare the 26th object from the 100pagesize request with the 1st object from the 2nd page where pagesize is 25.
			if respQuotes[25] != obj26 {
				t.Fatalf("got %+v but expected %+v", respQuotes[25], obj26)
			}

		})

	})
}
