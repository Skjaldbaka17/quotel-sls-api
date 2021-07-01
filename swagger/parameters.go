package docs

// swagger:parameters GetAuthors
type getAuthorsWrapper struct {
	// The structure of the request for getting authors by their ids
	// in: body
	// required: true
	Body struct {
		// A list of the authors's ids that you want
		//
		// Required: true
		// Example: [24952,19161]
		Ids []int `json:"ids"`
		// The api-key you use to access the api
		//
		// Required: true
		// Example: 91fd6d19-2c32-4081-8729-4d9786d43b95
		ApiKey string `json:"apiKey"`
	}
}

// swagger:parameters ListAuthors
type authorsListWrapper struct {
	// The structure of the request for getting a list of authors
	// in: body
	Body struct {
		// The api-key you use to access the api
		//
		// Required: true
		// Example: 91fd6d19-2c32-4081-8729-4d9786d43b95
		ApiKey string `json:"apiKey"`
		// Response is paged. This parameter controls the number of Authors to be returned on each "page"
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
		// is set on the quotes' language. Note if ordering by nrOfQuotes if this parameter is set then only the amount of
		// quotes the author has in the given language counts towards the final ordering.
		// Example: English
		Language string `json:"language"`
		//Model
		OrderConfig orderConfigListAuthorsModel `json:"orderConfig"`
	}
}

// swagger:parameters GetRandomAuthor
type randomAuthorWrapper struct {
	// The structure of the request for getting a random author
	// in: body
	Body struct {
		// The api-key you use to access the api
		//
		// Required: true
		// Example: 91fd6d19-2c32-4081-8729-4d9786d43b95
		ApiKey string `json:"apiKey"`
		// The random author must have quotes in the given language ("english" or "icelandic") if left empty then no
		// constraint on language is set
		//
		// Example: English
		Language string `json:"language"`
		// How many quotes, maximum, to be returned from this author
		//
		// Example: 10
		// Maximum: 50
		// default: 1
		MaxQuotes int `json:"maxQuotes"`
	}
}

// swagger:parameters GetAuthorOfTheDay GetQuoteOfTheDay
type ofTheDayWrapper struct {
	// The structure of the request for getting the author / quote of the day
	// in: body
	Body struct {
		// The api-key you use to access the api
		//
		// Required: true
		// Example: 91fd6d19-2c32-4081-8729-4d9786d43b95
		ApiKey string `json:"apiKey"`
		// Get the author / quote of the day for the given language ("icelandic" or "english")
		//
		// Default: English
		// Example: English
		Language string `json:"language"`
	}
}

// swagger:parameters GetAODHistory GetQODHistory
type historyAODWrapper struct {
	// The structure of the request for getting the history of AODs / QODs
	// in: body
	Body []struct {
		// The api-key you use to access the api
		//
		// Required: true
		// Example: 91fd6d19-2c32-4081-8729-4d9786d43b95
		ApiKey string `json:"apiKey"`
		// Get the history of the AODS / QODs for the given language ("icelandic" or "english")
		//
		// Default: English
		// Example: icelandic
		Language string `json:"language"`
		// The earliest date to return. All authors / quotes between minimum and today will be returned.
		// Example: 2020-12-21
		Minimum string `json:"minimum"`
	}
}

// swagger:parameters SetAuthorOfTheDay
type setAODWrapper struct {
	// The structure of the request for setting AODs
	// in: body
	Body []struct {
		// The api-key you use to access the api
		//
		// Required: true
		// Example: 91fd6d19-2c32-4081-8729-4d9786d43b95
		ApiKey string          `json:"apiKey"`
		Aods   []ofTheDayModel `json:"aods"`
	}
}

// swagger:parameters SetQuoteOfTheDay
type setQuoteOfTheDayWrapper struct {
	// The structure of the request for setting the QOD
	// in: body
	Body struct {
		// The api-key you use to access the api
		//
		// Required: true
		// Example: 91fd6d19-2c32-4081-8729-4d9786d43b95
		ApiKey string          `json:"apiKey"`
		Qods   []ofTheDayModel `json:"qods"`
	}
}

// swagger:parameters GetQuotes
type getQuotesByWrapper struct {
	// The structure of the request to get quotes. There are two ways to use this route. 1. Send the ids of the quotes to be
	// retrieved or 2. send the id of the author of the quotes you want (if you use this option the quotes are paginated)
	// in: body
	// required: true
	Body struct {
		// The api-key you use to access the api
		//
		// Required: true
		// Example: 91fd6d19-2c32-4081-8729-4d9786d43b95
		ApiKey string `json:"apiKey"`
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
		// The api-key you use to access the api
		//
		// Required: true
		// Example: 91fd6d19-2c32-4081-8729-4d9786d43b95
		ApiKey string `json:"apiKey"`
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
		// is set on the quotes' language.
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
		// The api-key you use to access the api
		//
		// Required: true
		// Example: 91fd6d19-2c32-4081-8729-4d9786d43b95
		ApiKey string `json:"apiKey"`
		// The random quote returned must be in the given language
		//
		// Example: English
		Language string `json:"language"`
		// The random quote returned must contain a match with the searchstring
		//
		// Example: float
		SearchString string `json:"searchString"`
		// The random quote returned must be a part of the topic with the given topicId
		//
		// Example: 10
		TopicId int `json:"topicId"`
		// The random quote returned must be from the author with the given authorId
		//
		//example: 24952
		Authorid int `json:"authorId"`
	}
}

// swagger:parameters SearchByString SearchQuotesByString
type getSearchByStringWrapper struct {
	// The structure of the request for searching quotes/authors
	// in: body
	// required: true
	Body struct {
		// The api-key you use to access the api
		//
		// Required: true
		// Example: 91fd6d19-2c32-4081-8729-4d9786d43b95
		ApiKey string `json:"apiKey"`
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
		// Should search in the specified topic for the searchString
		//
		// Example: 10
		TopicId int `json:"topicId"`
	}
}

// swagger:parameters SearchAuthorsByString
type getSearchAuthorsByStringWrapper struct {
	// The structure of the request for searching quotes/authors
	// in: body
	// required: true
	Body struct {
		// The api-key you use to access the api
		//
		// Required: true
		// Example: 91fd6d19-2c32-4081-8729-4d9786d43b95
		ApiKey string `json:"apiKey"`
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

// swagger:parameters GetTopic
type quotesFromTopicWrapper struct {
	// The structure of the request for listing topics
	// in: body
	Body struct {
		// The api-key you use to access the api
		//
		// Required: true
		// Example: 91fd6d19-2c32-4081-8729-4d9786d43b95
		ApiKey string `json:"apiKey"`
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

// swagger:parameters SignUp
type signUpParameterWrapper struct {
	// The structure of the sign up request
	// in: body
	Body struct {

		// The email for the created user (used for validation and password retrieval)
		// example: example@gmail.com
		Email string `json:"email"`
		// The name of the creator
		// example: Robert Huldars
		Name string `json:"name"`
		// The password for the user, used to get the ApiKey and look at stats (coming soon). Must be a relatively strong pass
		// example: 1234567890
		Password string `json:"password"`
		// The confirmation for the password. Must be same as "password"
		// example: 1234567890
		PasswordConfirmation string `json:"passwordConfirmation"`
	}
}

// swagger:parameters Login
type LoginParameterWrapper struct {
	// The structure of the Login request
	// in: body
	Body struct {
		// The email for the user you want to login
		// example: example@gmail.com
		Email string `json:"email"`
		// The password for the user
		// example: 1234567890
		Password string `json:"password"`
	}
}
