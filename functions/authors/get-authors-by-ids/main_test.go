package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
	"github.com/aws/aws-lambda-go/events"
)

var testingHandler = RequestHandler{}

//Returns AODs and AODICEs, in that order, put into the DB
func Setup(handler *RequestHandler, t *testing.T) []structs.AuthorDBModel {
	handler.InitializeDB()
	var authors []structs.AuthorDBModel
	err := handler.Db.Table("authors").Limit(3).Find(&authors).Error
	if err != nil {
		t.Fatalf("Setup error: %s", err)
	}
	//CleanUp
	t.Cleanup(func() {
		handler.Db.Table("authors").Model(&authors).Update("count = ?", 0)
	})
	return authors
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
	authors := Setup(&testingHandler, t)

	t.Run("Time Test for getting authors by ids", func(t *testing.T) {
		maxTime := 20
		t.Run(fmt.Sprintf("should return Authors with ids %d, %d and %d withing 20ms", authors[0].ID, authors[1].ID, authors[2].ID), func(t *testing.T) {
			start := time.Now()
			var jsonStr = fmt.Sprintf(`{"ids": [%d,%d,%d]}`, authors[0].ID, authors[1].ID, authors[2].ID)
			GetRequest(jsonStr, nil, t)
			end := time.Now()
			duration := end.Sub(start)

			if duration.Milliseconds() > int64(maxTime) {
				t.Fatalf("Expected getting authors by ids to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}
		})
	})
	t.Run("Get authors", func(t *testing.T) {
		t.Run("should return Author with id "+strconv.Itoa(int(authors[0].ID)), func(t *testing.T) {
			var jsonStr = fmt.Sprintf(`{"ids": [%d]}`, authors[0].ID)
			var respAuthors []structs.AuthorAPIModel
			GetRequest(jsonStr, &respAuthors, t)

			if len(respAuthors) > 1 {
				t.Fatalf("Expected %d authors, but got %d authors", 1, len(respAuthors))
			}
			firstAuthor := respAuthors[0]
			if firstAuthor.Id != authors[0].ID {
				t.Fatalf("got %d, want %d", firstAuthor.Id, authors[0].ID)
			}
		})

		t.Run(fmt.Sprintf("should return Authors with ids %d, %d and %d", authors[0].ID, authors[1].ID, authors[2].ID), func(t *testing.T) {
			var jsonStr = fmt.Sprintf(`{"ids": [%d,%d,%d]}`, authors[0].ID, authors[1].ID, authors[2].ID)
			var respAuthors []structs.AuthorAPIModel
			GetRequest(jsonStr, &respAuthors, t)
			if len(respAuthors) > 3 {
				t.Fatalf("Expected %d authors, but got %d authors", 1, len(respAuthors))
			}

			//Check that all the authors we wanted are returned
			for _, author := range authors {
				for i, re := range respAuthors {
					if re.Id == author.ID {
						break
					}
					if i == len(respAuthors)-1 {
						t.Fatalf("got %+v, want %d", respAuthors, author.ID)
					}
				}
			}
		})

	})

}
