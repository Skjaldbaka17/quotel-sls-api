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
func TestHandler(t *testing.T) {
	var testingHandler = RequestHandler{}
	authors := Setup(&testingHandler, t)

	t.Run("Time Test for getting authors by ids", func(t *testing.T) {
		t.Run(fmt.Sprintf("should return Authors with ids %d, %d and %d withing 20ms", authors[0].ID, authors[1].ID, authors[2].ID), func(t *testing.T) {
			start := time.Now()
			var jsonStr = []byte(fmt.Sprintf(`{"ids": [%d,%d,%d]}`, authors[0].ID, authors[1].ID, authors[2].ID))
			_, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}
			end := time.Now()
			duration := end.Sub(start)
			maxTime := 20
			if duration.Milliseconds() > int64(maxTime) {
				t.Fatalf("Expected getting authors by ids to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}
		})
	})
	t.Run("Get authors", func(t *testing.T) {
		t.Run("should return Author with id "+strconv.Itoa(int(authors[0].ID)), func(t *testing.T) {
			var jsonStr = []byte(fmt.Sprintf(`{"ids": [%d]}`, authors[0].ID))
			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}
			var respAuthors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &respAuthors)

			if len(respAuthors) > 1 {
				t.Fatalf("Expected %d authors, but got %d authors", 1, len(respAuthors))
			}
			firstAuthor := respAuthors[0]
			if firstAuthor.Id != authors[0].ID {
				t.Fatalf("got %d, want %d", firstAuthor.Id, authors[0].ID)
			}
		})

		t.Run(fmt.Sprintf("should return Authors with ids %d, %d and %d", authors[0].ID, authors[1].ID, authors[2].ID), func(t *testing.T) {
			var jsonStr = []byte(fmt.Sprintf(`{"ids": [%d,%d,%d]}`, authors[0].ID, authors[1].ID, authors[2].ID))
			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}
			var respAuthors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &respAuthors)
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
