package structs

type TestApiResponse struct {
	Id                  int    `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	HasIcelandicQuotes  bool   `json:"hasIcelandicQuotes,omitempty"`
	NrOfIcelandicQuotes int    `json:"nrOfIcelandicQuotes,omitempty"`
	NrOfEnglishQuotes   int    `json:"nrOfEnglishQuotes,omitempty"`
	Count               int    `json:"count,omitempty"`
	AuthorId            int    `json:"authorId,omitempty"`
	QuoteId             int    `json:"quoteId,omitempty" `
	TopicId             int    `json:"topicId,omitempty" `
	TopicName           string `json:"topicName,omitempty" `
	Quote               string `json:"quote,omitempty"`
	IsIcelandic         bool   `json:"isIcelandic,omitempty"`
	Date                string `json:"date,omitempty"`
	QuoteCount          int    `json:"quoteCount,omitempty"`
	AuthorCount         int    `json:"authorCount,omitempty"`
}
