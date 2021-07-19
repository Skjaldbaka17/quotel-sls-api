package structs

type QuoteDBModel struct {
	Id          uint   `json:"id,OMITEMPTY"`
	AuthorId    uint   `json:"author_id,OMITEMPTY"`
	Quote       string `json:"quote,OMITEMPTY"`
	Count       int    `json:"count,OMITEMPTY"`
	IsIcelandic bool   `json:"is_icelandic,OMITEMPTY"`

	Nationality string `json:"nationality,OMITEMPTY"`
	Profession  string `json:"profession,OMITEMPTY"`
	BirthYear   int    `json:"birth_year,OMITEMPTY"`
	BirthMonth  string `json:"birth_month,OMITEMPTY"`
	BirthDate   int    `json:"birth_date,OMITEMPTY"`
	DeathYear   int    `json:"death_year,OMITEMPTY"`
	DeathMonth  string `json:"death_month,OMITEMPTY"`
	DeathDate   int    `json:"death_date,OMITEMPTY"`
	Name        string `json:"name,OMITEMPTY"`
	TopicName   string `json:"topic_name,OMITEMPTY"`
	TopicId     int    `json:"topic_id,OMITEMPTY"`
}

type QuoteAPIModel struct {
	QuoteId     uint   `json:"quoteId,OMITEMPTY"`
	AuthorId    uint   `json:"authorId,OMITEMPTY"` //Author_id
	Quote       string `json:"quote,OMITEMPTY"`
	Count       int    `json:"count,OMITEMPTY"`
	IsIcelandic bool   `json:"isIcelandic,OMITEMPTY"`

	Nationality string `json:"nationality,OMITEMPTY"`
	Profession  string `json:"profession,OMITEMPTY"`
	BirthDate   string `json:"birthDate,OMITEMPTY"`
	DeathDate   string `json:"deathDate,OMITEMPTY"`
	Name        string `json:"name,OMITEMPTY"`
	TopicName   string `json:"topicName,OMITEMPTY"`
	TopicId     int    `json:"topicId,OMITEMPTY"`
}

type QodDBModel struct {
	QuoteId  uint   `json:"quote_id,OMITEMPTY"`
	Quote    string `json:"quote,OMITEMPTY"`
	AuthorId uint   `json:"author_id,OMITEMPTY"`

	Name        string `json:"name,OMITEMPTY"`
	Nationality string `json:"nationality,OMITEMPTY"`
	Profession  string `json:"profession,OMITEMPTY"`
	BirthYear   int    `json:"birth_year,OMITEMPTY"`
	BirthMonth  string `json:"birth_month,OMITEMPTY"`
	BirthDate   int    `json:"birth_date,OMITEMPTY"`
	DeathYear   int    `json:"death_year,OMITEMPTY"`
	DeathMonth  string `json:"death_month,OMITEMPTY"`
	DeathDate   int    `json:"death_date,OMITEMPTY"`
	Date        string `json:"date,OMITEMPTY"`

	IsIcelandic bool   `json:"is_icelandic,OMITEMPTY"`
	TopicId     uint   `json:"topic_id,OMITEMPTY"`
	TopicName   string `json:"topic_name,OMITEMPTY"`
}

type QodAPIModel struct {
	QuoteId  uint   `json:"quoteId,OMITEMPTY"`
	Quote    string `json:"quote,OMITEMPTY"`
	AuthorId uint   `json:"authorId,OMITEMPTY"` //author_id

	Name        string `json:"name,OMITEMPTY"`
	Nationality string `json:"nationality,OMITEMPTY"`
	Profession  string `json:"profession,OMITEMPTY"`
	BirthDate   string `json:"birthDate,OMITEMPTY"`
	DeathDate   string `json:"deathDate,OMITEMPTY"`
	Date        string `json:"date,OMITEMPTY"`

	IsIcelandic bool   `json:"isIcelandic,OMITEMPTY"`
	TopicId     uint   `json:"topicId,OMITEMPTY"`
	TopicName   string `json:"topicName,OMITEMPTY"`
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
		BirthDate:   getDate(dbModel.BirthYear, dbModel.BirthMonth, dbModel.BirthDate),
		DeathDate:   getDate(dbModel.DeathYear, dbModel.DeathMonth, dbModel.DeathDate),
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
		BirthDate:   getDate(dbModel.BirthYear, dbModel.BirthMonth, dbModel.BirthDate),
		DeathDate:   getDate(dbModel.DeathYear, dbModel.DeathMonth, dbModel.DeathDate),
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
