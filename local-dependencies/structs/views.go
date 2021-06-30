package structs

import "encoding/json"

type SearchViewDBModel struct {
	AuthorId    int    `json:"author_id,omitempty"`
	Name        string `json:"name,omitempty"`
	QuoteId     int    `json:"quote_id,omitempty" `
	Quote       string `json:"quote,omitempty"`
	IsIcelandic bool   `json:"is_icelandic,omitempty"`
	QuoteCount  int    `json:"quote_count,omitempty"`
	AuthorCount int    `json:"author_count,omitempty"`
}

type SearchViewAPIModel struct {
	// The author's id
	//Unique: true
	//example: 24952
	AuthorId int `json:"authorId,omitempty"`
	// Name of author
	//example: Muhammad Ali
	Name string `json:"name,omitempty"`
	// The quote's id
	//Unique: true
	//example: 582676
	QuoteId int `json:"quoteId,omitempty" `
	// The quote
	//example: Float like a butterfly, sting like a bee.
	Quote string `json:"quote,omitempty"`
	// Whether or not this quote is in Icelandic or not
	// example: false
	IsIcelandic bool `json:"isIcelandic,omitempty"`
	//swagger:ignore
	QuoteCount int `json:"quoteCount,omitempty"`
	//swagger:ignore
	AuthorCount int `json:"authorCount,omitempty"`
}

func (view *SearchViewAPIModel) ToString() string {
	out, _ := json.Marshal(view)
	return string(out)
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
