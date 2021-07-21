package structs

type QuoteDBModel struct {
	Id          uint   `json:"id,omitempty"`
	AuthorId    uint   `json:"author_id,omitempty"`
	Quote       string `json:"quote,omitempty"`
	Count       int    `json:"count,omitempty"`
	IsIcelandic bool   `json:"is_icelandic,omitempty"`

	Nationality string `json:"nationality,omitempty"`
	Profession  string `json:"profession,omitempty"`
	BirthYear   int    `json:"birth_year,omitempty"`
	BirthMonth  string `json:"birth_month,omitempty"`
	BirthDate   int    `json:"birth_date,omitempty"`
	DeathYear   int    `json:"death_year,omitempty"`
	DeathMonth  string `json:"death_month,omitempty"`
	DeathDate   int    `json:"death_date,omitempty"`
	Name        string `json:"name,omitempty"`
	TopicName   string `json:"topic_name,omitempty"`
	TopicId     int    `json:"topic_id,omitempty"`
}

type QuoteAPIModel struct {
	QuoteId     uint   `json:"quoteId,omitempty"`
	AuthorId    uint   `json:"authorId,omitempty"` //Author_id
	Quote       string `json:"quote,omitempty"`
	Count       int    `json:"count,omitempty"`
	IsIcelandic bool   `json:"isIcelandic,omitempty"`

	Nationality string `json:"nationality,omitempty"`
	Profession  string `json:"profession,omitempty"`
	Born        string `json:"born,omitempty"`
	Died        string `json:"died,omitempty"`
	Name        string `json:"name,omitempty"`
	TopicName   string `json:"topicName,omitempty"`
	TopicId     int    `json:"topicId,omitempty"`
}

type QodDBModel struct {
	QuoteId  uint   `json:"quote_id,omitempty"`
	Quote    string `json:"quote,omitempty"`
	AuthorId uint   `json:"author_id,omitempty"`

	Name        string `json:"name,omitempty"`
	Nationality string `json:"nationality,omitempty"`
	Profession  string `json:"profession,omitempty"`
	BirthYear   int    `json:"birth_year,omitempty"`
	BirthMonth  string `json:"birth_month,omitempty"`
	BirthDate   int    `json:"birth_date,omitempty"`
	DeathYear   int    `json:"death_year,omitempty"`
	DeathMonth  string `json:"death_month,omitempty"`
	DeathDate   int    `json:"death_date,omitempty"`
	Date        string `json:"date,omitempty"`

	IsIcelandic bool   `json:"is_icelandic,omitempty"`
	TopicId     uint   `json:"topic_id,omitempty"`
	TopicName   string `json:"topic_name,omitempty"`
}

type QodAPIModel struct {
	QuoteId  uint   `json:"quoteId,omitempty"`
	Quote    string `json:"quote,omitempty"`
	AuthorId uint   `json:"authorId,omitempty"` //author_id

	Name        string `json:"name,omitempty"`
	Nationality string `json:"nationality,omitempty"`
	Profession  string `json:"profession,omitempty"`
	Born        string `json:"born,omitempty"`
	Died        string `json:"died,omitempty"`
	Date        string `json:"date,omitempty"`

	IsIcelandic bool   `json:"isIcelandic,omitempty"`
	TopicId     uint   `json:"topicId,omitempty"`
	TopicName   string `json:"topicName,omitempty"`
}

//------------------- STRUCT CONVERSIONS -------------------//

/* quoteDB to quoteAPI conversion */
func (dbModel *QuoteDBModel) ConvertToAPIModel() QuoteAPIModel {
	return QuoteAPIModel{
		QuoteId:     dbModel.Id,
		AuthorId:    dbModel.AuthorId,
		Quote:       dbModel.Quote,
		IsIcelandic: dbModel.IsIcelandic,
		Nationality: dbModel.Nationality,
		Profession:  dbModel.Profession,
		Born:        getDate(dbModel.BirthYear, dbModel.BirthMonth, dbModel.BirthDate),
		Died:        getDate(dbModel.DeathYear, dbModel.DeathMonth, dbModel.DeathDate),
		Name:        dbModel.Name,
		TopicName:   dbModel.TopicName,
		TopicId:     dbModel.TopicId,
	}
}

/* quoteDB to qodDB conversion */
func (dbModel *QuoteDBModel) ConvertToQODDBModel(date string) QodDBModel {
	return QodDBModel{
		QuoteId:     dbModel.Id,
		AuthorId:    dbModel.AuthorId,
		Quote:       dbModel.Quote,
		BirthYear:   dbModel.BirthYear,
		BirthMonth:  dbModel.BirthMonth,
		BirthDate:   dbModel.BirthDate,
		DeathYear:   dbModel.DeathYear,
		DeathMonth:  dbModel.DeathMonth,
		DeathDate:   dbModel.DeathDate,
		Nationality: dbModel.Nationality,
		Profession:  dbModel.Profession,
		Name:        dbModel.Name,
		Date:        date,
		IsIcelandic: dbModel.IsIcelandic,
		TopicId:     uint(dbModel.TopicId),
		TopicName:   dbModel.TopicName,
	}
}

/* qodDB to qodAPI conversion */
func (dbModel *QodDBModel) ConvertToAPIModel() QodAPIModel {
	return QodAPIModel{
		QuoteId:     dbModel.QuoteId,
		AuthorId:    dbModel.AuthorId,
		Quote:       dbModel.Quote,
		Nationality: dbModel.Nationality,
		Profession:  dbModel.Profession,
		Born:        getDate(dbModel.BirthYear, dbModel.BirthMonth, dbModel.BirthDate),
		Died:        getDate(dbModel.DeathYear, dbModel.DeathMonth, dbModel.DeathDate),
		Name:        dbModel.Name,
		Date:        dbModel.Date,
		IsIcelandic: dbModel.IsIcelandic,
		TopicId:     dbModel.TopicId,
		TopicName:   dbModel.TopicName,
	}
}

//------------------- SLICE CONVERSIONS -------------------//

/* quoteDBs to quoteAPIs conversion */
func ConvertToQuotesAPIModel(quotes []QuoteDBModel) []QuoteAPIModel {
	quotesAPI := []QuoteAPIModel{}
	for _, author := range quotes {
		quotesAPI = append(quotesAPI, author.ConvertToAPIModel())
	}
	return quotesAPI
}

/* qodDBs to qodAPIs conversion */
func ConvertToQodAPIModel(quotes []QodDBModel) []QodAPIModel {
	quotesAPI := []QodAPIModel{}
	for _, author := range quotes {
		quotesAPI = append(quotesAPI, author.ConvertToAPIModel())
	}
	return quotesAPI
}
