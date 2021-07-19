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

func GetRequest(jsonStr string, obj interface{}, t *testing.T) string {
	response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: jsonStr})
	if err != nil {
		t.Fatalf("Expected 3 quotes but got an error: %+v", err)
	}
	json.Unmarshal([]byte(response.Body), &obj)
	return response.Body
}

//Returns AODs and AODICEs, in that order, put into the DB
func Setup(handler *RequestHandler, t *testing.T) ([]structs.AodDBModel, []structs.AodDBModel) {
	handler.InitializeDB()

	//Get the first 3 english- and icelandic-authors to be set as AODs/AODICEs
	var englishAuthors []structs.AuthorDBModel
	var icelandicAuthors []structs.AuthorDBModel
	err := handler.Db.Table("authors").Where("nr_of_english_quotes > 0").Limit(3).Find(&englishAuthors).Error
	if err != nil {
		t.Fatalf("Setup error: %s", err)
	}
	err = handler.Db.Table("authors").Where("nr_of_icelandic_quotes > 0").Limit(3).Find(&icelandicAuthors).Error
	if err != nil {
		t.Fatalf("Setup error: %s", err)
	}

	// Insert AODs and AODICEs for 2021-06-16, 2021-06-16,2019-06-16
	dates := []string{"2021-06-16", "2020-06-16", "2019-06-16"}
	var AODs []structs.AodDBModel
	var AODICEs []structs.AodDBModel
	for idx, date := range dates {
		AODs = append(AODs, englishAuthors[idx].ConvertToAODDBModel(date, false))
		AODICEs = append(AODICEs, icelandicAuthors[idx].ConvertToAODDBModel(date, true))
	}
	createAOD := append(AODs, AODICEs...)
	err = handler.Db.Table("aods").Create(&createAOD).Error
	if err != nil {
		t.Fatalf("Setup error 2: %s", err)
	}

	//CleanUp
	t.Cleanup(func() {
		handler.Db.Exec("delete from aods")
	})

	return AODs, AODICEs
}
func TestHandler(t *testing.T) {
	AODs, AODICEs := Setup(&testingHandler, t)
	t.Run("Time test History", func(t *testing.T) {
		maxTime := 50
		t.Run("Should get history, when there is no history i.e. need to create AOD for today at least, in less than 50ms", func(t *testing.T) {
			start := time.Now()
			//Get History:
			jsonStr := []byte(fmt.Sprintf(`{"language":"%s"}`, "english"))
			GetRequest(string(jsonStr), nil, t)

			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(maxTime) {
				t.Fatalf("Expected getting history of AODS to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}
		})
	})
	t.Run("AOD/AODICE History", func(t *testing.T) {

		t.Run("Should get complete history of AODs", func(t *testing.T) {
			//Get History:
			jsonStr := []byte(fmt.Sprintf(`{"language":"%s"}`, "english"))
			responseBod := GetRequest(string(jsonStr), nil, t)

			for _, author := range AODs {
				reg := regexp.MustCompile(author.Date)
				if !reg.MatchString(responseBod) {
					t.Fatalf("Missing Aod for date: %s", author.Date)
				}
			}

		})

		t.Run("Should get complete history of AODICEs", func(t *testing.T) {
			//Get History:
			jsonStr := []byte(fmt.Sprintf(`{"language":"%s"}`, "icelandic"))
			responseBod := GetRequest(string(jsonStr), nil, t)

			for _, author := range AODICEs {
				reg := regexp.MustCompile(author.Date)
				if !reg.MatchString(responseBod) {
					t.Fatalf("Missing Aodice for date: %s", author.Date)
				}
			}

		})

		t.Run("Should get complete history of AODs from 2020-01-01", func(t *testing.T) {
			//Get History:
			jsonStr := []byte(fmt.Sprintf(`{"language":"%s","minimum":"%s"}`, "english", "2020-01-01"))
			responseBod := GetRequest(string(jsonStr), nil, t)

			shouldMatchReg := regexp.MustCompile(AODs[1].Date) //Regex for 2020-06-16 AODICE
			if !shouldMatchReg.MatchString(responseBod) {
				t.Fatalf("Expected the history of AODs to contain input AODICE for date %s but got body %s", AODs[1].Date, responseBod)
			}
			shouldNotMatchReg := regexp.MustCompile(AODs[2].Date) //Regex for 2019-06-16 AODICE
			if shouldNotMatchReg.MatchString(responseBod) {
				t.Fatalf("Expected the hisory of AODs only from 2020-01-01 but got body %s", responseBod)
			}

		})

		t.Run("Should get complete history of AODICEs from 2020-01-01", func(t *testing.T) {
			//Get History:
			jsonStr := []byte(fmt.Sprintf(`{"language":"%s","minimum":"%s"}`, "icelandic", "2020-01-01"))
			responseBod := GetRequest(string(jsonStr), nil, t)

			shouldMatchReg := regexp.MustCompile(AODICEs[1].Date) //Regex for 2020-06-16 AODICE
			if !shouldMatchReg.MatchString(responseBod) {
				t.Fatalf("Expected the history of AODICEs to contain input AODICE for date %s but got body %s", AODICEs[1].Date, responseBod)
			}
			shouldNotMatchReg := regexp.MustCompile(AODICEs[2].Date) //Regex for 2019-06-16 AODICE
			if shouldNotMatchReg.MatchString(responseBod) {
				t.Fatalf("Expected the hisory of AODICEs only from 2020-01-01 but got body %s", responseBod)
			}

		})

	})

}
