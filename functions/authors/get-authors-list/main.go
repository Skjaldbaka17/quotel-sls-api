package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gorm.io/gorm"
)

type RequestHandler struct {
	utils.RequestHandler
}

var theReqHandler = RequestHandler{}

func checkTime(requestBody *structs.Request, dbPointer *gorm.DB) *gorm.DB {
	if requestBody.Time != (structs.Time{}) {
		if requestBody.Time.Born != (structs.BeforeAfter{}) {
			if requestBody.Time.Born.Year != 0 {
				dbPointer = dbPointer.Where("birth_year = ?", requestBody.Time.Born.Year)
			}

			if requestBody.Time.Born.Month != "" {
				dbPointer = dbPointer.Where("birth_month = ?", strings.Title(strings.ToLower(requestBody.Time.Born.Month)))
			}

			if requestBody.Time.Born.Date != 0 {
				dbPointer = dbPointer.Where("birth_date = ?", requestBody.Time.Born.Date)
			}

			if requestBody.Time.Born.Before != "" {
				beforeDate := requestBody.Time.Born.Before + " 11:59PM" //Made to get those born on the day
				//birth_year > 0 because authorwise we get all the unknown when born authors first (order by DATEOFBIRTH asc!)
				dbPointer = dbPointer.Where("birth_year > 0").Where("birth_day <= ?", beforeDate)
			}

			if requestBody.Time.Born.After != "" {
				//birth_year > 0 because authorwise we get all the unknown when born authors first (order by DATEOFBIRTH asc!)
				dbPointer = dbPointer.Where("birth_year > 0").Where("birth_day >= ?", requestBody.Time.Born.After)
			}

		}

		if requestBody.Time.Died != (structs.BeforeAfter{}) {
			if requestBody.Time.Died.Year != 0 {
				dbPointer = dbPointer.Where("death_year = ?", requestBody.Time.Died.Year)
			}

			if requestBody.Time.Died.Month != "" {
				dbPointer = dbPointer.Where("death_month = ?", strings.Title(strings.ToLower(requestBody.Time.Died.Month)))
			}

			if requestBody.Time.Died.Date != 0 {
				dbPointer = dbPointer.Where("death_date = ?", requestBody.Time.Died.Date)
			}

			if requestBody.Time.Died.Before != "" {
				beforeDate := requestBody.Time.Died.Before + " 11:59PM" //Made to get those born on the day
				//birth_year > 0 because authorwise we get all the unknown when born authors first (order by DATEOFBIRTH asc!)
				dbPointer = dbPointer.Where("death_year > 0").Where("death_day <= ?", beforeDate)
			}

			if requestBody.Time.Died.After != "" {
				//birth_year > 0 because authorwise we get all the unknown when born authors first (order by DATEOFBIRTH asc!)
				dbPointer = dbPointer.Where("death_year > 0").Where("death_day >= ?", requestBody.Time.Died.After)
			}
		}

		if requestBody.Time.IsAlive {
			dbPointer = dbPointer.Where("birth_year > 1900").Not("death_year > 0")
		}

		if requestBody.Time.IsDead {
			dbPointer = dbPointer.Where("(death_year > 0 or birth_year < 1910)")
		}

		if requestBody.Time.Age != (structs.Age{}) {
			if requestBody.Time.Age.OlderThan != 0 {
				//Different for IsDead=true because then user likely wants authors who died at least the age of .Age.OlderThan
				if requestBody.Time.IsDead {
					dbPointer = dbPointer.Where("date_part('year', age(death_day,birth_day)) >= ?", requestBody.Time.Age.OlderThan).Where("birth_year > 0").Where("death_year > 0")
				} else {
					dbPointer = dbPointer.Where("date_part('year', age(birth_day)) >= ?", requestBody.Time.Age.OlderThan).Where("birth_year > 0")
				}
			}

			if requestBody.Time.Age.YoungerThan != 0 {
				//Different for IsDead=true because then user likely wants authors who died at younger than the age of .Age.YoungerThan
				if requestBody.Time.IsDead {
					dbPointer = dbPointer.Where("date_part('year', age(death_day,birth_day)) <= ?", requestBody.Time.Age.YoungerThan).Where("birth_year > 0").Where("death_year > 0")
				} else {
					dbPointer = dbPointer.Where("date_part('year', age(birth_day)) <= ?", requestBody.Time.Age.YoungerThan).Where("birth_year > 0")
				}

			}

			if requestBody.Time.Age.Exactly != 0 {
				if requestBody.Time.IsDead {
					dbPointer = dbPointer.Where("date_part('year', age(death_day,birth_day)) = ?", requestBody.Time.Age.Exactly).Where("birth_year > 0").Where("death_year > 0")
				}

				if requestBody.Time.IsAlive {
					dbPointer = dbPointer.Where("date_part('year', age(birth_day)) = ?", requestBody.Time.Age.Exactly)
				}

				if !requestBody.Time.IsAlive && !requestBody.Time.IsDead {
					dbPointer = dbPointer.Where("(date_part( 'year', age(death_day,birth_day) ) = ? and death_year > 0) or (date_part('year', age(birth_day)) = ? and death_year = 0) ", requestBody.Time.Age.Exactly, requestBody.Time.Age.Exactly).Where("birth_year > 0")
				}

			}
		}
	}
	return dbPointer
}

