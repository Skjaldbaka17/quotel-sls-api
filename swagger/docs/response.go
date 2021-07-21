package docs

// ------------------------------------ AUTHORS ------------------------------------ //

// Data structure representing the error response to internal server errors
// swagger:response internalServerErrorResponse
type internalServerErrorResponseWrapper struct {
	// The structure of the error response
	// in: body
	Body struct {
		// The error message describing what happened
		// Example: Please try again later.
		Message string `json:"message"`
		// The http status code for this error
		// Example: 500
		StatusCode int `json:"statusCode"`
	}
}

// Data structure representing the error response to a wrongly structured request body
// swagger:response incorrectBodyStructureResponse
type incorrectBodyStructureResponseWrapper struct {
	// The error response to a wrongly structured request body
	// in: body
	Body struct {
		// The error message
		// Example: request body is not structured correctly.
		Message string `json:"message"`
		// The http status code for this error
		// Example: 400
		StatusCode int `json:"statusCode"`
	}
}

// Data structure representing the response for the history of AODs
// swagger:response aodHistoryResponse
type aodHistoryResponseWrapper struct {
	// The response to the get history of AODs request
	// in: body
	Body []AodAPIModel
}

// Data structure representing the response for the author of the day
// swagger:response aodResponse
type aodResponseWrapper struct {
	// The response to the author of the day request
	// in: body
	Body AodAPIModel
}

// Data structure representing the response for authors
// swagger:response authorsResponse
type authorsResponseWrapper struct {
	// The structure of authors objects
	// in: body
	Body []AuthorAPIModel
}

// ------------------------------------ META ------------------------------------ //

// Data structure for the nationalities in the database
// swagger:response listNationalities
type listNationalitiesWrapper struct {
	// The nationalities supported by the api
	// in: body
	Body []struct {
		// The nationalities supported
		// example: ["American", "Italian"]
		Nationalities []string `json:"nationalities"`
	}
}

// Data structure for the professions in the database
// swagger:response listProfessions
type listProfessionsWrapper struct {
	// The professions supported by the api
	// in: body
	Body []struct {
		// The professions supported
		// example: ["Rapper", "Politician"]
		Professions []string `json:"professions"`
	}
}

// Data structure for supported languages information
// swagger:response listOfStrings
type listOfStringsWrapper struct {
	// The languages supported by the api
	// in: body
	Body []struct {
		// The languages supported
		// example: ["English", "Icelandic"]
		Languages []string `json:"languages"`
	}
}

// ------------------------------------ QUOTES ------------------------------------ //

// Data structure representing the response for quotes
// swagger:response quotesApiResponse
type quotesApiResponseWrapper struct {
	// Quotes response
	// in: body
	Body []QuoteAPIModel
}

// Data structure representing the response for the history of QODS
// swagger:response qodHistoryResponse
type qodHistoryResponseWrapper struct {
	// The response to the history of QODs
	// in: body
	Body []QodAPIModel
}

// Data structure representing the response for the QOD
// swagger:response qodResponse
type qodResponseWrapper struct {
	// The response to the quote of the day request
	// in: body
	Body QodAPIModel
}

// ------------------------------------ TOPICS ------------------------------------ //

// Data structure representing the response for a quote based on a particular topic
// swagger:response topicApiResponse
type topicApiResponseWrapper struct {
	// A quote struct
	// in: body
	Body TopicQuoteAPIModel
}

// Data structure representing a list response for topics
// swagger:response topicsResponse
type topicsResponseWrapper struct {
	// List of topics
	// in: body
	Body []TopicAPIModel
}

// ------------------------------------ OTHER ------------------------------------ //

// Data structure representing the error response to a not found error
// swagger:response notFoundResponse
type notFoundResponseWrapper struct {
	// The error response to a not found error
	// in: body
	Body struct {
		// The error message
		// Example: No quote exists that matches the given parameters.
		Message string `json:"message"`
	}
}
