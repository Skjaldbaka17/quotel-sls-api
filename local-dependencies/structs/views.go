package structs

type SearchViewDBModel struct {
	AuthorId    int    `json:"author_id"`
	Name        string `json:"name"`
	QuoteId     int    `json:"quote_id" `
	Quote       string `json:"quote"`
	IsIcelandic bool   `json:"is_icelandic"`
	QuoteCount  int    `json:"quote_count"`
	AuthorCount int    `json:"author_count"`
}

type SearchViewAPIModel struct {
	// The author's id
	//Unique: true
	//example: 24952
	AuthorId int `json:"authorId"`
	// Name of author
	//example: Muhammad Ali
	Name string `json:"name"`
	// The quote's id
	//Unique: true
	//example: 582676
	QuoteId int `json:"quoteId" `
	// The quote
	//example: Float like a butterfly, sting like a bee.
	Quote string `json:"quote"`
	// Whether or not this quote is in Icelandic or not
	// example: false
	IsIcelandic bool `json:"isIcelandic"`
	//swagger:ignore
	QuoteCount int `json:"quoteCount"`
	//swagger:ignore
	AuthorCount int `json:"authorCount"`
}

func (dbModel *SearchViewDBModel) ConvertToAPIModel() SearchViewAPIModel {
	return SearchViewAPIModel(*dbModel)
}

func (apiModel *SearchViewAPIModel) ConvertToDBModel() SearchViewDBModel {
	return SearchViewDBModel(*apiModel)
}

func ConvertToSearchViewsAPIModel(views []SearchViewDBModel) []SearchViewAPIModel {
	viewsAPI := []SearchViewAPIModel{}
	for _, view := range views {
		viewsAPI = append(viewsAPI, SearchViewAPIModel(view))
	}
	return viewsAPI
}

func ConvertToSearchViewsDBModel(views []SearchViewAPIModel) []SearchViewDBModel {
	viewsDB := []SearchViewDBModel{}
	for _, view := range views {
		viewsDB = append(viewsDB, SearchViewDBModel(view))
	}
	return viewsDB
}
