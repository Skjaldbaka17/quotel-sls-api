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

func (dbModel *TopicDBModel) ConvertToAPIModel() TopicAPIModel {
	return TopicAPIModel(*dbModel)
}

func (apiModel *TopicAPIModel) ConvertToDBModel() TopicDBModel {
	return TopicDBModel(*apiModel)
}

func ConvertToTopicsAPIModel(authors []TopicDBModel) []TopicAPIModel {
	authorsAPI := []TopicAPIModel{}
	for _, author := range authors {
		authorsAPI = append(authorsAPI, TopicAPIModel(author))
	}
	return authorsAPI
}

func ConvertToTopicsDBModel(authors []TopicAPIModel) []TopicDBModel {
	authorsDB := []TopicDBModel{}
	for _, author := range authors {
		authorsDB = append(authorsDB, TopicDBModel(author))
	}
	return authorsDB
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
	IsIcelandic bool `json:"isIcelandi,omitempty"`
	// The topic's name (if topic id / name not supplied this will return empty string "")
	// example: inspirational
	TopicName string `json:"topicName,omitempty"`
	// The topic's id (if topic id / name not supplied this will return a zero id)
	// example: 10
	TopicId int `json:"topicId,omitempty"`
}

func (dbModel *TopicViewDBModel) ConvertToAPIModel() TopicViewAPIModel {
	return TopicViewAPIModel(*dbModel)
}

func (apiModel *TopicViewAPIModel) ConvertToDBModel() TopicViewDBModel {
	return TopicViewDBModel(*apiModel)
}

func ConvertToTopicViewsAPIModel(views []TopicViewDBModel) []TopicViewAPIModel {
	viewsAPI := []TopicViewAPIModel{}
	for _, view := range views {
		viewsAPI = append(viewsAPI, TopicViewAPIModel(view))
	}
	return viewsAPI
}

func ConvertToTopicViewsDBModel(views []TopicViewAPIModel) []TopicViewDBModel {
	viewsDB := []TopicViewDBModel{}
	for _, view := range views {
		viewsDB = append(viewsDB, TopicViewDBModel(view))
	}
	return viewsDB
}
