package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
	"github.com/aws/aws-lambda-go/events"
)

var testingHandler = RequestHandler{}

func GetRequest(jsonStr string, obj interface{}, t *testing.T) string {
	response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: jsonStr})
	if err != nil {
		t.Fatalf("Expected 3 quotes but got an error: %+v", err)
	}
	json.Unmarshal([]byte(response.Body), &obj)
	return response.Body
}

func TestHandler(t *testing.T) {

	t.Run("Time Test for searching by string", func(t *testing.T) {
		// maxTime := 25

		// t.Run("easy search should return list of quotes with Muhammad Ali as first author", func(t *testing.T) {
		// 	start := time.Now()
		// 	var jsonStr = fmt.Sprintf(`{"searchString":"%s" }`, "Float like a butterfly sting like a bee")
		// 	var respQuotes []structs.QuoteAPIModel
		// 	GetRequest(jsonStr, &respQuotes, t)
		// 	end := time.Now()
		// 	duration := end.Sub(start)
		// 	if duration.Milliseconds() > int64(maxTime) {
		// 		t.Fatalf("Expected search for author to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
		// 	}
		// })

	})

	t.Run("Search by quotes string", func(t *testing.T) {

		t.Run("Should return the possible English topics as a list of objects", func(t *testing.T) {
			var language string = "english"
			var jsonStr = fmt.Sprintf(`{"language": "%s" }`, language)
			var respTopics []structs.TopicAPIModel
			GetRequest(jsonStr, &respTopics, t)

			if len(respTopics) <= 100 {
				t.Fatalf("got %d number of topics, expected more than %d", len(respTopics), 100)
			}
		})

		t.Run("Should return the possible Icelandic topics as a list of objects", func(t *testing.T) {

			var language string = "icelandic"
			var jsonStr = fmt.Sprintf(`{"language": "%s" }`, language)
			var respTopics []structs.TopicAPIModel
			GetRequest(jsonStr, &respTopics, t)

			if (len(respTopics) <= 5) || len(respTopics) >= 10 {
				t.Fatalf("got %d number of topics, expected at least %d and at most %d", len(respTopics), 5, 20)
			}
		})

		t.Run("Should return all possible topics as a list of objects", func(t *testing.T) {
			var respTopics []structs.TopicAPIModel
			GetRequest(`{ }`, &respTopics, t)

			if len(respTopics) != 131 {
				t.Fatalf("got %d number of topics, expected at least %d", len(respTopics), 120)
			}
		})

	})
}
