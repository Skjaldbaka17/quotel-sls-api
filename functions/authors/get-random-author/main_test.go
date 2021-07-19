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

	t.Run("Time Test for getting authors by ids", func(t *testing.T) {
		maxTime := 60

		t.Run("Should return a random author with only a single quote (i.e. default)", func(t *testing.T) {
			start := time.Now()
			GetRequest("{}", nil, t)
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(maxTime/2) {
				t.Fatalf("Expected getting authors by ids to take less than %dms but it took %dms", maxTime/2, duration.Milliseconds())
			}
		})

		t.Run("Should return a random Author with only quotes from him in English", func(t *testing.T) {
			start := time.Now()
			language := "english"
			var jsonStr = fmt.Sprintf(`{"language":"%s"}`, language)
			GetRequest(jsonStr, nil, t)
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(maxTime) {
				t.Fatalf("Expected getting authors by ids to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}
		})

	})
	t.Run("Get random author", func(t *testing.T) {

		t.Run("Should return a random author with only a single quote (i.e. default)", func(t *testing.T) {

			var authors []structs.QuoteAPIModel
			GetRequest("{}", &authors, t)

			if len(authors) != 1 {
				t.Fatalf("Expected only a single quote from the random author but got %d", len(authors))
			}
			firstAuthor := authors[0]
			if firstAuthor.Id == 0 {
				t.Fatal("got an author with id 0, want author with valid id")
			}

			GetRequest("{}", &authors, t)
			if firstAuthor.Id == authors[0].Id {
				t.Fatalf("Expected two different authors but got the same author twice which is higly improbable, got author with id %d and name %s", firstAuthor.Id, firstAuthor.Name)
			}

		})

		t.Run("Should return a random Author with only quotes from him in Icelandic", func(t *testing.T) {
			language := "icelandic"
			var jsonStr = fmt.Sprintf(`{"language":"%s"}`, language)
			var authors []structs.QuoteAPIModel
			GetRequest(jsonStr, &authors, t)
			firstAuthor := authors[0]
			if firstAuthor.Name == "" {
				t.Fatalf("Expected a random firstauthor but got an empty name for author")
			}

			if !firstAuthor.IsIcelandic {
				t.Fatalf("Expected the quotes returned to be in icelandic")
			}

			GetRequest(jsonStr, &authors, t)
			secondAuthor := authors[0]
			if firstAuthor.Id == secondAuthor.Id {
				t.Fatalf("Expected two different authors but got the same author twice which is higly improbable, got author with id %d and name %s", firstAuthor.Id, firstAuthor.Name)
			}
		})

		t.Run("Should return a random Author with only quotes from him in English", func(t *testing.T) {

			language := "english"
			var jsonStr = fmt.Sprintf(`{"language":"%s"}`, language)
			var authors []structs.QuoteAPIModel
			GetRequest(jsonStr, &authors, t)
			firstAuthor := authors[0]
			if firstAuthor.Name == "" {
				t.Fatalf("Expected a random firstauthor but got an empty name for author")
			}

			if firstAuthor.IsIcelandic {
				t.Fatalf("Expected the quotes returned to not be in icelandic")
			}

			GetRequest(jsonStr, &authors, t)
			secondAuthor := authors[0]
			if firstAuthor.Id == secondAuthor.Id {
				t.Fatalf("Expected two different authors but got the same author twice which is higly improbable, got author with id %d and name %s", firstAuthor.Id, firstAuthor.Name)
			}

		})

		t.Run("Should return author with a maximum of 2 of his quotes", func(t *testing.T) {
			maxQuotes := 2
			var jsonStr = fmt.Sprintf(`{"maxQuotes":%d}`, maxQuotes)
			var authors []structs.QuoteAPIModel
			GetRequest(jsonStr, &authors, t)

			firstAuthor := authors[0]
			if firstAuthor.Name == "" {
				t.Fatalf("Expected a random author but got an empty name for author")
			}

			if len(authors) != 2 {
				t.Fatalf("Expected 2 quotes but got %d", len(authors))
			}
		})

	})

}
