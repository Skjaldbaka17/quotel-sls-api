package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
	"github.com/aws/aws-lambda-go/events"
)

//Returns AODs and AODICEs, in that order, put into the DB
func Setup(handler *RequestHandler, t *testing.T) {
	handler.InitializeDB()
	//CleanUp
	t.Cleanup(func() {
		handler.Db.Unscoped().Exec("delete from aods")
		handler.Db.Unscoped().Exec("delete from aodices")
	})
}
func TestHandler(t *testing.T) {
	var testingHandler = RequestHandler{}
	Setup(&testingHandler, t)
	t.Run("Time test get AOD", func(t *testing.T) {
		t.Run("Should get aod in less than 50ms", func(t *testing.T) {
			start := time.Now()
			//Get History:
			jsonStr := []byte(fmt.Sprintf(`{"language":"%s"}`, "english"))
			_, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected the AOD but got an error: %+v", err)
			}

			end := time.Now()
			duration := end.Sub(start)
			maxTime := int64(50)
			if duration.Milliseconds() > maxTime {
				t.Fatalf("Expected getting the AOD to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}
		})
	})

	t.Run("Get AOD/AODICE", func(t *testing.T) {
		t.Run("Should get aod", func(t *testing.T) {
			jsonStr := []byte(fmt.Sprintf(`{"language":"%s"}`, "english"))
			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected the AOD but got an error: %+v", err)
			}

			var aod structs.AodAPIModel
			json.Unmarshal([]byte(response.Body), &aod)

			//Check wether or not the aod returned is the aod for today
			reg := regexp.MustCompile("^[0-9]+-[0-9]+-[0-9]{2}")
			today := reg.Find([]byte(time.Now().String()))
			mustMatch := regexp.MustCompile(string(today))
			if !mustMatch.MatchString(aod.Date) {
				t.Fatalf("Expected aod for today: %s, but got for date: %s", today, aod.Date)
			}

			if aod.Id <= 0 {
				t.Fatalf("Expected a valid AOD but got author with id: %d", aod.Id)
			}

			if aod.Name == "" {
				t.Fatalf("Expected an AOD with a valid name but got author with name: %s, and id: %d", aod.Name, aod.Id)
			}

			if aod.Nationality == "" {
				t.Fatalf("Expected an AOD with a valid nationality but got author with nationality: %s, and id: %d", aod.Nationality, aod.Id)
			}
			log.Println(aod)
		})

		t.Run("Should get aodice", func(t *testing.T) {
			jsonStr := []byte(fmt.Sprintf(`{"language":"%s"}`, "icelandic"))
			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected the AODICE but got an error: %+v", err)
			}

			var aod structs.AodAPIModel
			json.Unmarshal([]byte(response.Body), &aod)

			//Check wether or not the aodice returned is the aodice for today
			reg := regexp.MustCompile("^[0-9]+-[0-9]+-[0-9]{2}")
			today := reg.Find([]byte(time.Now().String()))
			mustMatch := regexp.MustCompile(string(today))
			if !mustMatch.MatchString(aod.Date) {
				t.Fatalf("Expected aodice for today: %s, but got for date: %s", today, aod.Date)
			}

			if aod.Id <= 0 {
				t.Fatalf("Expected a valid AODICE but got author with id: %d", aod.Id)
			}

			if aod.Name == "" {
				t.Fatalf("Expected an AODICE with a valid name but got author with name: %s, and id: %d", aod.Name, aod.Id)
			}
		})
	})

}
