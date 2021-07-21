package docs

// ------------------------------------ REQUEST ------------------------------------ //

// swagger:model Time
type timeModel struct {
	// The structure for getting a list of authors based on some time values. Birtdate, year of death, age etc.
	//Model
	Born beforeAfterModel `json:"born"`
	//Model
	Died beforeAfterModel `json:"died"`
	// If true only return authors that are known to be alive
	// Example: true
	// Default: false
	IsAlive bool `json:"isAlive"`
	// If true only return authors that are known to be dead
	// Example: true
	// Default: false
	IsDead bool `json:"isDead"`
	// Model
	Age ageModel `json:"age"`
}

// swagger:model ageModel
type ageModel struct {
	// Only return authors that are exactly the given age
	// Example: 25
	Exactly int `json:"exactly"`
	// Only return authors that are older than or of equal age as the given age
	// Example: 25
	OlderThan int `json:"olderThan"`
	// Only return authors that are younger than or of equal age as the given age
	// Example: 25
	YoungerThan int `json:"youngerThan"`
}

// swagger:model BeforeAfter
type beforeAfterModel struct {
	// The date the author should be born/or have died on or before (i.e. <= 1998-06-16)
	// Example: 1998-06-16
	Before string `json:"before"`
	// The date the author should be born/or have died on or after (i.e. >= 1998-06-16)
	// Example: 1998-06-16
	After string `json:"after"`
	// Only return authors that were born/died in the given year
	// Example: 1998
	Year int `json:"year"`
	// Only return authors that were born/died in the given month
	// Example: June
	Month string `json:"month"`
	// Only return authors that were born/died on the given date
	// Example: 16
	Date int `json:"date"`
}

// swagger:model ListAuthorsOrderConfiguration
type orderConfigListAuthorsModel struct {
	// What to order by: 	'alphabetical', 'popularity', 'nrOfQuotes', 'dateOfBirth', 'dateOfDeath' or 'age'
	// example: popularity
	OrderBy string `json:"orderBy"`
	// Where to start the ordering (if empty it starts from the beginning, for example start at 'A' for alphabetical ascending order).
	// Note this key is always a string, for example if ordering by nrOfQuotes (or popularity) of minimum 10 quotes you need to
	// set "minimum":"10" in the request body
	// example: 10
	Minimum string `json:"minimum"`
	// Where to end the ordering (if empty it ends at the logical end, for example end at 'Z' for alphabetical ascending order).
	// Note this key is always a string, for example if ordering by nrOfQuotes (or popularity) of maximum 11 quotes you need to
	// set "maximum":"11" in the request body
	// example: 11
	Maximum string `json:"maximum"`
	// Whether to order the list in reverse or not (true is Descending and false is Ascending, false is default)
	// example: true
	Reverse bool `json:"reverse"`
}

// swagger:model ListQuotesOrderConfiguration
type orderConfigListQuotesModel struct {
	// What to order by, 'quoteId', 'popularity' or 'length'
	// example: popularity
	OrderBy string `json:"orderBy"`
	// Where to start the ordering (if empty it starts from beginning, for example start at 1 for quoteid ascending order).
	// Note this key is always a string.
	// example: 10
	Minimum string `json:"minimum"`
	// Where to end the ordering (if empty it ends at the logical end, for example end at the highest quoteid for quoteid ascending order).
	// Note this key is always a string.
	// example: 11
	Maximum string `json:"maximum"`
	// Whether to order the list in reverse or not (true is Descending and false is Ascending, false is default)
	// example: true
	Reverse bool `json:"reverse"`
}

// ------------------------------------ AUTHORS ------------------------------------ //

//swagger:model aodModel
type AodAPIModel struct {
	// The name of the author
	// example: Muhammad Ali
	Name string `json:"name,omitempty"`
	// The date this author was born
	// example: 1942-January-17
	Born string `json:"born,omitempty"`
	// The date this author died
	// example: 2016-June-3
	Died string `json:"died,omitempty"`
	// The author's main profession
	// example: Boxer
	Profession string `json:"profession,omitempty"`
	// The author's nationality
	// example: American
	Nationality string `json:"nationality,omitempty"`
	// The date when this author was the author of the day
	// example: 2021-06-12T00:00:00Z
	Date string `json:"date,omitempty"`
	// Whether the author is icelandic or not
	// example: false
	IsIcelandic bool `json:"isIcelandic,omitempty"`
	// The author's id
	// example: 29333
	AuthorId uint `json:"authorId,omitempty"`
}

