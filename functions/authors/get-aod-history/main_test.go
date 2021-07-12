package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func Setup() {

}
func TestHandler(t *testing.T) {
	var testingHandler = RequestHandler{}
	t.Run("Author of the day", func(t *testing.T) {

		t.Run("Should get complete history of Author of the day", func(t *testing.T) {
			//Get History:
			jsonStr := []byte(fmt.Sprintf(`{"language":"%s"}`, "english"))
			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})

			if err != nil {
				t.Fatalf("Expected the history of AOD but got an error: %+v", err)
			}

			log.Println("HERE:", response.Body)

			// if len(authors) == 0 {
			// 	t.Fatalf("Expected the history of AOD but got an empty list: %+v", authors)
			// }

			// containsBirfdayAuthor := false
			// containsTodayAuthor := false
			// const layout = "2006-01-02T15:04:05Z" //The date needed for reference always
			// for _, author := range authors {
			// 	if author.Id == 0 {
			// 		t.Fatalf("Expected all authors to have id > 0 but got: %+v", authors)
			// 	}
			// 	date, _ := time.Parse(layout, author.Date)
			// 	if date.Format("01-02-2006") == time.Now().Format("01-02-2006") {
			// 		containsTodayAuthor = true
			// 	}

			// 	if date.Format("01-02-2006") == "06-16-1998" {
			// 		containsBirfdayAuthor = true
			// 	}
			// }

			// if !containsBirfdayAuthor {
			// 	t.Fatalf("AOD history should contain the AOD for birfday but does not: %+v", authors)
			// }

			// if !containsTodayAuthor {
			// 	t.Fatalf("AOD history should contain the AOD for today but does not: %+v", authors)
			// }

		})

		// t.Run("Should get history of AOD starting from June 4th 2021", func(t *testing.T) {

		// 	//Input a quote in history for testing
		// 	authorId := 666
		// 	date := "2021-06-04"
		// 	var jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","aods": [{"id":%d, "date":"%s"}]}`, godUser.ApiKey, authorId, date))
		// 	_, response := requestAndReturnArray(jsonStr, SetAuthorOfTheDay)
		// 	if response.StatusCode != 200 {
		// 		t.Fatalf("Expected a succesful insert but got %+v", response)
		// 	}

		// 	//Get History:

		// 	minimum := "2021-06-04"
		// 	jsonStr = []byte(fmt.Sprintf(`{"apiKey":"%s","language":"%s", "minimum":"%s"}`, user.ApiKey, "english", minimum))
		// 	authors, _ := requestAndReturnArray(jsonStr, GetAODHistory)

		// 	if len(authors) == 0 {
		// 		t.Fatalf("Expected the history of AOD but got an empty list: %+v", authors)
		// 	}

		// 	const layout = "2006-01-02T15:04:05Z" //The date needed for reference always
		// 	compareDate, _ := time.Parse(layout, "2021-06-04")
		// 	compareYear := compareDate.Year()
		// 	compareMonth := compareDate.Month()
		// 	compareDay := compareDate.Day()
		// 	containsAuthorNotInRange := false
		// 	containsFourthOfJuneAuthor := false
		// 	for _, author := range authors {
		// 		date, _ := time.Parse(layout, author.Date)
		// 		yearOfAuthor := date.Year()
		// 		monthOfAuthor := date.Month()
		// 		dayOfAuthor := date.Day()

		// 		if yearOfAuthor < compareYear || (yearOfAuthor == compareYear && monthOfAuthor < compareMonth) || (yearOfAuthor == compareYear && monthOfAuthor == compareMonth && dayOfAuthor < compareDay) {
		// 			containsAuthorNotInRange = true
		// 		}

		// 		if date.Format("2006-01-02") == "2021-06-04" {
		// 			containsFourthOfJuneAuthor = true
		// 		}

		// 		if author.Id == 0 {
		// 			t.Fatalf("Expected all authors to have id > 0 but got: %+v", authors)
		// 		}

		// 	}

		// 	if containsAuthorNotInRange {
		// 		t.Fatalf("AOD history contains an earlier quote than was requested: %+v", authors)
		// 	}

		// 	if !containsFourthOfJuneAuthor {
		// 		t.Fatalf("QOD history should contain the QOD for 4th of june 2021 but does not: %+v", authors)
		// 	}

		// })

	})

}
