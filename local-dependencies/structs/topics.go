package structs

type TopicDBModel struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	IsIcelandic bool   `json:"is_icelandic"`
}

type TopicAPIModel struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	IsIcelandic bool   `json:"isIcelandic"`
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
	AuthorId    int    `json:"author_id"`
	Name        string `json:"name"`
	QuoteId     int    `json:"quote_id" `
	Quote       string `json:"quote"`
	IsIcelandic bool   `json:"is_icelandic"`
	TopicName   string `json:"topic_name"`
	TopicId     int    `json:"topic_id"`
}

type TopicViewAPIModel struct {
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
	// The topic's name (if topic id / name not supplied this will return empty string "")
	// example: inspirational
	TopicName string `json:"topicName"`
	// The topic's id (if topic id / name not supplied this will return a zero id)
	// example: 10
	TopicId int `json:"topicId"`
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
