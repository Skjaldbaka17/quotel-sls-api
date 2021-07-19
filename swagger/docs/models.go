package docs

// swagger:model OrderConfiguration
type orderConfigListAuthorsModel struct {
	// What to order by, 'alphabetical', 'popularity','nrOfQuotes','dateOfBirth', 'dateOfDeath' or 'age'
	// example: popularity
	OrderBy string `json:"orderBy"`
	// Where to start the ordering (if empty it starts from beginning, for example start at 'A' for alphabetical ascending order).
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

// swagger:model quotesResponse
type baseQuotesResponseModel struct {
	// The author's id
	//Unique: true
	//example: 24952
	Authorid int `json:"authorid"`
	// Name of author
	//example: Muhammad Ali
	Name string `json:"name"`
	// The quote's id
	//Unique: true
	//example: 582676
	Quoteid int `json:"quoteid" `
	// The quote
	//example: Float like a butterfly, sting like a bee.
	Quote string `json:"quote"`
	// Whether or not this quote is in Icelandic or not
	// example: false
	Isicelandic bool `json:"isicelandic"`
}

// swagger:model qodResponseModel
type qodResponseModel struct {
	// The author's id
	//Unique: true
	//example: 24952
	Authorid int `json:"authorid"`
	// Name of the author
	//example: Muhammad Ali
	Name string `json:"name"`
	// The date when this author was the author of the day
	// example: 2021-06-12T00:00:00Z
	Date string `json:"date"`
	// The quote's id
	//Unique: true
	//example: 582676
	Quoteid int `json:"quoteid" `
	// The quote for the day
	// example: Float like a butterfly, sting like a bee
	Quote string `json:"quote"`
	// Whether the quote is in icelandic
	// example: false
	Isicelandic bool `json:"isicelandic"`
}

// swagger:model OfTheDayModel
type ofTheDayModel struct {
	// The id of the author / quote
	// example: 1
	Id int `json:"id"`
	// The date when this author / quote should be 'of the day'
	// example: 2020-06-12
	Date string `json:"date"`
	// The language of this author / quote
	//
	// Default: English
	// Example: icelandic
	Language string `json:"language"`
}

// swagger:model OrderConfiguration
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

// swagger:model Time
type timeModel struct {
	//Model
	Born beforeAfterModel `json:"born"`
	//Model
	Died beforeAfterModel `json:"died"`
	// If true only return authors that are alive
	// Example: true
	// Default: false
	IsAlive bool `json:"isAlive"`
	// If true only return authors that are dead
	// Example: true
	// Default: false
	IsDead bool `json:"isDead"`
	// Model
	Age ageModel `json:"age"`
}

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
