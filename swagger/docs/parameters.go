package docs

// ------------------------------------ AUTHORS ------------------------------------ //

// swagger:parameters GetAODHistory
type historyAODWrapper struct {
	// The structure of the request for getting the history of AODs
	// in: body
	Body []struct {
		// Get the history of the AODS for the given language ("icelandic" or "english")
		//
		// Default: English
		// Example: icelandic
		Language string `json:"language"`
		// The earliest date to return. All authors between minimum and today will be returned.
		// Example: 2020-12-21
		Minimum string `json:"minimum"`
	}
}

// swagger:parameters GetAuthorOfTheDay
type authorOfTheDayWrapper struct {
	// The structure of the request for getting the author / quote of the day
	// in: body
	Body struct {
		// Get the author / quote of the day for the given language ("icelandic" or "english")
		//
		// Default: English
		// Example: English
		// Get the author / quote of the day for the given topic by supplying its topicId
		//
		// Example: 1
		Language string `json:"language"`
	}
}

// swagger:parameters GetAuthors
type getAuthorsWrapper struct {
	// The structure of the request for getting authors by their ids
	// in: body
	// required: true
	Body struct {
		// A list of the authors's ids that you want
		//
		// Required: true
		// Example: [ 29333, 19161]
		Ids []int `json:"ids"`
	}
}

// swagger:parameters ListAuthors
type authorsListWrapper struct {
	// The structure of the request for getting a list of authors
	// in: body
	Body struct {
		// The Response is paged. This parameter controls the number of Authors to be returned on each "page"
		//
		// Maximum: 200
		// Minimum: 1
		// Default: 25
		// Example: 30
		PageSize int `json:"pageSize"`
		// Response is paged. This parameter controls the page you are asking for, starts with 0.
		//
		// Minimum: 0
		// Example: 0
		Page int `json:"page"`
		// Only return authors that have quotes in the given language ("english" or "icelandic") if left empty then no constraint
		// is set on the quotes' language. Note: if also ordering by nrOfQuotes and this parameter is set then only the amount of
		// quotes the author has in the given language counts towards the final ordering.
		// Example: English
		Language string `json:"language"`
		// Only return the author's in any of the given professions
		// Example: Designer
		Professions []string `json:"professions"`
		// Only return the author's with any of the given nationalities
		// Example: American
		Nationalities []string `json:"nationalities"`
		//Model
		OrderConfig orderConfigListAuthorsModel `json:"orderConfig"`
		//Model
		Time timeModel `json:"time"`
	}
}

// swagger:parameters GetRandomAuthor
type randomAuthorWrapper struct {
	// The structure of the request for getting a random author
	// in: body
	Body struct {
		// The random author must have quotes in the given language ("english" or "icelandic") if left empty then no
		// constraint on language is set
		//
		// Example: English
		Language string `json:"language"`
		// How many of the author's quotes, maximum, should this request return
		//
		// Example: 10
		// Maximum: 50
		// default: 1
		MaxQuotes int `json:"maxQuotes"`
	}
}

// swagger:parameters SearchAuthorsByString
type getSearchAuthorsByStringWrapper struct {
	// The structure of the request for searching quotes/authors
	// in: body
	// required: true
	Body struct {
		// The string to be used in the search
		//
		// Required: true
		// Example: Ali Muhammad
		SearchString string `json:"searchString"`
		// The number of quotes to be returned on each "page"
		//
		// Maximum: 200
		// Minimum: 1
		// Default: 25
		// Example: 30
		PageSize int `json:"pageSize"`
		// The page you are asking for, starts with 0.
		//
		// Minimum: 0
		// Example: 0
		Page int `json:"page"`
		// The particular language that the quote should be in
		// example: English
		Language string `json:"language"`
	}
}

// ------------------------------------ QUOTES ------------------------------------ //

// swagger:parameters GetQODHistory
type historyQODWrapper struct {
	// The structure of the request for getting the history of QODs
	// in: body
	Body []struct {
		// Get the history of the QODs for the given language ("icelandic" or "english")
		//
		// Default: English
		// Example: icelandic
		Language string `json:"language"`
		// The earliest date to return. All quotes between minimum and today will be returned.
		// Example: 2020-12-21
		Minimum string `json:"minimum"`
		// Get the QOD for the specified topic
		// Example: 1
		TopicId int `json:"topicId"`
	}
}

