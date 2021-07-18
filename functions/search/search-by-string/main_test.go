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
		testingHandler.Db.Table("topics").Update("count", 0)
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
		longTime := 700

		t.Run("searching for author", func(t *testing.T) {
			t.Run("easy search should return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {
				start := time.Now()
				var jsonStr = fmt.Sprintf(`{"searchString": "%s"}`, "Friedrich Nietzsche")
				var respQuotes []structs.QuoteAPIModel
				GetRequest(jsonStr, &respQuotes, t)
				end := time.Now()
				duration := end.Sub(start)
				if duration.Milliseconds() > int64(maxTime) {
					t.Fatalf("Expected search for author to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
				}
			})

			t.Run("hard search should return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {
				start := time.Now()
				var jsonStr = fmt.Sprintf(`{"searchString": "%s"}`, "Nietshe Friedrik")
				var respQuotes []structs.QuoteAPIModel
				GetRequest(jsonStr, &respQuotes, t)
				end := time.Now()
				duration := end.Sub(start)
				if duration.Milliseconds() > int64(longTime) {
					t.Fatalf("Expected search for author to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
				}
			})
		})

		t.Run("searching for quote", func(t *testing.T) {
			t.Run("easy search should return list of quotes with Martin Luther as first author", func(t *testing.T) {
				start := time.Now()
				var jsonStr = fmt.Sprintf(`{"searchString": "%s"}`, "If you are not allowed to Laugh in Heaven")
				var respQuotes []structs.QuoteAPIModel
				GetRequest(jsonStr, &respQuotes, t)
				end := time.Now()
				duration := end.Sub(start)
				if duration.Milliseconds() > int64(maxTime) {
					t.Fatalf("Expected search for author to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
				}
			})
		})

		t.Run("General Search inside topic 'inspirational' or 'motivational' by supplying its id, should return 'Michael Jordan' Quote", func(t *testing.T) {
			start := time.Now()
			motivational := topics[0]
			inspirational := topics[1]
			var jsonStr = fmt.Sprintf(`{"searchString": "Jordan Michel", "topicIds":[%d,%d]}`, motivational.Id, inspirational.Id)
			var respQuotes []structs.QuoteAPIModel
			GetRequest(jsonStr, &respQuotes, t)
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(longTime) {
				t.Fatalf("Expected search for author to take less than %dms but it took %dms", longTime, duration.Milliseconds())
			}
		})

	})

	t.Run("Search by string", func(t *testing.T) {

		t.Run("searching for author", func(t *testing.T) {
			t.Run("easy search should return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {
				var jsonStr = fmt.Sprintf(`{"searchString": "%s"}`, "Friedrich Nietzsche")
				var respQuotes []structs.QuoteAPIModel
				GetRequest(jsonStr, &respQuotes, t)
				firstAuthor := respQuotes[0].Name //Use index 1 because in index 0 there is an author talking extensively about Nietzsche
				want := "Friedrich Nietzsche"
				if firstAuthor != want {
					t.Fatalf("got %q, want %q", firstAuthor, want)
				}
			})

			t.Run("hard search should return list of quotes with Friedrich Nietzsche as first author", func(t *testing.T) {
				var jsonStr = fmt.Sprintf(`{"searchString": "%s"}`, "Nietshe Friedrik")
				var respQuotes []structs.QuoteAPIModel
				GetRequest(jsonStr, &respQuotes, t)
				firstAuthor := respQuotes[0].Name //Use index 1 because in index 0 there is an author talking extensively about Nietzsche
				want := "Friedrich Nietzsche"
				if firstAuthor != want {
					t.Fatalf("got %q, want %q", firstAuthor, want)
				}
			})
		})

		t.Run("searching for quote", func(t *testing.T) {
			t.Run("easy search should return list of quotes with Martin Luther as first author", func(t *testing.T) {

				var jsonStr = fmt.Sprintf(`{"searchString": "%s"}`, "If you are not allowed to Laugh in Heaven")
				var respQuotes []structs.QuoteAPIModel
				GetRequest(jsonStr, &respQuotes, t)
				firstAuthor := respQuotes[0].Name
				want := "Martin Luther"
				if firstAuthor != want {
					t.Fatalf("got %q, want %q", firstAuthor, want)
				}
			})
		})

		t.Run("General Search inside topic 'inspirational' or 'motivational' by supplying its id, should return 'Michael Jordan' Quote", func(t *testing.T) {
			motivational := topics[0]
			inspirational := topics[1]
			var jsonStr = fmt.Sprintf(`{"searchString": "Jordan Michel", "topicIds":[%d,%d]}`, motivational.Id, inspirational.Id)
			var respQuotes []structs.QuoteAPIModel
			GetRequest(jsonStr, &respQuotes, t)
			firstAuthorName := respQuotes[0].Name
			want_author := "Michael Jordan"
			if firstAuthorName != want_author {
				t.Fatalf("got %q, want %q", firstAuthorName, want_author)
			}

			if respQuotes[0].TopicId != inspirational.Id && respQuotes[0].TopicId != motivational.Id {
				t.Fatalf("got quote with topicId %d, but expected with topicID either %d or %d. Quote got: %+v", respQuotes[0].TopicId, motivational.Id, inspirational.Id, respQuotes[0])
			}
		})

	})
}
