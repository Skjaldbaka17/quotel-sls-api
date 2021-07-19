package structs

import (
	"strconv"

	"gorm.io/gorm"
)

type AuthorDBModel struct {
	gorm.Model
	Nationality string `json:"nationality,omitempty"`
	Profession  string `json:"profession,omitempty"`
	BirthYear   int    `json:"birth_year,omitempty"`
	BirthMonth  string `json:"birth_month,omitempty"`
	BirthDate   int    `json:"birth_date,omitempty"`
	DeathYear   int    `json:"death_year,omitempty"`
	DeathMonth  string `json:"death_month,omitempty"`
	DeathDate   int    `json:"death_date,omitempty"`

	Name                string `json:"name,omitempty"`
	HasIcelandicQuotes  bool   `json:"has_icelandic_quotes,omitempty"`
	NrOfIcelandicQuotes int    `json:"nr_of_icelandic_quotes,omitempty"`
	NrOfEnglishQuotes   int    `json:"nr_of_english_quotes,omitempty"`
	Count               int    `json:"count,omitempty"`
}

type AuthorAPIModel struct {
	// The author's id
	// unique: true
	// example: 24952
	Id uint `json:"id,omitempty"`
	// Name of the author
	// example: Muhammad Ali
	Name string `json:"name,omitempty"`
	// Whether or not this author has some icelandic quotes
	// example: true
	HasIcelandicQuotes bool `json:"hasIcelandicQuotes,omitempty"`
	// How many quotes in Icelandic this author has
	// example: 6
	NrOfIcelandicQuotes int `json:"nrOfIcelandicQuotes,omitempty"`
	// How many quotes in English this author has
	// example: 78
	NrOfEnglishQuotes int `json:"nrOfEnglishQuotes,omitempty"`
	// The popularity index of the author
	// example: 1111
	Count int `json:"count,omitempty"`
	// The date of birth for this author - some may only have a birth year or a birth year and month
	// example: 1998-June-16
	BirthDate string `json:"birthDate,omitempty"`
	// The date of death for this author - some may only have a birth year or a birth year and month
	// example: 1998-June-16
	DeathDate string `json:"deathDate,omitempty"`
	// The author's profession
	// example: Musician
	Profession string `json:"profession,omitempty"`
	// The author's nationality
	// example: American
	Nationality string `json:"nationality,omitempty"`
}

func (dbModel *AuthorDBModel) ConvertToAPIModel() AuthorAPIModel {
	return AuthorAPIModel{
		Id:                  dbModel.ID,
		Name:                dbModel.Name,
		HasIcelandicQuotes:  dbModel.NrOfIcelandicQuotes > 0,
		NrOfIcelandicQuotes: dbModel.NrOfIcelandicQuotes,
		NrOfEnglishQuotes:   dbModel.NrOfEnglishQuotes,
		Count:               dbModel.Count,
		BirthDate:           getDate(dbModel.BirthYear, dbModel.BirthMonth, dbModel.BirthDate),
		DeathDate:           getDate(dbModel.DeathYear, dbModel.DeathMonth, dbModel.DeathDate),
		Profession:          dbModel.Profession,
		Nationality:         dbModel.Nationality,
	}
}

func getDate(year int, month string, day int) string {
	var date string
	if year > 0 {
		date += strconv.Itoa(year)
		if month != "" {
			date += "-" + month
		}
		if day > 0 {
			date += "-" + strconv.Itoa(day)
		}
	} else if month != "" {
		date += month
		if day > 0 {
			date += "-" + strconv.Itoa(day)
		}
	}
	return date
}

func ConvertToAuthorsAPIModel(authors []AuthorDBModel) []AuthorAPIModel {
	authorsAPI := []AuthorAPIModel{}
	for _, author := range authors {
		authorsAPI = append(authorsAPI, author.ConvertToAPIModel())
	}
	return authorsAPI
}

type AodRequest struct {
	// the date for which this author is the AOD, if left empty this quote is today's AOD.
	//
	// Example: 12-22-2020
	Date string `json:"date,omitempty"`
	// The id of the author to be set as this dates AOD
	//
	// Example: 1
	Id int `json:"id,omitempty"`
	// The language of the AOD
	// Example: icelandic
	Language string `json:"language,omitempty"`
}

type AodDBModel struct {
	gorm.Model
	Nationality string `json:"nationality,omitempty"`
	Profession  string `json:"profession,omitempty"`
	BirthYear   int    `json:"birth_year,omitempty"`
	BirthMonth  string `json:"birth_month,omitempty"`
	BirthDate   int    `json:"birth_date,omitempty"`
	DeathYear   int    `json:"death_year,omitempty"`
	DeathMonth  string `json:"death_month,omitempty"`
	DeathDate   int    `json:"death_date,omitempty"`
	Name        string `json:"name,omitempty"`

	Date string `json:"date,omitempty"`

	IsIcelandic bool `json:"is_icelandic,OMITEMPTY"`
}

type AodAPIModel struct {
	// The author's id
	// unique: true
	// example: 24952
	Id uint `json:"id,omitempty"`
	// Name of the author
	// example: Muhammad Ali
	Name string `json:"name,omitempty"`
	// The date of birth for this author - some may only have a birth year or a birth year and month
	// example: 1998-June-16
	BirthDate string `json:"birthDate,omitempty"`
	// The date of death for this author - some may only have a birth year or a birth year and month
	// example: 1998-June-16
	DeathDate string `json:"deathDate,omitempty"`
	// The author's profession
	// example: Musician
	Profession string `json:"profession,omitempty"`
	// The author's nationality
	// example: American
	Nationality string `json:"nationality,omitempty"`
	// The date when this author was the author of the day
	// example: 2021-06-12T00:00:00Z
	Date        string `json:"date,omitempty"`
	IsIcelandic bool   `json:"isIcelandic,OMITEMPTY"`
	TopicId     uint   `json:"topicId,OMITEMPTY"`
	TopicName   string `json:"topicName,OMITEMPTY"`
}

func (dbModel *AodDBModel) ConvertToAPIModel() AodAPIModel {
	return AodAPIModel{
		Id:          dbModel.ID,
		Name:        dbModel.Name,
		BirthDate:   getDate(dbModel.BirthYear, dbModel.BirthMonth, dbModel.BirthDate),
		DeathDate:   getDate(dbModel.DeathYear, dbModel.DeathMonth, dbModel.DeathDate),
		Profession:  dbModel.Profession,
		Nationality: dbModel.Nationality,
		Date:        dbModel.Date,
		IsIcelandic: dbModel.IsIcelandic,
	}
}

func ConvertToAodAPIModel(authors []AodDBModel) []AodAPIModel {
	authorsAPI := []AodAPIModel{}
	for _, author := range authors {
		authorsAPI = append(authorsAPI, author.ConvertToAPIModel())
	}
	return authorsAPI
}

func (dbModel *AuthorDBModel) ConvertToAODDBModel(date string, isIcelandic bool) AodDBModel {
	model := AodDBModel{
		Name:        dbModel.Name,
		BirthYear:   dbModel.BirthYear,
		BirthMonth:  dbModel.BirthMonth,
		BirthDate:   dbModel.BirthDate,
		DeathYear:   dbModel.DeathYear,
		DeathMonth:  dbModel.DeathMonth,
		DeathDate:   dbModel.DeathDate,
		Profession:  dbModel.Profession,
		Nationality: dbModel.Nationality,
		Date:        date,
		IsIcelandic: isIcelandic,
	}
	model.ID = dbModel.ID
	return model
}