// swagger:parameters GetQuoteOfTheDay
type quoteOfTheDayWrapper struct {
	// The structure of the request for getting the quote of the day
	// in: body
	Body struct {
		// Get the quote of the day for the given language ("icelandic" or "english")
		//
		// Default: English
		// Example: English
		Language string `json:"language"`
		// Get the quote of the day for the given topic by supplying its topicId
		//
		// Example: 100
		TopicId int `json:"topicId"`
	}
}

// swagger:parameters GetQuotes
type getQuotesByWrapper struct {
	// The structure of the request to get quotes. There are two ways to use this route. 1. Send the ids of the quotes to be
	// retrieved or 2. send the id of the author of the quotes you want (if you use this option the quotes are paginated)
	// in: body
	// required: true
	Body struct {
		// The list of quotes's ids you want
		//
		// Example: [582676,443976]
		Ids []int `json:"ids"`
		// The id of the author of the quotes you want.
		// Example: 24952
		AuthorId int `json:"authorId"`
		// If using authorId the response is paged. This parameter controls the number of Authors to be returned on each "page"
		//
		// Maximum: 200
		// Minimum: 1
		// Default: 25
		// Example: 30
		PageSize int `json:"pageSize"`
		// If using authorId the response is paged. This parameter controls the page you are asking for, starts with 0.
		//
		// Minimum: 0
		// Example: 0
		Page int `json:"page"`
	}
}

// swagger:parameters GetQuotesList
type quotesListWrapper struct {
	// The structure of the request for getting a list of quotes
	// in: body
	Body struct {
		// Response is paged. This parameter controls the number of Quotes to be returned on each "page"
		//
		// Maximum: 200
		// Minimum: 1
		// Default: 25
		// Example: 30
		PageSize int `json:"pageSize"`
		// Response is paged. This parameter controls the page you are asking for, starts with 0.
		//
		// Minimum: 0
		// Example: 0
		Page int `json:"page"`
		// Only return quotes that have quotes in the given language ("english" or "icelandic") if left empty then no constraint
		// is set on quotes' language.
		// Example: English
		Language string `json:"language"`
		//Model
		OrderConfig orderConfigListQuotesModel `json:"orderConfig"`
	}
}

// swagger:parameters GetRandomQuote
type getRandomQuoteResponseWrapper struct {
	// The structure of the request for a random quote
	// in: body
	Body struct {
		// The random quote returned must be in the given language
		//
		// Example: English
		Language string `json:"language"`
		// The random quote returned must contain a match with the searchstring
		//
		// Example: float
		SearchString string `json:"searchString"`
		// The random quote returned must be a part of one of the topics with id in the topicIds array
		//
		// Example: [10,11]
		TopicIds []int `json:"topicIds"`
		// The random quote returned must be from the author with the given authorId
		//
		//example: 29333
		Authorid int `json:"authorId"`
	}
}

// ------------------------------------ SEARCH ------------------------------------ //

// swagger:parameters SearchByString SearchQuotesByString
type getSearchByStringWrapper struct {
	// The structure of the request for searching quotes/authors
	// in: body
	// required: true
	Body struct {
		// The string to be used in the search
		//
		// Required: true
		// Example: sting like butterfly
		SearchString string `json:"searchString"`
		// The number of quotes to be returned on each "page"
		//
		// Maximum: 200
		// Minimum: 1
		// Default: 25
		// Example: 30
		PageSize int `json:"pageSize"`
		// The page you are asking for, starts with 0.
		//
		// Minimum: 0
		// Example: 0
		Page int `json:"page"`
		// The particular language that the quote should be in
		// example: English
		Language string `json:"language"`
		// Should search in the specified topics for the searchString
		//
		// Example: [10]
		TopicIds []int `json:"topicIds"`
	}
}

// ------------------------------------ TOPICS ------------------------------------ //

// swagger:parameters GetTopic
type quotesFromTopicWrapper struct {
	// The structure of the request for listing topics
	// in: body
	Body struct {
		// Name of the topic, if left empty then the id is used
		//
		// required: false
		// Example: Motivational
		Topic string `json:"topic"`
		// The topic's id, if left empty then the topic name is used
		//
		// Example: 10
		Id int `json:"id"`
		// The number of quotes to be returned on each "page"
		//
		// Maximum: 200
		// Minimum: 1
		// Default: 25
		// Example: 30
		PageSize int `json:"pageSize"`
		// The page you are asking for, starts with 0.
		//
		// Minimum: 0
		// Example: 0
		Page int `json:"page"`
	}
}

// swagger:parameters GetTopics
type listTopicsWrapper struct {
	// The structure of the request for listing topics
	// in: body
	Body struct {
		// The language of the topics. If left empty all topics from all languages are returned
		//
		// Example: English
		Language string `json:"language"`
	}
}
