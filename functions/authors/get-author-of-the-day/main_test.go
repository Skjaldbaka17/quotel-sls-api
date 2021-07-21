package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
	"github.com/aws/aws-lambda-go/events"
)

var testingHandler = RequestHandler{}

//Returns AODs and AODICEs, in that order, put into the DB
func Setup(handler *RequestHandler, t *testing.T) (structs.AodDBModel, structs.AodDBModel) {
	handler.InitializeDB()

	var englishAuthor structs.AuthorDBModel
	var icelandicAuthor structs.AuthorDBModel
	err := handler.Db.Table("authors").Where("nr_of_english_quotes > 0").Limit(1).Find(&englishAuthor).Error
	if err != nil {
		t.Fatalf("Setup error: %s", err)
	}
	err = handler.Db.Table("authors").Where("nr_of_icelandic_quotes > 0").Limit(1).Find(&icelandicAuthor).Error
	if err != nil {
		t.Fatalf("Setup error: %s", err)
	}

	year, month, day := time.Now().Date()
	today := fmt.Sprintf("%d-%d-%d", year, month, day)
	AOD := englishAuthor.ConvertToAODDBModel(today, false)
	AODICE := icelandicAuthor.ConvertToAODDBModel(today, true)
	createAODs := []structs.AodDBModel{
		AOD,
		AODICE,
	}
	err = handler.Db.Table("aods").Create(&createAODs).Error
	if err != nil {
		t.Fatalf("Setup error 2: %s", err)
	}

	//CleanUp
	t.Cleanup(func() {
		handler.Db.Exec("delete from aods")
	})

	return AOD, AODICE
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
	AOD, AODICE := Setup(&testingHandler, t)
	t.Run("Time test get AOD", func(t *testing.T) {
		maxTime := int64(50)
		t.Run("Should get aod in less than 50ms", func(t *testing.T) {
			start := time.Now()
			//Get History:
			jsonStr := fmt.Sprintf(`{"language":"%s"}`, "english")
			GetRequest(jsonStr, nil, t)

			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > maxTime {
				t.Fatalf("Expected getting the AOD to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}
		})
	})

	t.Run("Get AOD/AODICE", func(t *testing.T) {
		t.Run("Should get aod", func(t *testing.T) {
			jsonStr := fmt.Sprintf(`{"language":"%s"}`, "english")
			var aod structs.AodAPIModel
			GetRequest(jsonStr, &aod, t)

			//Check wether or not the aod returned is the aod for today
			reg := regexp.MustCompile("^[0-9]+-[0-9]+-[0-9]{2}")
			today := reg.Find([]byte(time.Now().String()))
			mustMatch := regexp.MustCompile(string(today))
			if !mustMatch.MatchString(aod.Date) {
				t.Fatalf("Expected aod for today: %s, but got for date: %s", today, aod.Date)
			}

			if aod.AuthorId != AOD.AuthorId {
				t.Fatalf("Expected a valid AOD with id %d but got author with id: %d", AOD.AuthorId, aod.AuthorId)
			}
		})

		t.Run("Should get aodice", func(t *testing.T) {
			jsonStr := []byte(fmt.Sprintf(`{"language":"%s"}`, "icelandic"))
			var aod structs.AodAPIModel
			GetRequest(string(jsonStr), &aod, t)

			//Check wether or not the aodice returned is the aodice for today
			reg := regexp.MustCompile("^[0-9]+-[0-9]+-[0-9]{2}")
			today := reg.Find([]byte(time.Now().String()))
			mustMatch := regexp.MustCompile(string(today))
			if !mustMatch.MatchString(aod.Date) {
				t.Fatalf("Expected aodice for today: %s, but got for date: %s", today, aod.Date)
			}

			if aod.AuthorId != AODICE.AuthorId {
				t.Fatalf("Expected a valid AODICE with ID %d but got author with id: %d", AODICE.AuthorId, aod.AuthorId)
			}
		})
	})

}
