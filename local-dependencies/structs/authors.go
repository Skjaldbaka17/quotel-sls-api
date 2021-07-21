package structs

import (
	"strconv"

	"gorm.io/gorm"
)

type AuthorDBModel struct {
	gorm.Model
	Nationality         string `json:"nationality,omitempty"`
	Profession          string `json:"profession,omitempty"`
	BirthYear           int    `json:"birth_year,omitempty"`
	BirthMonth          string `json:"birth_month,omitempty"`
	BirthDate           int    `json:"birth_date,omitempty"`
	DeathYear           int    `json:"death_year,omitempty"`
	DeathMonth          string `json:"death_month,omitempty"`
	DeathDate           int    `json:"death_date,omitempty"`
	Name                string `json:"name,omitempty"`
	HasIcelandicQuotes  bool   `json:"has_icelandic_quotes,omitempty"`
	NrOfIcelandicQuotes int    `json:"nr_of_icelandic_quotes,omitempty"`
	NrOfEnglishQuotes   int    `json:"nr_of_english_quotes,omitempty"`
	Count               int    `json:"count,omitempty"`
}

type AuthorAPIModel struct {
	AuthorId            uint   `json:"authorId,omitempty"`
	Name                string `json:"name,omitempty"`
	HasIcelandicQuotes  bool   `json:"hasIcelandicQuotes,omitempty"`
	NrOfIcelandicQuotes int    `json:"nrOfIcelandicQuotes,omitempty"`
	NrOfEnglishQuotes   int    `json:"nrOfEnglishQuotes,omitempty"`
	Count               int    `json:"count,omitempty"`
	BirthDate           string `json:"birthDate,omitempty"`
	DeathDate           string `json:"deathDate,omitempty"`
	Profession          string `json:"profession,omitempty"`
	Nationality         string `json:"nationality,omitempty"`
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
	AuthorId    uint   `json:"author_id,omitempty"`
	Date        string `json:"date,omitempty"`
	IsIcelandic bool   `json:"is_icelandic,omitempty"`
}

type AodAPIModel struct {
	Name        string `json:"name,omitempty"`
	BirthDate   string `json:"birthDate,omitempty"`
	DeathDate   string `json:"deathDate,omitempty"`
	Profession  string `json:"profession,omitempty"`
	Nationality string `json:"nationality,omitempty"`
	Date        string `json:"date,omitempty"`
	IsIcelandic bool   `json:"isIcelandic,omitempty"`
	AuthorId    uint   `json:"authorId,omitempty"`
}

//------------------- STRUCT CONVERSIONS -------------------//

/* Converting authorDB to authorAPI */
func (dbModel *AuthorDBModel) ConvertToAPIModel() AuthorAPIModel {
	return AuthorAPIModel{
		AuthorId:            dbModel.ID,
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

/* Converting aodDB to aodAPI */
func (dbModel *AodDBModel) ConvertToAPIModel() AodAPIModel {
	return AodAPIModel{
		AuthorId:    dbModel.AuthorId,
		Name:        dbModel.Name,
		BirthDate:   getDate(dbModel.BirthYear, dbModel.BirthMonth, dbModel.BirthDate),
		DeathDate:   getDate(dbModel.DeathYear, dbModel.DeathMonth, dbModel.DeathDate),
		Profession:  dbModel.Profession,
		Nationality: dbModel.Nationality,
		Date:        dbModel.Date,
		IsIcelandic: dbModel.IsIcelandic,
	}
}

/*
	Converting authorDB to aodDB
	Used for inserting new aods
*/
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
		AuthorId:    dbModel.ID,
	}
	return model
}

//------------------- SLICE CONVERSIONS -------------------//

/* Converting aodDBs to aodAPIs */
func ConvertToAodAPIModel(authors []AodDBModel) []AodAPIModel {
	authorsAPI := []AodAPIModel{}
	for _, author := range authors {
		authorsAPI = append(authorsAPI, author.ConvertToAPIModel())
	}
	return authorsAPI

}

/* Converting authorDBs to authorAPIs */
func ConvertToAuthorsAPIModel(authors []AuthorDBModel) []AuthorAPIModel {
	authorsAPI := []AuthorAPIModel{}
	for _, author := range authors {
		authorsAPI = append(authorsAPI, author.ConvertToAPIModel())
	}
	return authorsAPI
}

//------------------- USEFUL FUNCTIONS -------------------//

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