//swagger:model authorModel
type AuthorAPIModel struct {
	// The name of the author
	// example: Muhammad Ali
	Name string `json:"name,omitempty"`
	// The date this author was born
	// example: 1942-January-17
	Born string `json:"born,omitempty"`
	// The date this author died
	// example: 2016-June-3
	Died string `json:"died,omitempty"`
	// The author's main profession
	// example: Boxer
	Profession string `json:"profession,omitempty"`
	// The author's nationality
	// example: American
	Nationality string `json:"nationality,omitempty"`
	// The author's id
	// example: 29333
	AuthorId uint `json:"authorId,omitempty"`
	// Whether or not this author has some icelandic quotes
	// example: true
	HasIcelandicQuotes bool `json:"hasIcelandicQuotes,omitempty"`
	// How many quotes in Icelandic this author has
	// example: 6
	NrOfIcelandicQuotes int `json:"nrOfIcelandicQuotes,omitempty"`
	// How many quotes in English this author has
	// example: 114
	NrOfEnglishQuotes int `json:"nrOfEnglishQuotes,omitempty"`
	// The popularity index of the author
	// example: 1111
	Count int `json:"count,omitempty"`
}

// ------------------------------------ QUOTES ------------------------------------ //

// swagger:model QuoteAPIModel
type QuoteAPIModel struct {
	// The quote's id
	// example: 879890
	QuoteId uint `json:"quoteId,omitempty"`
	// The quote
	// example: I hated every minute of training, but I said, 'Don't quit. Suffer now and live the rest of your life as a champion.'
	Quote string `json:"quote,omitempty"`
	// Whether or not this quote is in icelandic
	// example: false
	IsIcelandic bool `json:"isIcelandic,omitempty"`
	// The name of the author
	// example: Muhammad Ali
	Name string `json:"name,omitempty"`
	// The date this author was born
	// example: 1942-January-17
	Born string `json:"born,omitempty"`
	// The date this author died
	// example: 2016-June-3
	Died string `json:"died,omitempty"`
	// The author's main profession
	// example: Boxer
	Profession string `json:"profession,omitempty"`
	// The author's nationality
	// example: American
	Nationality string `json:"nationality,omitempty"`
	// The popularity index of the quote
	// example: 1111
	Count int `json:"count,omitempty"`
	// The author's id
	// example: 29333
	AuthorId uint `json:"authorId,omitempty"`
}

// swagger:model QodPIModel
type QodAPIModel struct {
	// The quote's id
	// example: 879890
	QuoteId int `json:"quoteId,omitempty"`
	// The name of the author
	// example: Muhammad Ali
	Name string `json:"name,omitempty"`
	// The quote for the day
	// example: I hated every minute of training, but I said, 'Don't quit. Suffer now and live the rest of your life as a champion.'
	Quote string `json:"quote,omitempty"`
	// The author's id
	// example: 29333
	AuthorId int `json:"authorId,omitempty"`
	// Whether or not this quote is in icelandic
	// example: false
	IsIcelandic bool `json:"isIcelandic,omitempty"`
	// The date when this quote was the quote of the day
	// example: 2021-06-12
	Date string `json:"date,omitempty"`

	// The date this author was born
	// example: 1942-January-17
	Born string `json:"born,omitempty"`
	// The date this author died
	// example: 2016-June-3
	Died string `json:"died,omitempty"`
	// The author's main profession
	// example: Boxer
	Profession string `json:"profession,omitempty"`
	// The author's nationality
	// example: American
	Nationality string `json:"nationality,omitempty"`
	// The topic's Id that this QOD is quote of the day of
	// Example: 100
	TopicId uint `json:"topicId,omitempty"`
	// The topic's name that this QOD is quote of the day of
	// Example: Wisdom
	TopicName string `json:"topicName,omitempty"`
}

// ------------------------------------ TOPICS ------------------------------------ //

//swagger:model topicModel
type TopicAPIModel struct {
	Id          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	IsIcelandic bool   `json:"isIcelandic,omitempty"`
}

//swagger:model topicQuoteModel
type TopicQuoteAPIModel struct {
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

	Count int `json:"count,omitempty"`

	// The date this author was born
	// example: 1942-January-17
	Born string `json:"born,omitempty"`
	// The date this author died
	// example: 2016-June-3
	Died string `json:"died,omitempty"`
	// The author's main profession
	// example: Boxer
	Profession string `json:"profession,omitempty"`
	// The author's nationality
	// example: American
	Nationality string `json:"nationality,omitempty"`
	// The topic's Id that this QOD is quote of the day of
	// Example: 100
	TopicId uint `json:"topicId,omitempty"`
	// The topic's name that this QOD is quote of the day of
	// Example: Wisdom
	TopicName string `json:"topicName,omitempty"`
}
