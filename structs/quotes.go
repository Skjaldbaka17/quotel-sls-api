package structs

type QuoteDBModel struct {
	Id          int    `json:"id"`
	AuthorId    int    `json:"author_id"`
	Quote       string `json:"quote"`
	Count       int    `json:"count"`
	IsIcelandic string `json:"is_icelandic"`
}

type QuoteAPIModel struct {
	Id          int    `json:"id"`
	AuthorId    int    `json:"authorId"`
	Quote       string `json:"quote"`
	Count       int    `json:"count"`
	IsIcelandic string `json:"isIcelandic"`
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
	QuoteId     int    `json:"quote_id"`
	Name        string `json:"name"`
	Quote       string `json:"quote"`
	AuthorId    int    `json:"author_id"`
	IsIcelandic bool   `json:"is_icelandic"`
	Date        string `json:"date"`
}

type QodViewAPIModel struct {
	// The quote's id
	// example: 582676
	QuoteId int `json:"quoteId"`
	// Name of the author
	// example: Muhammad Ali
	Name string `json:"name"`
	// The quote for the day
	// example: Float like a butterfly, sting like a bee
	Quote string `json:"quote"`
	// The author's id
	AuthorId int `json:"authorId"`
	// Whether the quote is in icelandic
	// false
	IsIcelandic bool `json:"isIcelandic"`
	// The date when this quote was the quote of the day
	// example: 2021-06-12T00:00:00Z
	Date string `json:"date"`
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
