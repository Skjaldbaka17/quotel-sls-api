package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
	"github.com/aws/aws-lambda-go/events"
)

func Setup(handler *RequestHandler, t *testing.T) ([]structs.AuthorDBModel, []structs.TopicDBModel) {
	handler.InitializeDB()
	var authors []structs.AuthorDBModel
	var author structs.AuthorDBModel
	err := handler.Db.Table("authors").Where("name = ?", "Theodore Roosevelt").First(&author).Error
	if err != nil {
		t.Fatalf("got error in setup: %s", err)
	}
	authors = append(authors, author)

	var topics []structs.TopicDBModel
	var topic structs.TopicDBModel
	err = handler.Db.Table("topics").Where("name = ?", "Motivational").First(&topic).Error
	if err != nil {
		t.Fatalf("got error in setup motivational: %s", err)
	}

	topics = append(topics, topic)
	topic = structs.TopicDBModel{}
	err = handler.Db.Table("topics").Where("name = ?", "Smile").Find(&topic).Error
	if err != nil {
		t.Fatalf("got error in setup smile: %s", err)
	}

	topics = append(topics, topic)
	topic = structs.TopicDBModel{}
	err = handler.Db.Table("topics").Where("name = ?", "Happiness").First(&topic).Error
	if err != nil {
		t.Fatalf("got error in setup happiness: %s", err)
	}

	topics = append(topics, topic)

	topic = structs.TopicDBModel{}
	err = handler.Db.Table("topics").Where("name = ?", "Inspirational").First(&topic).Error
	if err != nil {
		t.Fatalf("got error in setup happiness: %s", err)
	}

	topics = append(topics, topic)

	t.Cleanup(func() {
		handler.Db.Table("authors").Model(&authors).Update("count", 0)
		handler.Db.Table("topics").Model(&authors).Update("count", 0)
	})
	return authors, topics
}

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

	authors, topics := Setup(&testingHandler, t)
	t.Run("Time Test for getting quotes", func(t *testing.T) {
		maxTime := 25
		longTime := 250

		t.Run("Should return a random quote", func(t *testing.T) {
			start := time.Now()
			var firstRespQuote structs.QuoteAPIModel
			GetRequest("{}", &firstRespQuote, t)
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(maxTime) {
				t.Fatalf("Expected getting random quote to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}
		})

		t.Run("Should return a random quote from Teddy Roosevelt (given authorId)", func(t *testing.T) {
			start := time.Now()
			teddy := authors[0]
			var jsonStr = fmt.Sprintf(`{"authorId": %d}`, teddy.ID)
			var firstRespQuote structs.QuoteAPIModel
			GetRequest(jsonStr, &firstRespQuote, t)
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(maxTime) {
				t.Fatalf("Expected getting random quote to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}
		})

		t.Run("Should return a random quote from topic 'motivational','smile' or 'happiness' (given topicId)", func(t *testing.T) {
			start := time.Now()
			motivational := topics[0]
			smile := topics[1]
			happiness := topics[2]
			var jsonStr = fmt.Sprintf(`{"topicIds": [%d,%d,%d]}`, motivational.Id, smile.Id, happiness.Id)
			var firstRespQuote structs.QuoteAPIModel
			GetRequest(jsonStr, &firstRespQuote, t)
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(maxTime) {
				t.Fatalf("Expected getting random quote to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}
		})

		t.Run("Should return a random Icelandic quote", func(t *testing.T) {
			start := time.Now()
			language := "Icelandic"
			var jsonStr = fmt.Sprintf(`{"language": "%s"}`, language)
			var firstRespQuote structs.QuoteAPIModel
			GetRequest(jsonStr, &firstRespQuote, t)
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(longTime) {
				t.Fatalf("Expected getting random quote to take less than %dms but it took %dms", longTime, duration.Milliseconds())
			}
		})

		t.Run("Should return a random quote containing the searchString 'love'", func(t *testing.T) {
			start := time.Now()
			searchString := "love"
			var jsonStr = fmt.Sprintf(`{"searchString":"%s"}`, searchString)
			var firstRespQuote structs.QuoteAPIModel
			GetRequest(jsonStr, &firstRespQuote, t)
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(longTime) {
				t.Fatalf("Expected getting random quote to take less than %dms but it took %dms", longTime, duration.Milliseconds())
			}

		})

		t.Run("Should return a random quote containing the searchString 'strong' from the topic 'inspirational' (given topicId)", func(t *testing.T) {
			start := time.Now()
			inspirational := topics[3]
			searchString := "strong"
			var jsonStr = fmt.Sprintf(`{"searchString":"%s","topicIds": [%d]}`, searchString, inspirational.Id)
			var firstRespQuote structs.QuoteAPIModel
			GetRequest(jsonStr, &firstRespQuote, t)
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(maxTime) {
				t.Fatalf("Expected getting random quote to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}
		})

	})

	t.Run("Get quotes", func(t *testing.T) {

		//The test calls the function twice to test if the function returns two different quotes
		t.Run("Should return a random quote", func(t *testing.T) {
			var firstRespQuote structs.QuoteAPIModel
			GetRequest("{}", &firstRespQuote, t)
			var secondRespQuote structs.QuoteAPIModel
			GetRequest("{}", &secondRespQuote, t)
			if secondRespQuote.QuoteId == firstRespQuote.QuoteId {
				t.Fatalf("Expected two different quotes but got the same quote twice which is higly improbable")
			}
		})

		t.Run("Should return a random quote from Teddy Roosevelt (given authorId)", func(t *testing.T) {

			teddy := authors[0]
			var jsonStr = fmt.Sprintf(`{"authorId": %d}`, teddy.ID)
			var firstRespQuote structs.QuoteAPIModel
			GetRequest(jsonStr, &firstRespQuote, t)

			if firstRespQuote.Id != teddy.ID {
				t.Fatalf("Expected author from %s but got author from %s", teddy.Name, firstRespQuote.Name)
			}

			var secondRespQuote structs.QuoteAPIModel
			GetRequest(jsonStr, &secondRespQuote, t)

			if secondRespQuote.Id != firstRespQuote.Id {
				t.Fatalf("got author with id %d, expected author with id %d", secondRespQuote.Id, firstRespQuote.Id)
			}

			if secondRespQuote.QuoteId == firstRespQuote.QuoteId {
				t.Fatalf("got quote %s, expected a random different quote", secondRespQuote.Quote)
			}

		})

		t.Run("Should return a random quote from topic 'motivational' (given topicId)", func(t *testing.T) {

			motivational := topics[0]
			var jsonStr = fmt.Sprintf(`{"topicIds": [%d]}`, motivational.Id)
			var firstRespQuote structs.QuoteAPIModel
			GetRequest(jsonStr, &firstRespQuote, t)
			if firstRespQuote.TopicName != motivational.Name {
				t.Fatalf("got topicname: %s, expected %s", firstRespQuote.TopicName, motivational.Name)
			}
			var secondRespQuote structs.QuoteAPIModel
			GetRequest(jsonStr, &secondRespQuote, t)
			if secondRespQuote.TopicId != firstRespQuote.TopicId {
				t.Fatalf("got topic with id %d, expected topic with id %d", secondRespQuote.TopicId, firstRespQuote.TopicId)
			}

			if secondRespQuote.QuoteId == firstRespQuote.QuoteId {
				t.Fatalf("got quote %s, expected a random different quote", secondRespQuote.Quote)
			}
		})

		t.Run("Should return a random quote from topic 'motivational','smile' or 'happiness' (given topicId)", func(t *testing.T) {

			motivational := topics[0]
			smile := topics[1]
			happiness := topics[2]
			var jsonStr = fmt.Sprintf(`{"topicIds": [%d,%d,%d]}`, motivational.Id, smile.Id, happiness.Id)
			var firstRespQuote structs.QuoteAPIModel
			GetRequest(jsonStr, &firstRespQuote, t)
			if firstRespQuote.TopicName != motivational.Name && firstRespQuote.TopicName != smile.Name && firstRespQuote.TopicName != happiness.Name {
				t.Fatalf("got topicname: %s, expected any of: %s,%s,%s ", firstRespQuote.TopicName, motivational.Name, smile.Name, happiness.Name)
			}
			var secondRespQuote structs.QuoteAPIModel
			GetRequest(jsonStr, &secondRespQuote, t)

			if secondRespQuote.QuoteId == firstRespQuote.QuoteId {
				t.Fatalf("got quote %s, expected a random different quote", secondRespQuote.Quote)
			}
		})

		t.Run("Should return a random English quote", func(t *testing.T) {

			language := "english"
			var jsonStr = fmt.Sprintf(`{"language": "%s"}`, language)
			var firstRespQuote structs.QuoteAPIModel
			GetRequest(jsonStr, &firstRespQuote, t)
			if firstRespQuote.IsIcelandic {
				t.Fatalf("first response, got an IcelandicQuote but expected an English quote")
			}
			var secondRespQuote structs.QuoteAPIModel
			GetRequest(jsonStr, &secondRespQuote, t)
			if secondRespQuote.IsIcelandic {
				t.Fatalf("second response, got an IcelandicQuote but expected an English quote")
			}

			if secondRespQuote.QuoteId == firstRespQuote.QuoteId {
				t.Fatalf("got quote %s, expected a random different quote", secondRespQuote.Quote)
			}
		})

		t.Run("Should return a random Icelandic quote", func(t *testing.T) {

			language := "Icelandic"
			var jsonStr = fmt.Sprintf(`{"language": "%s"}`, language)
			var firstRespQuote structs.QuoteAPIModel
			GetRequest(jsonStr, &firstRespQuote, t)
			if !firstRespQuote.IsIcelandic {
				t.Fatalf("first response, got an EnglishQuote but expected an Icelandic quote")
			}
			var secondRespQuote structs.QuoteAPIModel
			GetRequest(jsonStr, &secondRespQuote, t)
			if !secondRespQuote.IsIcelandic {
				t.Fatalf("second response, got an EnglishQuote, %+v, but expected an Icelandic quote", secondRespQuote)
			}

			if secondRespQuote.QuoteId == firstRespQuote.QuoteId {
				t.Fatalf("got quote %s, expected a random different quote", secondRespQuote.Quote)
			}
		})

		t.Run("Should return a random quote containing the searchString 'love'", func(t *testing.T) {

			searchString := "love"
			var jsonStr = fmt.Sprintf(`{"searchString":"%s"}`, searchString)
			var firstRespQuote structs.QuoteAPIModel
			GetRequest(jsonStr, &firstRespQuote, t)
			regexStub := searchString[:3]
			m1 := regexp.MustCompile(regexStub)
			if !m1.Match([]byte(strings.ToLower(firstRespQuote.Quote))) {
				t.Fatalf("first response, got the quote %+v that does not contain the searchString %s", firstRespQuote, regexStub)
			}

		})

		t.Run("Should return a random Icelandic quote containing the searchString 'þitt'", func(t *testing.T) {

			searchString := "þitt"
			var jsonStr = fmt.Sprintf(`{"searchString":"%s"}`, searchString)
			var firstRespQuote structs.QuoteAPIModel
			GetRequest(jsonStr, &firstRespQuote, t)
			m1 := regexp.MustCompile(searchString)
			if !m1.Match([]byte(firstRespQuote.Quote)) {
				t.Fatalf("first response, got the quote %+v that does not contain the searchString %s", firstRespQuote, searchString)
			}

			if !firstRespQuote.IsIcelandic {
				t.Fatalf("first response, got the quote %+v which is in English but expected it to be in icelandic", firstRespQuote)
			}

		})

		t.Run("Should return a random quote containing the searchString 'strong' from the topic 'inspirational' (given topicId)", func(t *testing.T) {

			inspirational := topics[3]
			searchString := "strong"
			var jsonStr = fmt.Sprintf(`{"searchString":"%s","topicIds": [%d]}`, searchString, inspirational.Id)
			var firstRespQuote structs.QuoteAPIModel
			GetRequest(jsonStr, &firstRespQuote, t)

			if firstRespQuote.TopicName != inspirational.Name {
				t.Fatalf("got %s, expected %s", firstRespQuote.TopicName, inspirational.Name)
			}

			m1 := regexp.MustCompile(searchString)
			if !m1.Match([]byte(firstRespQuote.Quote)) {
				t.Fatalf("first response, got the quote %+v that does not contain the searchString %s", firstRespQuote, searchString)
			}
		})

	})
}
