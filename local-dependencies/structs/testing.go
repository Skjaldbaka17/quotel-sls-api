package structs

type TestApiResponse struct {
	Id                  int    `json:"id"`
	Name                string `json:"name"`
	HasIcelandicQuotes  bool   `json:"hasIcelandicQuotes"`
	NrOfIcelandicQuotes int    `json:"nrOfIcelandicQuotes"`
	NrOfEnglishQuotes   int    `json:"nrOfEnglishQuotes"`
	Count               int    `json:"count"`
	AuthorId            int    `json:"authorId"`
	QuoteId             int    `json:"quoteId" `
	TopicId             int    `json:"topicId" `
	TopicName           string `json:"topicName" `
	Quote               string `json:"quote"`
	IsIcelandic         bool   `json:"isIcelandic"`
	Date                string `json:"date"`
	QuoteCount          int    `json:"quoteCount"`
	AuthorCount         int    `json:"authorCount"`
}
