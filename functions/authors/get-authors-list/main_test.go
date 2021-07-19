package main

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
	"github.com/aws/aws-lambda-go/events"
)

//Returns AODs and AODICEs, in that order, put into the DB
func Setup(handler *RequestHandler, t *testing.T) {
	handler.InitializeDB()
	authorId := 1
	countIncrease := 10
	handler.Db.Table("authors").Update("count = count + ?", countIncrease).Where("id = ?", authorId)

	//CleanUp
	t.Cleanup(func() {
		handler.Db.Exec("update authors set count = 0")
	})
}
func TestHandler(t *testing.T) {
	var testingHandler = RequestHandler{}
	Setup(&testingHandler, t)

	t.Run("Time Test for getting authors", func(t *testing.T) {

		maxTime := 50
		longTime := 100
		t.Run("Should return first authors starting from 'F' (i.e. greater than or equal to 'F' alphabetically)", func(t *testing.T) {
			start := time.Now()
			minimum := "f"
			var jsonStr = []byte(fmt.Sprintf(`{ "orderConfig":{"orderBy":"alphabetical","minimum":"%s"}}`, minimum))

			_, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(longTime) {
				t.Fatalf("Expected getting authors to take less than %dms but it took %dms", longTime, duration.Milliseconds())
			}

		})

		t.Run("Should return first authors with less than 10 quotes in total in reversed order (start with those with 10 quotes)", func(t *testing.T) {
			maximum := 10
			start := time.Now()
			var jsonStr = []byte(fmt.Sprintf(`{"orderConfig":{"orderBy":"nrOfQuotes","maximum":"%d","reverse":true}}`, maximum))

			_, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}
			end := time.Now()
			duration := end.Sub(start)

			if duration.Milliseconds() > int64(maxTime) {
				t.Fatalf("Expected getting authors to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}

		})

		t.Run("Should return first 50 authors (ordered by most popular, i.e. DESC count)", func(t *testing.T) {
			start := time.Now()
			var jsonStr = []byte(fmt.Sprintf(`{"orderConfig":{"orderBy":"%s"}}`, "popularity"))

			_, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(maxTime) {
				t.Fatalf("Expected getting authors to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}

		})

		t.Run("Should return first 50 authors of profession 'Rapper' or 'Designer' ", func(t *testing.T) {
			start := time.Now()
			profession1 := "rApPeR"
			profession2 := "DESIGNER"
			var jsonStr = []byte(fmt.Sprintf(`{"professions":["%s","%s"]}`, profession2, profession1))

			_, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(maxTime) {
				t.Fatalf("Expected getting authors to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}
		})

		t.Run("Should return first 50 authors born in June", func(t *testing.T) {
			start := time.Now()
			birthMonth := "JUNE"

			var jsonStr = []byte(fmt.Sprintf(`{"time":{"born":{"month":"%s"}}}`, birthMonth))

			_, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(maxTime) {
				t.Fatalf("Expected getting authors to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}
		})

		t.Run("Should return first 50 authors alive today and older than 70", func(t *testing.T) {
			start := time.Now()
			minAge := 70
			var jsonStr = []byte(fmt.Sprintf(`{"time":{"isAlive":%t,"age":{"olderThan":%d}}, "orderConfig":{"orderBy":"DATEOFDEATH", "reverse":true}}`, true, minAge))

			_, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(maxTime) {
				t.Fatalf("Expected getting authors to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}

		})

		t.Run("Should return first 50 authors order by age desc", func(t *testing.T) {
			start := time.Now()
			var jsonStr = []byte(fmt.Sprintf(`{"orderConfig":{"orderBy":"AGE", "reverse":%t}}`, true))

			_, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}
			end := time.Now()
			duration := end.Sub(start)
			if duration.Milliseconds() > int64(maxTime) {
				t.Fatalf("Expected getting authors to take less than %dms but it took %dms", maxTime, duration.Milliseconds())
			}
		})

	})

	t.Run("Get authors", func(t *testing.T) {

		t.Run("Should return first 50 authors (alphabetically)", func(t *testing.T) {

			pageSize := 50
			var jsonStr = []byte(fmt.Sprintf(`{"pageSize": %d}`, pageSize))
			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)

			if len(authors) != 50 {
				t.Fatalf("got list of length %d, but expected list of length %d", len(authors), pageSize)
			}

			firstAuthor := authors[0]
			if firstAuthor.Name[0] != '2' {
				t.Fatalf("got %s, want name that starts with '2'", firstAuthor.Name)
			}

		})

		t.Run("Should return first authors, with only English quotes starting from A, (alphabetically)", func(t *testing.T) {

			language := "english"
			var jsonStr = []byte(fmt.Sprintf(`{"language": "%s","orderConfig":{"minimum":"A"}}`, language))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)

			firstAuthor := authors[0]

			if firstAuthor.NrOfIcelandicQuotes > 0 {
				t.Fatalf("got %+v, but expected an author that has no icelandic quotes", firstAuthor)
			}

			if firstAuthor.Name[0] != 'A' {
				t.Fatalf("got %s, want name that starts with 'A'", firstAuthor.Name)
			}

		})

		t.Run("Should return first English authors in reverse alphabetical order (i.e. first author starts with Z)", func(t *testing.T) {

			language := "english"
			var jsonStr = []byte(fmt.Sprintf(`{"language": "%s", "orderConfig":{"orderBy":"alphabetical", "reverse":true, "maximum":"Z"}}`, language))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)

			firstAuthor := authors[0]

			if firstAuthor.NrOfIcelandicQuotes > 0 {
				t.Fatalf("got %+v, but expected an author that has no icelandic quotes", firstAuthor)
			}

			if firstAuthor.Name[0] != 'Z' {
				t.Fatalf("got %s, want name that starts with 'Z'", firstAuthor.Name)
			}

		})

		t.Run("Should return first authors starting from 'F' (i.e. greater than or equal to 'F' alphabetically)", func(t *testing.T) {
			minimum := "f"
			var jsonStr = []byte(fmt.Sprintf(`{ "orderConfig":{"orderBy":"alphabetical","minimum":"%s"}}`, minimum))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)

			firstAuthor := authors[0]

			if firstAuthor.Name[0] != strings.ToUpper(minimum)[0] {
				t.Fatalf("got %s, want name that starts with 'F'", firstAuthor.Name)
			}

		})

		t.Run("Should return authors with less than or equal to 1 quotes in total", func(t *testing.T) {

			maximum := 1
			var jsonStr = []byte(fmt.Sprintf(`{"orderConfig":{"orderBy":"nrOfQuotes","maximum":"%d"}}`, maximum))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)

			firstAuthor := authors[0]

			if firstAuthor.NrOfIcelandicQuotes+firstAuthor.NrOfEnglishQuotes > 1 {
				t.Fatalf("got %+v, but expected an author that has no more than 1 quotes", firstAuthor)
			}

		})

		t.Run("Should return first authors with more than 10 quotes but less than or equal to 11 in total", func(t *testing.T) {

			minimum := 10
			maximum := 11
			var jsonStr = []byte(fmt.Sprintf(`{"orderConfig":{"orderBy":"nrOfQuotes","maximum":"%d", "minimum":"%d"}}`, maximum, minimum))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)

			firstAuthor := authors[0]

			if firstAuthor.NrOfIcelandicQuotes+firstAuthor.NrOfEnglishQuotes != 10 {
				t.Fatalf("got %+v, but expected an author that has no fewer than 10 quotes", firstAuthor)
			}

		})

		t.Run("Should return first authors with less than 10 quotes in total in reversed order (start with those with 10 quotes)", func(t *testing.T) {

			maximum := 10
			var jsonStr = []byte(fmt.Sprintf(`{"orderConfig":{"orderBy":"nrOfQuotes","maximum":"%d","reverse":true}}`, maximum))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)

			firstAuthor := authors[0]

			if firstAuthor.NrOfIcelandicQuotes+firstAuthor.NrOfEnglishQuotes != 10 {
				t.Fatalf("got %+v, but expected an author that has 10 quotes", firstAuthor)
			}

		})

		t.Run("Should return first authors (reverse order DESC by nr of quotes) only icelandic quotes", func(t *testing.T) {
			language := "icelandic"
			var jsonStr = []byte(fmt.Sprintf(`{"language":"%s", "orderConfig":{"orderBy":"nrOfQuotes","reverse":true}}`, language))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)

			firstAuthor := authors[0]

			if firstAuthor.Name != "Óþekktur höfundur" {
				t.Fatalf("got %+v, but expected the óþekktur höfundur author", firstAuthor)
			}
		})

		t.Run("Should return first 50 authors (ordered by most popular, i.e. DESC count)", func(t *testing.T) {

			var jsonStr = []byte(fmt.Sprintf(`{"orderConfig":{"orderBy":"%s"}}`, "popularity"))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)

			firstAuthor := authors[0]

			if firstAuthor.Count == 0 {
				t.Fatalf("got %+v, but expected an author that does not have 0 popularity count", firstAuthor)
			}

		})

		t.Run("Should return first 50 authors in reverse popularity order (i.e. least popular first i.e. ASC count)", func(t *testing.T) {

			var jsonStr = []byte(fmt.Sprintf(`{"orderConfig":{"orderBy":"%s","reverse":true}}`, "popularity"))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)

			firstAuthor := authors[0]

			if firstAuthor.Count != 0 {
				t.Fatalf("got %+v, but expected an author that has 0 popularity count", firstAuthor)
			}

		})

		t.Run("Should return first 100 authors", func(t *testing.T) {
			pageSize := 100
			var jsonStr = []byte(fmt.Sprintf(`{"pageSize":%d}`, pageSize))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)

			if len(authors) != 100 {
				t.Fatalf("got %d nr of authors, but expected %d authors", len(authors), pageSize)
			}
		})

		t.Run("Should return the next 50 authors starting from 'F' (i.e. pagination, page 1, alphabetical order)", func(t *testing.T) {

			pageSize := 100
			minimum := "F"
			var jsonStr = []byte(fmt.Sprintf(`{"pageSize":%d, "orderConfig":{"minimum":"%s"}}`, pageSize, minimum))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)

			objToFetch := authors[50]

			if authors[0].Name[0] != minimum[0] {
				t.Fatalf("got %+v, but expected author starting with '%s'", authors[0], minimum)
			}

			pageSize = 50
			page := 1
			jsonStr = []byte(fmt.Sprintf(`{ "pageSize":%d, "page":%d, "orderConfig":{"minimum":"%s"}}`, pageSize, page, minimum))

			response, err = testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			json.Unmarshal([]byte(response.Body), &authors)

			if objToFetch.Id != authors[0].Id {
				t.Fatalf("got %+v, but expected %+v", authors[0], objToFetch)
			}

		})

		t.Run("Should return first 50 authors of profession 'Politician'", func(t *testing.T) {
			profession := "POLITICIAN"
			var jsonStr = []byte(fmt.Sprintf(`{"professions":["%s"]}`, profession))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)

			if authors[0].Profession != "Politician" {
				t.Fatalf("Got author of profession %s, expected author of profession %s", authors[0].Profession, "Politician")
			}
		})
		t.Run("Should return first 50 authors of profession 'Rapper' or 'Designer' ", func(t *testing.T) {
			profession1 := "rApPeR"
			profession2 := "DESIGNER"
			var jsonStr = []byte(fmt.Sprintf(`{"professions":["%s","%s"]}`, profession2, profession1))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)

			gotDesigner := false
			gotRapper := false
			for _, author := range authors {
				if author.Profession != "Rapper" && author.Profession != "Designer" {
					t.Fatalf("Got author of profession %s, expected author of profession %s or %s", author.Profession, "Designer", "Rapper")
				}
				if author.Profession == "Rapper" {
					gotRapper = true
				}
				if author.Profession == "Designer" {
					gotDesigner = true
				}
			}

			if !gotDesigner || !gotRapper {
				t.Fatalf("Expected to get at least one rapper and one designer but gotDesigner: %t and gotRapper: %t", gotDesigner, gotRapper)
			}

		})

		t.Run("Should return first 50 authors of nationality 'Zimbabwean'", func(t *testing.T) {
			nationality := "ZIMBABWEAN"
			var jsonStr = []byte(fmt.Sprintf(`{"nationalities":["%s"]}`, nationality))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			if authors[0].Nationality != "Zimbabwean" {
				t.Fatalf("Got author of nationality %s, expected author of nationality %s", authors[0].Nationality, "Zimbabwean")
			}

		})
		t.Run("Should return first 50 authors of nationality 'Italian' and 'French'", func(t *testing.T) {
			nationality1 := "ITALIAN"
			nationality2 := "FRENCH"
			var jsonStr = []byte(fmt.Sprintf(`{"nationalities":["%s","%s"]}`, nationality1, nationality2))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			gotItalian := false
			gotFrench := false
			for _, author := range authors {
				if author.Nationality != "Italian" && author.Nationality != "French" {
					t.Fatalf("Got author of nationality %s, expected author of nationality %s or %s", author.Nationality, "Italian", "French")
				}
				if author.Nationality == "Italian" {
					gotItalian = true
				}
				if author.Nationality == "French" {
					gotFrench = true
				}
			}

			if !gotItalian || !gotFrench {
				t.Fatalf("Expected to get at least one italian and one french but gotItalian: %t and gotFrench: %t", gotItalian, gotFrench)
			}

		})

		t.Run("Should return first 50 authors born in 1956", func(t *testing.T) {
			birthYear := 1956

			var jsonStr = []byte(fmt.Sprintf(`{"time":{"born":{"year":%d}}}`, birthYear))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]

			reg := regexp.MustCompile(strconv.Itoa(birthYear))
			if !reg.Match([]byte(firstAuthor.BirthDate)) {
				t.Fatalf("Expected to get author born in %d but got author born in %s", birthYear, firstAuthor.BirthDate)
			}
		})
		t.Run("Should return first 50 authors born in June", func(t *testing.T) {
			birthMonth := "JUNE"

			var jsonStr = []byte(fmt.Sprintf(`{"time":{"born":{"month":"%s"}}}`, birthMonth))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]

			reg := regexp.MustCompile("June")
			if !reg.Match([]byte(firstAuthor.BirthDate)) {
				t.Fatalf("Expected to get author born in %s but got author born in %s", "June", firstAuthor.BirthDate)
			}
		})

		t.Run("Should return first 50 authors born on June-16", func(t *testing.T) {
			birthMonth := "JUNE"
			birthDate := 16
			var jsonStr = []byte(fmt.Sprintf(`{"time":{"born":{"month":"%s", "date":%d}}}`, birthMonth, birthDate))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]

			reg := regexp.MustCompile("June-16")
			if !reg.Match([]byte(firstAuthor.BirthDate)) {
				t.Fatalf("Expected to get author born in %s but got author born in %s", "June-16", firstAuthor.BirthDate)
			}
		})

		t.Run("Should return first 50 authors born on 1998-June-16", func(t *testing.T) {
			birthMonth := "JUNE"
			birthDate := 16
			birthYear := 1998
			var jsonStr = []byte(fmt.Sprintf(`{"time":{"born":{"year":%d,"month":"%s", "date":%d}}}`, birthYear, birthMonth, birthDate))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]

			reg := regexp.MustCompile("1998-June-16")
			if !reg.Match([]byte(firstAuthor.BirthDate)) {
				t.Fatalf("Expected to get author born on %s but got author born on %s", "1998-June-16", firstAuthor.BirthDate)
			}
		})

		t.Run("Should return first 50 authors born in 1956 order by date of birth ASC", func(t *testing.T) {
			birthYear := 1956
			var jsonStr = []byte(fmt.Sprintf(`{"time":{"born":{"year":%d}}, "orderConfig":{"orderBy":"DATEOFBIRTH"}}`, birthYear))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]

			reg := regexp.MustCompile("1956-January-")
			if !reg.Match([]byte(firstAuthor.BirthDate)) {
				t.Fatalf("Expected to get author born in %s but got author born in %s", "1956-January-", firstAuthor.BirthDate)
			}
		})
		t.Run("Should return first 50 authors born in 1956 order by date of birth DESC", func(t *testing.T) {
			birthYear := 1956
			var jsonStr = []byte(fmt.Sprintf(`{"time":{"born":{"year":%d}}, "orderConfig":{"orderBy":"DATEOFBIRTH", "reverse":true}}`, birthYear))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]

			//Do it like this because the authors who where "born" the latest in that year are all the authors that did not have months/dates, only a birth_year
			//Therefore the expected returned DeathDate is "1956" (not with "-").
			reg := regexp.MustCompile("1956-")
			reg2 := regexp.MustCompile("1956")
			if reg.Match([]byte(firstAuthor.BirthDate)) || !reg2.Match([]byte(firstAuthor.BirthDate)) {
				t.Fatalf("Expected to get author born in %s but got author born in %s", "1956", firstAuthor.BirthDate)
			}
		})

		t.Run("Should return first 50 authors that died in 1956", func(t *testing.T) {
			deathYear := 1956

			var jsonStr = []byte(fmt.Sprintf(`{"time":{"died":{"year":%d}}}`, deathYear))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]

			reg := regexp.MustCompile(strconv.Itoa(deathYear))
			if !reg.Match([]byte(firstAuthor.DeathDate)) {
				t.Fatalf("Expected to get author died in %d but got author died in %s", deathYear, firstAuthor.DeathDate)
			}
		})
		t.Run("Should return first 50 authors died in June", func(t *testing.T) {
			deathMonth := "JUNE"

			var jsonStr = []byte(fmt.Sprintf(`{"time":{"died":{"month":"%s"}}}`, deathMonth))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]

			reg := regexp.MustCompile("June")
			if !reg.Match([]byte(firstAuthor.DeathDate)) {
				t.Fatalf("Expected to get author died in %s but got author died in %s", "June", firstAuthor.DeathDate)
			}
		})
		t.Run("Should return first 50 authors died on June-16", func(t *testing.T) {
			deathMonth := "JUNE"
			deathDate := 16
			var jsonStr = []byte(fmt.Sprintf(`{"time":{"died":{"month":"%s", "date":%d}}}`, deathMonth, deathDate))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]

			reg := regexp.MustCompile("June-16")
			if !reg.Match([]byte(firstAuthor.DeathDate)) {
				t.Fatalf("Expected to get author died on %s but got author died on %s", "June-16", firstAuthor.DeathDate)
			}
		})
		t.Run("Should return first 50 authors died on 1960-June-16", func(t *testing.T) {
			deathMonth := "JUNE"
			deathDate := 16
			deathYear := 1960
			var jsonStr = []byte(fmt.Sprintf(`{"time":{"died":{"year":%d,"month":"%s", "date":%d}}}`, deathYear, deathMonth, deathDate))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]

			reg := regexp.MustCompile(fmt.Sprintf("%d-%s-%d", deathYear, "June", deathDate))
			if !reg.Match([]byte(firstAuthor.DeathDate)) {
				t.Fatalf("Expected to get author died on %s but got author died on %s", fmt.Sprintf("%d-%s-%d", deathYear, deathMonth, deathDate), firstAuthor.DeathDate)
			}
		})
		t.Run("Should return first 50 authors born in 1956 order by date of death ASC", func(t *testing.T) {
			deathYear := 1956
			var jsonStr = []byte(fmt.Sprintf(`{"time":{"died":{"year":%d}}, "orderConfig":{"orderBy":"DATEOFDEATH"}}`, deathYear))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]

			reg := regexp.MustCompile("1956-January-")
			if !reg.Match([]byte(firstAuthor.DeathDate)) {
				t.Fatalf("Expected to get author died in %s but got author died in %s", "1956-January-", firstAuthor.DeathDate)
			}
		})
		t.Run("Should return first 50 authors born in 1956 order by date of death DESC", func(t *testing.T) {
			deathYear := 1956
			var jsonStr = []byte(fmt.Sprintf(`{"time":{"died":{"year":%d}}, "orderConfig":{"orderBy":"DATEOFDEATH", "reverse":true}}`, deathYear))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]

			//Do it like this because the authors who "died" the latest in that year are all the authors that did not have months/dates, only a birth_year.
			//Therefore the expected returned DeathDate is "1956" (not with "-").
			reg := regexp.MustCompile("1956-")
			reg2 := regexp.MustCompile("1956")
			if reg.Match([]byte(firstAuthor.DeathDate)) || !reg2.Match([]byte(firstAuthor.DeathDate)) {
				t.Fatalf("Expected to get author died in %s but got author died in %s", "1956", firstAuthor.DeathDate)
			}
		})

		t.Run("Should return first 50 authors born before 1967-01-27", func(t *testing.T) {
			beforeYear := "1967-01-27"
			var jsonStr = []byte(fmt.Sprintf(`{"time":{"born":{"before":"%s"}}, "orderConfig":{"orderBy":"DATEOFBIRTH", "reverse":true}}`, beforeYear))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]

			reg := regexp.MustCompile("1967-January-26")
			if !reg.Match([]byte(firstAuthor.BirthDate)) {
				t.Fatalf("Expected to get author born on %s but got author born on %s", "1967-January-26", firstAuthor.BirthDate)
			}
		})
		t.Run("Should return first 50 authors born after 1967-January-27", func(t *testing.T) {
			afterYear := "1967-01-27"
			var jsonStr = []byte(fmt.Sprintf(`{"time":{"born":{"after":"%s"}}, "orderConfig":{"orderBy":"DATEOFBIRTH", "reverse":false}}`, afterYear))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]

			reg := regexp.MustCompile("1967-January-28")
			if !reg.Match([]byte(firstAuthor.BirthDate)) {
				t.Fatalf("Expected to get author born on %s but got author born on %s", "1967-January-28", firstAuthor.BirthDate)
			}
		})
		t.Run("Should return first 50 authors born between 1967-January-29 <-> 1967-February-28", func(t *testing.T) {
			afterYear := "1967-01-29"
			beforeYear := "1967-02-28"
			var jsonStr = []byte(fmt.Sprintf(`{"time":{"born":{"after":"%s", "before":"%s"}}, "orderConfig":{"orderBy":"DATEOFBIRTH", "reverse":false}}`, afterYear, beforeYear))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]

			reg := regexp.MustCompile("1967-January-30")
			if !reg.Match([]byte(firstAuthor.BirthDate)) {
				t.Fatalf("Expected to get author born on %s but got author born on %s", "1967-January-30", firstAuthor.BirthDate)
			}

			//Now make the same request but reverse order to get first the authors closes to born on 1967-February-28

			jsonStr = []byte(fmt.Sprintf(`{"time":{"born":{"after":"%s", "before":"%s"}}, "orderConfig":{"orderBy":"DATEOFBIRTH", "reverse":true}}`, afterYear, beforeYear))

			response, err = testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor = authors[0]

			reg = regexp.MustCompile("1967-February-28")
			if !reg.Match([]byte(firstAuthor.BirthDate)) {
				t.Fatalf("Expected to get author born on %s but got author born on %s", "1967-February-28", firstAuthor.BirthDate)
			}

		})

		t.Run("Should return first 50 authors died before 1967-January-28", func(t *testing.T) {
			beforeYear := "1967-01-28"
			var jsonStr = []byte(fmt.Sprintf(`{"time":{"died":{"before":"%s"}}, "orderConfig":{"orderBy":"DATEOFDEATH", "reverse":true}}`, beforeYear))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]

			reg := regexp.MustCompile("1967-January-27")
			if !reg.Match([]byte(firstAuthor.DeathDate)) {
				t.Fatalf("Expected to get author died on %s but got author died on %s", "1967-January-27", firstAuthor.DeathDate)
			}
		})
		t.Run("Should return first 50 authors died after 1967-January-28", func(t *testing.T) {
			afterYear := "1967-01-28"
			var jsonStr = []byte(fmt.Sprintf(`{"time":{"died":{"after":"%s"}}, "orderConfig":{"orderBy":"DATEOFDEATH", "reverse":false}}`, afterYear))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]

			reg := regexp.MustCompile("1967-February-7")
			if !reg.Match([]byte(firstAuthor.DeathDate)) {
				t.Fatalf("Expected to get author died on %s but got author died on %s", "1967-February-07", firstAuthor.DeathDate)
			}
		})
		t.Run("Should return first 50 authors died between 1967-January-28 - 1967-February-28", func(t *testing.T) {
			afterYear := "1967-01-29"
			beforeYear := "1967-02-28"
			var jsonStr = []byte(fmt.Sprintf(`{"time":{"died":{"after":"%s", "before":"%s"}}, "orderConfig":{"orderBy":"DATEOFDEATH", "reverse":false}}`, afterYear, beforeYear))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]

			reg := regexp.MustCompile("1967-February-7")
			if !reg.Match([]byte(firstAuthor.DeathDate)) {
				t.Fatalf("Expected to get author died on %s but got author died on %s", "1967-February-7", firstAuthor.DeathDate)
			}

			//Now make the same request but reverse order to get first the authors closes to born on 1967-February-28

			jsonStr = []byte(fmt.Sprintf(`{"time":{"died":{"after":"%s", "before":"%s"}}, "orderConfig":{"orderBy":"DATEOFDEATH", "reverse":true}}`, afterYear, beforeYear))

			response, err = testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor = authors[0]

			reg = regexp.MustCompile("1967-February-28")
			if !reg.Match([]byte(firstAuthor.DeathDate)) {
				t.Fatalf("Expected to get author died on %s but got author died on %s", "1967-February-28", firstAuthor.DeathDate)
			}
		})

		t.Run("Should return first 50 authors alive today", func(t *testing.T) {
			var jsonStr = []byte(fmt.Sprintf(`{"time":{"isAlive":%t}, "orderConfig":{"orderBy":"DATEOFDEATH", "reverse":true}}`, true))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]
			if firstAuthor.DeathDate != "" {
				t.Fatalf("Expected to get author still alive but got author died on %s", firstAuthor.DeathDate)
			}
		})

		t.Run("Should return first 50 authors alive today exactly 50 years old", func(t *testing.T) {
			exactAge := 50
			var jsonStr = []byte(fmt.Sprintf(`{"time":{"isAlive":%t,"age":{"exactly":%d}}, "orderConfig":{"orderBy":"DATEOFDEATH", "reverse":true}}`, true, exactAge))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]
			if firstAuthor.DeathDate != "" {
				t.Fatalf("Expected to get author still alive but got author died on %s", firstAuthor.DeathDate)
			}

			birth, err := time.Parse("2006-January-2", firstAuthor.BirthDate)
			if err != nil {
				t.Fatalf("Expected given author to have valid birthdate but got an error: %+v", err)
			}

			age := math.Floor(time.Since(birth).Hours() / 24 / 365)
			if int(age) != exactAge {
				t.Fatalf("Expected author of age exactly %d, but got author of age %d", exactAge, int(age))
			}
		})

		t.Run("Should return first 50 authors alive today and older than 70", func(t *testing.T) {
			minAge := 70
			var jsonStr = []byte(fmt.Sprintf(`{"time":{"isAlive":%t,"age":{"olderThan":%d}}, "orderConfig":{"orderBy":"DATEOFDEATH", "reverse":true}}`, true, minAge))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]
			if firstAuthor.DeathDate != "" {
				t.Fatalf("Expected to get author still alive but got author died on %s", firstAuthor.DeathDate)
			}

			birth, err := time.Parse("2006-January-2", firstAuthor.BirthDate)
			if err != nil {
				t.Fatalf("Expected given author to have valid birthdate but got an error: %+v", err)
			}

			age := math.Floor(time.Since(birth).Hours() / 24 / 365)
			if int(age) < minAge {
				t.Fatalf("Expected author of age at least %d, but got author of age %d", minAge, int(age))
			}

		})

		t.Run("Should return first 50 authors alive today and younger than 40", func(t *testing.T) {
			maxAge := 40
			var jsonStr = []byte(fmt.Sprintf(`{"time":{"isAlive":%t,"age":{"youngerThan":%d}}, "orderConfig":{"orderBy":"DATEOFDEATH", "reverse":true}}`, true, maxAge))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]
			if firstAuthor.DeathDate != "" {
				t.Fatalf("Expected to get author still alive but got author died on %s", firstAuthor.DeathDate)
			}

			birth, err := time.Parse("2006-January-2", firstAuthor.BirthDate)
			if err != nil {
				t.Fatalf("Expected given author to have valid birthdate but got an error: %+v", err)
			}

			age := math.Floor(time.Since(birth).Hours() / 24 / 365)
			if int(age) > maxAge {
				t.Fatalf("Expected author of age at most %d, but got author of age %d", maxAge, int(age))
			}
		})

		t.Run("Should return first 50 authors dead and older than 70", func(t *testing.T) {
			minAge := 70
			var jsonStr = []byte(fmt.Sprintf(`{"time":{"isDead":%t,"age":{"olderThan":%d}}, "orderConfig":{"orderBy":"DATEOFDEATH", "reverse":true}}`, true, minAge))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]
			if firstAuthor.DeathDate == "" {
				t.Fatalf("Expected to get author who is dead but got author who has no death date, %+v", firstAuthor)
			}

			birth, err := time.Parse("2006-January-2", firstAuthor.BirthDate)
			if err != nil {
				t.Fatalf("Expected given author to have valid birthdate but got an error: %+v", err)
			}

			age := math.Floor(time.Since(birth).Hours() / 24 / 365)
			if int(age) < minAge {
				t.Fatalf("Expected author of age at least %d, but got author of age %d", minAge, int(age))
			}
		})
		t.Run("Should return first 50 authors dead and younger than 40", func(t *testing.T) {
			maxAge := 40
			var jsonStr = []byte(fmt.Sprintf(`{"time":{"isDead":%t,"age":{"youngerThan":%d}}, "orderConfig":{"orderBy":"DATEOFDEATH", "reverse":true}}`, true, maxAge))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]
			if firstAuthor.DeathDate == "" {
				t.Fatalf("Expected to get author who is dead but got author who has no death date, %+v", firstAuthor)
			}

			birth, err := time.Parse("2006-January-2", firstAuthor.BirthDate)
			if err != nil {
				t.Logf("Expected given author to have valid birthdate but got an error: %s, with author: %+v", err, firstAuthor)
				birth, err = time.Parse("2006", firstAuthor.BirthDate)
				if err != nil {
					t.Fatalf("Expected given author to have at least birth_year but got an error: %s, with author: %+v", err, firstAuthor)
				}
			}

			age := math.Floor(time.Since(birth).Hours() / 24 / 365)
			if int(age) > maxAge {
				t.Fatalf("Expected author of age at most %d, but got author of age %d", maxAge, int(age))
			}
		})
		t.Run("Should return first 50 authors dead today", func(t *testing.T) {

			var jsonStr = []byte(fmt.Sprintf(`{"time":{"isDead":%t}, "orderConfig":{"orderBy":"DATEOFDEATH", "reverse":true}}`, true))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]
			if firstAuthor.DeathDate == "" {
				t.Fatalf("Expected to get author dead but got author still alive %+v", firstAuthor)
			}

		})
		t.Run("Should return first 50 authors dead today and exactly 50 years old", func(t *testing.T) {
			exactAge := 50
			var jsonStr = []byte(fmt.Sprintf(`{"time":{"isDead":%t,"age":{"exactly":%d}}, "orderConfig":{"orderBy":"DATEOFDEATH", "reverse":true}}`, true, exactAge))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]
			if firstAuthor.DeathDate == "" {
				t.Fatalf("Expected to get author dead but got author alive: %+v", firstAuthor)
			}

			birth, err := time.Parse("2006-January-2", firstAuthor.BirthDate)
			if err != nil {
				t.Fatalf("Expected given author to have valid birthdate but got an error: %+v", err)
			}

			death, err := time.Parse("2006-January-2", firstAuthor.DeathDate)
			if err != nil {
				t.Fatalf("Expected given author to have valid deathdate but got an error: %+v", err)
			}

			age := math.Floor(death.Sub(birth).Hours() / 24 / 365)
			if int(age) != exactAge {
				t.Fatalf("Expected author of age exactly %d, but got author of age %d", exactAge, int(age))
			}
		})

		t.Run("Should return first 50 authors dead or alive and exactly 50 years old", func(t *testing.T) {
			exactAge := 50
			var jsonStr = []byte(fmt.Sprintf(`{"time":{"age":{"exactly":%d}}, "orderConfig":{"orderBy":"DATEOFDEATH", "reverse":true}}`, exactAge))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]

			birth, err := time.Parse("2006-January-2", firstAuthor.BirthDate)
			if err != nil {
				t.Fatalf("Expected given author to have valid birthdate but got an error: %+v", err)
			}

			death, err := time.Parse("2006-January-2", firstAuthor.DeathDate)
			if err != nil {
				t.Fatalf("Expected given author to have valid deathdate but got an error: %+v", err)
			}

			age := math.Floor(death.Sub(birth).Hours() / 24 / 365)
			if int(age) != exactAge {
				t.Fatalf("Expected author of age exactly %d, but got author of age %d", exactAge, int(age))
			}
		})

		t.Run("Should return first 50 authors dead or alive and younger than 40", func(t *testing.T) { t.Skip() })
		t.Run("Should return first 50 authors dead or alive and older than 70", func(t *testing.T) { t.Skip() })

		t.Run("Should return first 50 authors order by age asc", func(t *testing.T) {
			var jsonStr = []byte(fmt.Sprintf(`{"orderConfig":{"orderBy":"AGE", "reverse":%t}}`, false))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]

			birth, err := time.Parse("2006-January-2", firstAuthor.BirthDate)
			if err != nil {
				birth, err = time.Parse("2006", firstAuthor.BirthDate)
				if err != nil {
					t.Fatalf("Expected given author to have valid birthdate but got an error: %+v", err)
				}
			}

			lastAge := math.Floor(time.Since(birth).Hours() / 24 / 365)
			for _, author := range authors {
				birth, err = time.Parse("2006-January-2", author.BirthDate)
				if err != nil {
					birth, err = time.Parse("2006", firstAuthor.BirthDate)
					if err != nil {
						t.Fatalf("Expected given author to have valid birthdate but got an error: %+v", err)
					}
				}
				currentAge := math.Floor(time.Since(birth).Hours() / 24 / 365)
				if lastAge > currentAge {
					t.Fatalf("Expected authors in the ascending order of age but did not get that son!")
				}
				lastAge = currentAge
			}
		})
		t.Run("Should return first 50 authors order by age desc", func(t *testing.T) {
			var jsonStr = []byte(fmt.Sprintf(`{"orderConfig":{"orderBy":"AGE", "reverse":%t}}`, true))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var authors []structs.AuthorAPIModel
			json.Unmarshal([]byte(response.Body), &authors)
			firstAuthor := authors[0]

			birth, err := time.Parse("2006-January-2", firstAuthor.BirthDate)
			if err != nil {
				birth, err = time.Parse("2006", firstAuthor.BirthDate)
				if err != nil {
					t.Fatalf("Expected given author to have valid birthdate but got an error: %+v", err)
				}
			}

			lastAge := math.Floor(time.Since(birth).Hours() / 24 / 365)
			for _, author := range authors {
				birth, err = time.Parse("2006-January-2", author.BirthDate)
				if err != nil {
					birth, err = time.Parse("2006", firstAuthor.BirthDate)
					if err != nil {
						t.Fatalf("Expected given author to have valid birthdate but got an error: %+v", err)
					}
				}
				currentAge := math.Floor(time.Since(birth).Hours() / 24 / 365)
				if lastAge < currentAge {
					t.Fatalf("Expected authors in the descending order of age but did not get that son!")
				}
				lastAge = currentAge
			}
		})

		t.Run("Should return first 50 authors order by date of birth", func(t *testing.T) { t.Skip() })
		t.Run("Should return first 50 authors order by date of death", func(t *testing.T) { t.Skip() })

		t.Run("Should return error message because there are no authors that match the request", func(t *testing.T) {
			var jsonStr = []byte(fmt.Sprintf(`{"time":{"born":{"before":"2000-June-01","after":"2001-June-01"}},"orderConfig":{"orderBy":"AGE", "reverse":%t}}`, true))

			response, err := testingHandler.handler(events.APIGatewayProxyRequest{Body: string(jsonStr)})
			if err != nil {
				t.Fatalf("Expected given author but got an error: %+v", err)
			}

			var errorResponse structs.ErrorResponse
			json.Unmarshal([]byte(response.Body), &errorResponse)
			if errorResponse.Message == "" || errorResponse.StatusCode != 200 {
				t.Fatalf("Should have gotten error message with 200 httpstatus code because no author matches the search but got %+v", errorResponse)
			}
		})

	})

}
