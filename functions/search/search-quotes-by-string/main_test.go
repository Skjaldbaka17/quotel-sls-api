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

		t.Run("easy search should return list of quotes with Muhammad Ali as first author", func(t *testing.T) {
			start := time.Now()
			var jsonStr = fmt.Sprintf(`{"searchString":"%s" }`, "Float like a butterfly sting like a bee")
			var respQuotes []structs.QuoteAPIModel
			GetRequest(jsonStr, &respQuotes, t)
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(maxTime) {
				t.Fatalf("Expected search for author to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}
		})

		t.Run("intermediate search should return list of quotes with Muhammad Ali as first author", func(t *testing.T) {
			start := time.Now()
			var jsonStr = fmt.Sprintf(`{"searchString":"%s" }`, "bee sting like a butterfly")
			var respQuotes []structs.QuoteAPIModel
			GetRequest(jsonStr, &respQuotes, t)
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(maxTime) {
				t.Fatalf("Expected search for author to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}
		})

		t.Run("hard search should return list of quotes with Muhammad Ali as first author", func(t *testing.T) {
			start := time.Now()
			var jsonStr = fmt.Sprintf(`{"searchString":"%s" }`, "bee butterfly float")
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

		t.Run("easy search should return list of quotes with Muhammad Ali as first author", func(t *testing.T) {
			searchString := "Float like a butterfly sting like a bee"
			var jsonStr = fmt.Sprintf(`{"searchString":"%s" }`, searchString)
			var respQuotes []structs.QuoteAPIModel
			GetRequest(jsonStr, &respQuotes, t)
			firstAuthor := respQuotes[0].Name
			want := "Muhammad Ali"
			if firstAuthor != want {
				t.Fatalf("got %q, want %q", firstAuthor, want)
			}
		})

		t.Run("intermediate search should return list of quotes with Muhammad Ali as first author", func(t *testing.T) {
			searchString := "bee sting like a butterfly"
			var jsonStr = fmt.Sprintf(`{"searchString":"%s" }`, searchString)
			var respQuotes []structs.QuoteAPIModel
			GetRequest(jsonStr, &respQuotes, t)
			firstAuthor := respQuotes[0].Name
			want := "Muhammad Ali"
			if firstAuthor != want {
				t.Fatalf("got %q, want %q", firstAuthor, want)
			}
		})

		t.Run("hard search should return list of quotes with Muhammad Ali as first author", func(t *testing.T) {
			searchString := "bee butterfly float"
			var jsonStr = fmt.Sprintf(`{"searchString":"%s" }`, searchString)
			var respQuotes []structs.QuoteAPIModel
			GetRequest(jsonStr, &respQuotes, t)
			firstAuthor := respQuotes[0].Name
			want := "Muhammad Ali"
			if firstAuthor != want {
				t.Fatalf("got %q, want %q", firstAuthor, want)
			}
		})

		t.Run("Search for quote 'Happiness resides not in possessions...' inside topic 'inspirational' by supplying its topicid", func(t *testing.T) {
			searchString := "Happiness resides not in possessions"
			inspirational := topics[1]
			var jsonStr = fmt.Sprintf(`{"searchString":"%s","topicIds":[%d] }`, searchString, inspirational.Id)
			var respQuotes []structs.QuoteAPIModel
			GetRequest(jsonStr, &respQuotes, t)
			firstAuthorName := respQuotes[0].Name
			want_author := "Democritus"
			if firstAuthorName != want_author {
				t.Fatalf("got %q, want %q", firstAuthorName, want_author)
			}

			firstAuthorQuote := respQuotes[0].Quote
			want_quote := "Happiness resides not in possessions, and not in gold, happiness dwells in the soul."
			if firstAuthorQuote != want_quote {
				t.Fatalf("got %q, want %q", firstAuthorQuote, want_quote)
			}

			if respQuotes[0].TopicId != inspirational.Id {
				t.Fatalf("got quote with topicId %d, but expected with topicID %d. Quote got: %+v", respQuotes[0].TopicId, inspirational.Id, respQuotes[0])
			}
		})

		t.Run("Search Quotes By string pagination", func(t *testing.T) {

			searchString := "Hate"
			pageSize := 50
			var jsonStr = fmt.Sprintf(`{"searchString":"%s" ,"pageSize":%d}`, searchString, pageSize)
			var respQuotes []structs.QuoteAPIModel
			GetRequest(jsonStr, &respQuotes, t)
			obj26 := respQuotes[25]

			//Next request to check if same dude in position 0 given that pageSize is 25 and same search parameters
			pageSize = 25
			jsonStr = fmt.Sprintf(`{"searchString": "%s", "pageSize":%d, "page":1}`, searchString, pageSize)
			GetRequest(jsonStr, &respQuotes, t)

			if pageSize != len(respQuotes) {
				t.Fatalf("got list of length %d but expected %d", len(respQuotes), pageSize)
			}

			if respQuotes[0].QuoteId != obj26.QuoteId {
				t.Fatalf("got %+v, want %+v", respQuotes[0], obj26)
			}
		})

	})
}
