package structs

type Request struct {
	Ids          []int       `json:"ids,omitempty"`
	Id           int         `json:"id,omitempty"`
	Page         int         `json:"page,omitempty"`
	SearchString string      `json:"searchString,omitempty"`
	PageSize     int         `json:"pageSize,omitempty"`
	Language     string      `json:"language,omitempty"`
	Topic        string      `json:"topic,omitempty"`
	AuthorId     int         `json:"authorId"`
	QuoteId      int         `json:"quoteId"`
	TopicId      int         `json:"topicId"`
	MaxQuotes    int         `json:"maxQuotes"`
	OrderConfig  OrderConfig `json:"orderConfig"`
	Date         string      `json:"date"`
	Minimum      string      `json:"minimum"`
	Maximum      string      `json:"maximum"`
	Qods         []Qod       `json:"qods"`
	Aods         []Qod       `json:"aods"`
	ApiKey       string      `json:"apiKey"`
}

type Qod struct {
	// the date for which this quote is the QOD, if left empty this quote is today's QOD.
	//
	// Example: 12-22-2020
	Date string `json:"date"`
	// The id of the quote to be set as this dates QOD
	//
	// Example: 1
	Id int `json:"id"`
	// The language of the QOD
	// Example: icelandic
	Language string `json:"language"`
}

type OrderConfig struct {
	// What to order by, 'alphabetical', 'popularity' or 'nrOfQuotes'
	// example: popularity
	OrderBy string `json:"orderBy"`
	// Where to start the ordering (if empty it starts from beginning, for example start at 'A' for alphabetical ascending order)
	// example: F
	Minimum string `json:"minimum"`
	// Where to end the ordering (if empty it ends at the logical end, for example end at 'Z' for alphabetical ascending order)
	// example: Z
	Maximum string `json:"maximum"`
	// Whether to order the list in reverse or not (true is Descending and false is Ascending, false is default)
	// example: true
	Reverse bool `json:"reverse"`
}

type ErrorResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}
