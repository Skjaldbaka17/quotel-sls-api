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

func GetRequest(jsonStr string, obj interface{}, t *testing.T) string {
	response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: jsonStr})
	if err != nil {
		t.Fatalf("Expected 3 quotes but got an error: %+v", err)
	}
	json.Unmarshal([]byte(response.Body), &obj)
	return response.Body
}
func TestHandler(t *testing.T) {
	t.Cleanup(func() {
		testingHandler.Db.Table("authors").Update("count", 0)
		testingHandler.Db.Table("topics").Update("count", 0)
	})

	t.Run("Time Test for getting quotes", func(t *testing.T) {
		// maxTime := 25
		longerTime := 50
		longTime := 100

		t.Run("easy search should return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {
			start := time.Now()
			var respAuthors []structs.AuthorAPIModel
			GetRequest(`{"searchString": "Friedrich Nietzsche"}`, &respAuthors, t)
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(longerTime) {
				t.Fatalf("Expected search for author to take less than %dms but it took %dms", longerTime, duration.Milliseconds())
			}
		})

		t.Run("intermediate search should Return list of quotes with Joseph Stalin as first author", func(t *testing.T) {
			start := time.Now()
			var respAuthors []structs.AuthorAPIModel
			GetRequest(`{"searchString": "Stalin jseph"}`, &respAuthors, t)
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(longTime) {
				t.Fatalf("Expected search for author to take less than %dms but it took %dms", longTime, duration.Milliseconds())
			}
		})

		t.Run("hard search should return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {
			start := time.Now()
			var respAuthors []structs.AuthorAPIModel
			GetRequest(`{"searchString": "Niet Friedric"}`, &respAuthors, t)
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(longTime) {
				t.Fatalf("Expected search for author to take less than %dms but it took %dms", longTime, duration.Milliseconds())
			}
		})

	})

	t.Run("Search authors", func(t *testing.T) {

		t.Run("easy search should return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {

			var respAuthors []structs.AuthorAPIModel
			GetRequest(`{"searchString": "Friedrich Nietzsche"}`, &respAuthors, t)
			firstAuthor := respAuthors[0].Name
			want := "Friedrich Nietzsche"
			if firstAuthor != want {
				t.Fatalf("got %q, want %q", firstAuthor, want)
			}
		})

		t.Run("intermediate search should Return list of quotes with Joseph Stalin as first author", func(t *testing.T) {
			var respAuthors []structs.AuthorAPIModel
			GetRequest(`{"searchString": "Stalin jseph"}`, &respAuthors, t)
			firstAuthor := respAuthors[0].Name
			want := "Joseph Stalin"
			if firstAuthor != want {
				t.Fatalf("got %q, want %q", firstAuthor, want)
			}
		})

		t.Run("hard search should return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {
			var respAuthors []structs.AuthorAPIModel
			GetRequest(`{"searchString": "Niet Friedric"}`, &respAuthors, t)
			firstAuthor := respAuthors[0].Name
			want := "Friedrich Nietzsche"
			if firstAuthor != want {
				t.Fatalf("got %q, want %q", firstAuthor, want)
			}
		})

		t.Run("Search Authors By string pagination", func(t *testing.T) {
			searchString := "Martin"
			pageSize := 100
			var respAuthors []structs.AuthorAPIModel
			jsonStr := fmt.Sprintf(`{"searchString": "%s", "pageSize":%d}`, searchString, pageSize)
			GetRequest(jsonStr, &respAuthors, t)
			obj26 := respAuthors[25]

			//Next request to check if same dude in position 0 given that pageSize is 25 and same search parameters
			pageSize = 25
			jsonStr = fmt.Sprintf(`{"searchString": "%s", "pageSize":%d, "page":1}`, searchString, pageSize)
			GetRequest(jsonStr, &respAuthors, t)

			if pageSize != len(respAuthors) {
				t.Fatalf("got list of length %d but expected %d", len(respAuthors), pageSize)
			}

			if respAuthors[0].AuthorId != obj26.AuthorId {
				t.Fatalf("got %+v, want %+v", respAuthors[0], obj26)
			}
		})

	})
}
