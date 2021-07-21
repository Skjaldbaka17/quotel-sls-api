package structs

type TopicDBModel struct {
	Id          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	IsIcelandic bool   `json:"is_icelandic,omitempty"`
}

type TopicAPIModel struct {
	Id          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	IsIcelandic bool   `json:"isIcelandic,omitempty"`
}

type TopicViewDBModel struct {
	AuthorId    int    `json:"author_id,omitempty"`
	Name        string `json:"name,omitempty"`
	QuoteId     int    `json:"quote_id,omitempty" `
	Quote       string `json:"quote,omitempty"`
	IsIcelandic bool   `json:"is_icelandic,omitempty"`
	TopicName   string `json:"topic_name,omitempty"`
	TopicId     int    `json:"topic_id,omitempty"`
}

type TopicViewAPIModel struct {
	AuthorId    int    `json:"authorId,omitempty"`
	Name        string `json:"name,omitempty"`
	QuoteId     int    `json:"quoteId,omitempty" `
	Quote       string `json:"quote,omitempty"`
	IsIcelandic bool   `json:"isIcelandic,omitempty"`
	TopicName   string `json:"topicName,omitempty"`
	TopicId     int    `json:"topicId,omitempty"`
}

//------------------- SLICE CONVERSIONS -------------------//

func ConvertToTopicsAPIModel(authors []TopicDBModel) []TopicAPIModel {
	authorsAPI := []TopicAPIModel{}
	for _, author := range authors {
		authorsAPI = append(authorsAPI, TopicAPIModel(author))
	}
	return authorsAPI
}
