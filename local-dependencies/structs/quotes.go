package structs

type QuoteDBModel struct {
	Id          int    `json:"id,omitempty"`
	AuthorId    int    `json:"author_id,omitempty"`
	Quote       string `json:"quote,omitempty"`
	Count       int    `json:"count,omitempty"`
	IsIcelandic string `json:"is_icelandic,omitempty"`
}

type QuoteAPIModel struct {
	Id          int    `json:"id,omitempty"`
	AuthorId    int    `json:"authorId,omitempty"`
	Quote       string `json:"quote,omitempty"`
	Count       int    `json:"count,omitempty"`
	IsIcelandic string `json:"isIcelandic,omitempty"`
}

func (dbModel *QuoteDBModel) ConvertToAPIModel() QuoteAPIModel {
	return QuoteAPIModel(*dbModel)
}

func (apiModel *QuoteAPIModel) ConvertToDBModel() QuoteDBModel {
	return QuoteDBModel(*apiModel)
}

func ConvertToQuotesAPIModel(authors []QuoteDBModel) []QuoteAPIModel {
	authorsAPI := []QuoteAPIModel{}
	for _, author := range authors {
		authorsAPI = append(authorsAPI, QuoteAPIModel(author))
	}
	return authorsAPI
}

func ConvertToQuotesDBModel(authors []QuoteAPIModel) []QuoteDBModel {
	authorsDB := []QuoteDBModel{}
	for _, author := range authors {
		authorsDB = append(authorsDB, QuoteDBModel(author))
	}
	return authorsDB
}

type QodViewDBModel struct {
	QuoteId     int    `json:"quote_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Quote       string `json:"quote,omitempty"`
	AuthorId    int    `json:"author_id,omitempty"`
	IsIcelandic bool   `json:"is_icelandic,omitempty"`
	Date        string `json:"date,omitempty"`
}

type QodViewAPIModel struct {
	// The quote's id
	// example: 582676
	QuoteId int `json:"quoteId,omitempty"`
	// Name of the author
	// example: Muhammad Ali
	Name string `json:"name,omitempty"`
	// The quote for the day
	// example: Float like a butterfly, sting like a bee
	Quote string `json:"quote,omitempty"`
	// The author's id
	AuthorId int `json:"authorId,omitempty"`
	// Whether the quote is in icelandic
	// false
	IsIcelandic bool `json:"isIcelandic,omitempty"`
	// The date when this quote was the quote of the day
	// example: 2021-06-12T00:00:00Z
	Date string `json:"date,omitempty"`
}

func (dbModel *QodViewDBModel) ConvertToAPIModel() QodViewAPIModel {
	return QodViewAPIModel(*dbModel)
}

func (apiModel *QodViewAPIModel) ConvertToDBModel() QodViewDBModel {
	return QodViewDBModel(*apiModel)
}

func ConvertToQodViewsAPIModel(authors []QodViewDBModel) []QodViewAPIModel {
	authorsAPI := []QodViewAPIModel{}
	for _, author := range authors {
		authorsAPI = append(authorsAPI, QodViewAPIModel(author))
	}
	return authorsAPI
}

func ConvertToQodViewsDBModel(authors []QodViewAPIModel) []QodViewDBModel {
	authorsDB := []QodViewDBModel{}
	for _, author := range authors {
		authorsDB = append(authorsDB, QodViewDBModel(author))
	}
	return authorsDB
}

type Qod struct {
	// the date for which this quote is the QOD, if left empty this quote is today's QOD.
	//
	// Example: 12-22-2020
	Date string `json:"date,omitempty"`
	// The id of the quote to be set as this dates QOD
	//
	// Example: 1
	Id int `json:"id,omitempty"`
	// The language of the QOD
	// Example: icelandic
	Language string `json:"language,omitempty"`
}
