package docs

import "github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"

// Data structure representing the response for authors
// swagger:response authorsResponse
type authorsResponseWrapper struct {
	// An authors response
	// in: body
	Body []structs.AuthorAPIModel
}

// Data structure representing the response for quotes
// swagger:response searchViewsResponse
type searchViewsResponseWrapper struct {
	// Quotes response
	// in: body
	Body []structs.SearchViewAPIModel
}

// Data structure representing the response for a quote
// swagger:response searchViewResponse
type searchViewResponseWrapper struct {
	// A quote struct
	// in: body
	Body structs.SearchViewAPIModel
}

// Data structure representing the response for a quote based on a particular topic
// swagger:response topicViewResponse
type topicViewResponseWrapper struct {
	// A quote struct
	// in: body
	Body structs.TopicViewAPIModel
}

// Data structure representing the response for a quote based on a particular topic
// swagger:response topicViewsResponse
type topicViewsResponseWrapper struct {
	// A quote struct
	// in: body
	Body []structs.TopicViewAPIModel
}

// Data structure representing the response for the author of the day
// swagger:response aodResponse
type aodResponseWrapper struct {
	// The response to the author of the day request
	// in: body
	Body structs.AodAPIModel
}

// Data structure representing the response for the history of AODs
// swagger:response aodHistoryResponse
type aodHistoryResponseWrapper struct {
	// The response to the history of AODs request
	// in: body
	Body []structs.AodAPIModel
}

// Data structure representing the response for the quote of the day
// swagger:response qodResponse
type qodResponseWrapper struct {
	// The response to the quote of the day request
	// in: body
	Body structs.QodViewAPIModel
}

// Data structure representing the response for the history of QODS
// swagger:response qodHistoryResponse
type qodHistoryResponseWrapper struct {
	// The response to the history of QODs
	// in: body
	Body []structs.QodViewAPIModel
}

// swagger:response successResponse
type successResponseWrapper struct {
	// The successful response to a successful setting of an asset
	// in: body
	Body struct {
		// Example: This request was a success
		Message string `json:"message"`
		// HTTP status code
		//
		// Example: 200
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
	}
}

// Data structure representing the error response to an internal server error
// swagger:response internalServerErrorResponse
type internalServerErrorResponseWrapper struct {
	// The error response to an internal server
	// in: body
	Body struct {
		// The error message
		// Example: Please try again later.
		Message string `json:"message"`
	}
}

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

// Data structure representing a list response for topics
// swagger:response topicsResponse
type topicsResponseWrapper struct {
	// List of topics
	// in: body
	Body []structs.TopicAPIModel
}

// Data structure representing a user response
// swagger:response userResponse
type userResponseWrapper struct {
	// The necessary data for the user to use the API
	// in: body
	Body structs.UserResponse
}

// Data structure representing the error response to an incorrect Credentials error
// swagger:response incorrectCredentialsResponse
type incorrectCredentialsResponseWrapper struct {
	// The error response to an unothorized access
	// in: body
	Body struct {
		// The error message
		// Example: Valar Dohaeris
		Message string `json:"message"`
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