// swagger:route POST /authors/list authors ListAuthors
//
// List authors based on parameters
//
// Use this route to get a list of authors according to some ordering / parameters -- for example based on age, when they where born, on popularity, profession or nationalities and many more
//
// responses:
//	200: authorsResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// GetAuthorsList handles POST requests to get the authors that fit the parameters
func (requestHandler *RequestHandler) handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//Initialize DB if requestHandler.Db = nil
	if errResponse := requestHandler.InitializeDB(); errResponse != (structs.ErrorResponse{}) {
		return events.APIGatewayProxyResponse{
			Body:       errResponse.ToString(),
			StatusCode: errResponse.StatusCode,
		}, nil
	}

	requestBody, errResponse := requestHandler.ValidateRequest(request)

	if errResponse != (structs.ErrorResponse{}) {
		return events.APIGatewayProxyResponse{
			Body:       errResponse.ToString(),
			StatusCode: errResponse.StatusCode,
		}, nil
	}

	var authors []structs.AuthorDBModel
	//** ---------- Paramatere configuration for DB query begins ---------- **//
	dbPointer := requestHandler.Db.Table("authors")

	dbPointer = utils.AuthorLanguageSQL(requestBody.Language, dbPointer)

	if len(requestBody.Professions) > 0 {
		for idx, profession := range requestBody.Professions {
			requestBody.Professions[idx] = strings.Title(strings.ToLower(profession))
		}
		dbPointer = dbPointer.Where("profession in ?", requestBody.Professions)
	}

	if len(requestBody.Nationalities) > 0 {
		for idx, nationality := range requestBody.Nationalities {
			requestBody.Nationalities[idx] = strings.Title(strings.ToLower(nationality))
		}
		dbPointer = dbPointer.Where("nationality in ?", requestBody.Nationalities)
	}

	dbPointer = checkTime(&requestBody, dbPointer)
	orderDirection := "ASC"
	if requestBody.OrderConfig.Reverse {
		orderDirection = "DESC"
	}

	switch strings.ToLower(requestBody.OrderConfig.OrderBy) {
	case "popularity": //TODO: add popularity ordering
		orderDirection = "DESC"
		if requestBody.OrderConfig.Reverse {
			orderDirection = "ASC"
		}
		dbPointer = dbPointer.Order("count " + orderDirection)
	case "nrofquotes":
		switch strings.ToLower(requestBody.Language) {
		case "english":
			dbPointer = utils.SetMaxMinNumber(requestBody.OrderConfig, "nr_of_english_quotes", orderDirection, dbPointer)
		case "icelandic":
			dbPointer = utils.SetMaxMinNumber(requestBody.OrderConfig, "nr_of_icelandic_quotes", orderDirection, dbPointer)
		default:
			dbPointer = utils.SetMaxMinNumber(requestBody.OrderConfig, "nr_of_icelandic_quotes + nr_of_english_quotes", orderDirection, dbPointer)
		}
	case "dateofbirth":
		dbPointer = dbPointer.Order("birth_day " + orderDirection)
	case "dateofdeath":
		dbPointer = dbPointer.Order("death_day " + orderDirection)
	case "age":
		dbPointer = dbPointer.Where("birth_year > 0").Order("age(birth_day) " + orderDirection)
	default:
		//Minimum letter to start with (i.e. start from given minimum letter of the alphabet)
		if requestBody.OrderConfig.Minimum != "" {
			dbPointer = dbPointer.Where("left(name,1) >= ?", strings.ToUpper(requestBody.OrderConfig.Minimum))
		}
		//Maximum letter to start with (i.e. end at the given maximum letter of the alphabet)
		if requestBody.OrderConfig.Maximum != "" {
			dbPointer = dbPointer.Where("left(name, 1) <= ?", strings.ToUpper(requestBody.OrderConfig.Maximum))
		}
		dbPointer = dbPointer.Order("left(name,1) " + orderDirection)
	}

	//** ---------- Paramatere configuratino for DB query ends---------- **//
	err := utils.Pagination(requestBody, dbPointer).Order("id").
		Find(&authors).
		Error

	if err != nil {
		log.Printf("Got error when querying DB in GetAuthors: %s", err)
		errResponse := structs.ErrorResponse{
			Message: utils.InternalServerError,
		}
		return events.APIGatewayProxyResponse{
			Body:       errResponse.ToString(),
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	if len(authors) == 0 {
		errResponse := structs.ErrorResponse{
			Message:    "There are no authors matching your parameters",
			StatusCode: http.StatusOK,
		}
		return events.APIGatewayProxyResponse{
			Body:       errResponse.ToString(),
			StatusCode: http.StatusOK,
		}, nil
	}

	//Update popularity in background! TODO: Put into its own Lambda function
	go requestHandler.AuthorsAppearInSearchCountIncrement(authors)

	authorsAPI := structs.ConvertToAuthorsAPIModel(authors)
	out, _ := json.Marshal(authorsAPI)
	return events.APIGatewayProxyResponse{
		Body:       string(out),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(theReqHandler.handler)
}
