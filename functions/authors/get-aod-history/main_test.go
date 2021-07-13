package main

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
	"github.com/aws/aws-lambda-go/events"
)

//Returns AODs and AODICEs, in that order, put into the DB
func Setup(handler *RequestHandler, t *testing.T) ([]structs.AodDBModel, []structs.AodDBModel) {
	handler.InitializeDB()

	//Get the first 3 english- and icelandic-authors to be set as AODs/AODICEs
	var englishAuthors []structs.AuthorDBModel
	var icelandicAuthors []structs.AuthorDBModel
	err := handler.Db.Table("authors").Where("nr_of_english_quotes > 0").Find(&englishAuthors).Limit(3).Error
	if err != nil {
		t.Fatalf("Setup error: %s", err)
	}
	err = handler.Db.Table("authors").Where("nr_of_icelandic_quotes > 0").Find(&icelandicAuthors).Limit(3).Error
	if err != nil {
		t.Fatalf("Setup error: %s", err)
	}

	// Insert AODs and AODICEs for 2021-06-16, 2021-06-16,2019-06-16
	dates := []string{"2021-06-16", "2020-06-16", "2019-06-16"}
	var AODs []structs.AodDBModel
	var AODICEs []structs.AodDBModel
	for idx, date := range dates {
		AODs = append(AODs, englishAuthors[idx].ConvertToAODDBModel(date))
		AODICEs = append(AODICEs, icelandicAuthors[idx].ConvertToAODDBModel(date))
	}
	err = handler.Db.Table("aods").Create(&AODs).Error
	if err != nil {
		t.Fatalf("Setup error 2: %s", err)
	}
	err = handler.Db.Table("aodices").Create(&AODICEs).Error
	if err != nil {
		t.Fatalf("Setup error 2: %s", err)
	}

	//CleanUp
	t.Cleanup(func() {
		handler.Db.Unscoped().Table("aods").Delete(&AODs)
		handler.Db.Unscoped().Table("aodices").Delete(&AODICEs)
	})

	return AODs, AODICEs
}
func TestHandler(t *testing.T) {
	var testingHandler = RequestHandler{}
	AODs, AODICEs := Setup(&testingHandler, t)
	t.Run("AOD/AODICE History", func(t *testing.T) {

		t.Run("Should get complete history of AODs", func(t *testing.T) {
			//Get History:
			jsonStr := []byte(fmt.Sprintf(`{"language":"%s"}`, "english"))
			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected the history of AOD but got an error: %+v", err)
			}

			for _, author := range AODs {
				reg := regexp.MustCompile(author.Date)
				if !reg.MatchString(response.Body) {
					t.Fatalf("Missing Aod for date: %s", author.Date)
				}
			}

		})

		t.Run("Should get complete history of AODICEs", func(t *testing.T) {
			//Get History:
			jsonStr := []byte(fmt.Sprintf(`{"language":"%s"}`, "icelandic"))
			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected the history of AODICE but got an error: %+v", err)
			}

			for _, author := range AODICEs {
				reg := regexp.MustCompile(author.Date)
				if !reg.MatchString(response.Body) {
					t.Fatalf("Missing Aodice for date: %s", author.Date)
				}
			}

		})

		t.Run("Should get complete history of AODs from 2020-01-01", func(t *testing.T) {
			//Get History:
			jsonStr := []byte(fmt.Sprintf(`{"language":"%s","minimum":"%s"}`, "english", "2020-01-01"))
			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected the history of AOD but got an error: %+v", err)
			}

			shouldMatchReg := regexp.MustCompile(AODs[1].Date) //Regex for 2020-06-16 AODICE
			if !shouldMatchReg.MatchString(response.Body) {
				t.Fatalf("Expected the history of AODs to contain input AODICE for date %s but got body %s", AODs[1].Date, response.Body)
			}
			shouldNotMatchReg := regexp.MustCompile(AODs[2].Date) //Regex for 2019-06-16 AODICE
			if shouldNotMatchReg.MatchString(response.Body) {
				t.Fatalf("Expected the hisory of AODs only from 2020-01-01 but got body %s", response.Body)
			}

		})

		t.Run("Should get complete history of AODICEs from 2020-01-01", func(t *testing.T) {
			//Get History:
			jsonStr := []byte(fmt.Sprintf(`{"language":"%s","minimum":"%s"}`, "icelandic", "2020-01-01"))
			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected the history of AODICE but got an error: %+v", err)
			}

			shouldMatchReg := regexp.MustCompile(AODICEs[1].Date) //Regex for 2020-06-16 AODICE
			if !shouldMatchReg.MatchString(response.Body) {
				t.Fatalf("Expected the history of AODICEs to contain input AODICE for date %s but got body %s", AODICEs[1].Date, response.Body)
			}
			shouldNotMatchReg := regexp.MustCompile(AODICEs[2].Date) //Regex for 2019-06-16 AODICE
			if shouldNotMatchReg.MatchString(response.Body) {
				t.Fatalf("Expected the hisory of AODICEs only from 2020-01-01 but got body %s", response.Body)
			}

		})

	})

}
