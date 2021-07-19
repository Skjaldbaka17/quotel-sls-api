package structs

import "encoding/json"

type Request struct {
	Ids           []int       `json:"ids,omitempty"`
	Id            int         `json:"id,omitempty"`
	Page          int         `json:"page,omitempty"`
	SearchString  string      `json:"searchString,omitempty"`
	PageSize      int         `json:"pageSize,omitempty"`
	Language      string      `json:"language,omitempty"`
	Topic         string      `json:"topic,omitempty"`
	AuthorId      int         `json:"authorId,omitempty"`
	QuoteId       int         `json:"quoteId,omitempty"`
	TopicId       int         `json:"topicId,omitempty"`
	TopicIds      []int       `json:"topicIds,omitempty"`
	MaxQuotes     int         `json:"maxQuotes,omitempty"`
	OrderConfig   OrderConfig `json:"orderConfig,omitempty"`
	Date          string      `json:"date,omitempty"`
	Minimum       string      `json:"minimum,omitempty"`
	Maximum       string      `json:"maximum,omitempty"`
	StopRecursion bool        `json:"stopRecursion,omitempty"`
	Professions   []string    `json:"professions"`
	Nationalities []string    `json:"nationalities"`
	Time          Time        `json:"time"`
}

type Time struct {
	Born    BeforeAfter `json:"born"`
	Died    BeforeAfter `json:"died"`
	IsAlive bool        `json:"isAlive"`
	IsDead  bool        `json:"isDead"`
	Age     Age         `json:"age"`
}

type Age struct {
	Exactly     int `json:"exactly"`
	OlderThan   int `json:"olderThan"`
	YoungerThan int `json:"youngerThan"`
}

type BeforeAfter struct {
	Before string `json:"before"`
	After  string `json:"after"`
	Year   int    `json:"year"`
	Month  string `json:"month"`
	Date   int    `json:"date"`
}

type OrderConfig struct {
	OrderBy string `json:"orderBy,omitempty"`
	Minimum string `json:"minimum,omitempty"`
	Maximum string `json:"maximum,omitempty"`
	Reverse bool   `json:"reverse,omitempty"`
}

type ErrorResponse struct {
	Message    string `json:"message,omitempty"`
	StatusCode int    `json:"status_code,omitempty"`
}

func (errorResponse *ErrorResponse) ToString() string {
	out, _ := json.Marshal(errorResponse)
	return string(out)
}
