package structs

type AuthorDBModel struct {
	Id                  int    `json:"id"`
	Name                string `json:"name"`
	HasIcelandicQuotes  bool   `json:"has_icelandic_quotes"`
	NrOfIcelandicQuotes int    `json:"nr_of_icelandic_quotes"`
	NrOfEnglishQuotes   int    `json:"nr_of_english_quotes"`
	Count               int    `json:"count"`
}

type AuthorAPIModel struct {
	// The author's id
	// unique: true
	// example: 24952
	Id int `json:"id"`
	// Name of the author
	// example: Muhammad Ali
	Name string `json:"name"`
	// Whether or not this author has some icelandic quotes
	// example: true
	HasIcelandicQuotes bool `json:"hasIcelandicQuotes"`
	// How many quotes in Icelandic this author has
	// example: 6
	NrOfIcelandicQuotes int `json:"nrOfIcelandicQuotes"`
	// How many quotes in English this author has
	// example: 78
	NrOfEnglishQuotes int `json:"nrOfEnglishQuotes"`
	// The popularity index of the author
	// example: 1111
	Count int `json:"count"`
}

func (dbModel *AuthorDBModel) ConvertToAPIModel() AuthorAPIModel {
	return AuthorAPIModel(*dbModel)
}

func (apiModel *AuthorAPIModel) ConvertToDBModel() AuthorDBModel {
	return AuthorDBModel(*apiModel)
}

func ConvertToAuthorsAPIModel(authors []AuthorDBModel) []AuthorAPIModel {
	authorsAPI := []AuthorAPIModel{}
	for _, author := range authors {
		authorsAPI = append(authorsAPI, AuthorAPIModel(author))
	}
	return authorsAPI
}

func ConvertToAuthorsDBModel(authors []AuthorAPIModel) []AuthorDBModel {
	authorsDB := []AuthorDBModel{}
	for _, author := range authors {
		authorsDB = append(authorsDB, AuthorDBModel(author))
	}
	return authorsDB
}

type AodDBModel struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Date string `json:"date"`
}

type AodAPIModel struct {
	// The author's id
	// example: 24952
	Id int `json:"id"`
	// The name of the author
	// example: Muhammad Ali
	Name string `json:"name"`
	// The date when this author was the author of the day
	// example: 2021-06-12T00:00:00Z
	Date string `json:"date"`
}
