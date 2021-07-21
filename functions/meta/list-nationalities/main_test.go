package main

import (
	"encoding/json"
	"testing"
	"time"

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
	type Response struct {
		Nationalities []string
	}

	t.Run("Time Test for getting quotes", func(t *testing.T) {
		// maxTime := 20
		longTime := 100
		t.Run("Get all nationalities", func(t *testing.T) {
			start := time.Now()

			var response Response
			GetRequest(`{}`, &response, t)
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(longTime) {
				t.Fatalf("Expected search for author to take less than %dms but it took %dms", longTime, duration.Milliseconds())
			}
		})

	})

	t.Run("Get all nationalities", func(t *testing.T) {
		var response Response
		GetRequest(`{}`, &response, t)
		if len(response.Nationalities) != 99 {
			t.Fatalf("expected %d nationalities but got %d", 99, len(response.Nationalities))
		}
	})

}
