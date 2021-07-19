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
func Setup(handler *RequestHandler, t *testing.T) ([]structs.QodDBModel, []structs.QodDBModel) {
	handler.InitializeDB()

	//Get the first 3 english- and icelandic-authors to be set as AODs/AODICEs
	var englishQuotes []structs.QuoteDBModel
	var icelandicQuotes []structs.QuoteDBModel
	err := handler.Db.Table("quotes").Not("is_icelandic").Limit(3).Find(&englishQuotes).Error
	if err != nil {
		t.Fatalf("Setup error: %s", err)
	}
	err = handler.Db.Table("quotes").Where("is_icelandic").Limit(3).Find(&icelandicQuotes).Error
	if err != nil {
		t.Fatalf("Setup error: %s", err)
	}

	// Insert AODs and AODICEs for 2021-06-16, 2021-06-16,2019-06-16
	dates := []string{"2021-06-16", "2020-06-16", "2019-06-16"}
	var QODs []structs.QodDBModel
	var QODICEs []structs.QodDBModel
	for idx, date := range dates {
		QODs = append(QODs, englishQuotes[idx].ConvertToQODDBModel(date))
		QODICEs = append(QODICEs, icelandicQuotes[idx].ConvertToQODDBModel(date))
	}
	createQODs := append(QODs, QODICEs...)
	err = handler.Db.Table("qods").Create(&createQODs).Error
	if err != nil {
		t.Fatalf("Setup error 2: %s", err)
	}

	//CleanUp
	t.Cleanup(func() {
		handler.Db.Exec("delete from qods")
	})

	return QODs, QODICEs
}
func TestHandler(t *testing.T) {
	QODs, QODICEs := Setup(&testingHandler, t)

	t.Run("Time Test for getting qod history", func(t *testing.T) {
		maxTime := 50
		t.Run("Time: Should get complete history of QODs", func(t *testing.T) {
			start := time.Now()
			//Get History:
			jsonStr := fmt.Sprintf(`{"language":"%s"}`, "english")
			GetRequest(jsonStr, nil, t)
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(maxTime) {
				t.Fatalf("Expected getting history of quotes to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}
		})
	})

	t.Run("Get quotes", func(t *testing.T) {
		t.Run("Should get complete history of QODs", func(t *testing.T) {
			//Get History:
			jsonStr := fmt.Sprintf(`{"language":"%s"}`, "english")
			responseBod := GetRequest(jsonStr, nil, t)
			for _, qod := range QODs {
				reg := regexp.MustCompile(qod.Date)
				if !reg.MatchString(responseBod) {
					t.Fatalf("Missing Qod for date: %s, got response %s", qod.Date, responseBod)
				}
			}
		})

		t.Run("Should get complete history of QODICEs", func(t *testing.T) {
			//Get History:
			jsonStr := fmt.Sprintf(`{"language":"%s"}`, "icelandic")
			responseBod := GetRequest(jsonStr, nil, t)

			for _, qod := range QODICEs {
				reg := regexp.MustCompile(qod.Date)
				if !reg.MatchString(responseBod) {
					t.Fatalf("Missing Qodice for date: %s, got response: %s", qod.Date, responseBod)
				}
			}

		})

		t.Run("Should get complete history of QODs from 2020-01-01", func(t *testing.T) {
			//Get History:
			jsonStr := fmt.Sprintf(`{"language":"%s","minimum":"%s"}`, "english", "2020-01-01")
			responseBod := GetRequest(jsonStr, nil, t)

			shouldMatchReg := regexp.MustCompile(QODs[1].Date) //Regex for 2020-06-16 QODs
			if !shouldMatchReg.MatchString(responseBod) {
				t.Fatalf("Expected the history of QODs to contain date %s but got body %s", QODs[1].Date, responseBod)
			}
			shouldNotMatchReg := regexp.MustCompile(QODs[2].Date) //Regex for 2019-06-16 QODs
			if shouldNotMatchReg.MatchString(responseBod) {
				t.Fatalf("Expected the hisory of QODs only from 2020-01-01 but got body %s", responseBod)
			}

		})

		t.Run("Should get complete history of QODICEs from 2020-01-01", func(t *testing.T) {
			//Get History:
			jsonStr := fmt.Sprintf(`{"language":"%s","minimum":"%s"}`, "icelandic", "2020-01-01")
			responseBod := GetRequest(jsonStr, nil, t)

			shouldMatchReg := regexp.MustCompile(QODICEs[1].Date) //Regex for 2020-06-16 AODICE
			if !shouldMatchReg.MatchString(responseBod) {
				t.Fatalf("Expected the history of QODICEs to contain input AODICE for date %s but got body %s", QODICEs[1].Date, responseBod)
			}
			shouldNotMatchReg := regexp.MustCompile(QODICEs[2].Date) //Regex for 2019-06-16 AODICE
			if shouldNotMatchReg.MatchString(responseBod) {
				t.Fatalf("Expected the hisory of QODICEs only from 2020-01-01 but got body %s", responseBod)
			}

		})

	})

}
