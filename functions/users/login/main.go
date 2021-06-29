package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/crypto/bcrypt"
)

type RequestHandler struct {
	utils.RequestHandler
}

var theReqHandler = RequestHandler{}

// swagger:route POST /users/login USERS Login
// Login to get the apiKey for the user
// responses:
//	200: userResponse
//  400: incorrectBodyStructureResponse
//  401: incorrectCredentialsResponse
//  500: internalServerErrorResponse

// Login handles post requests to login to a user and receive his ApiKey
func (requestHandler *RequestHandler) handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//Initialize DB if requestHandler.Db = nil
	if errResponse := requestHandler.InitializeDB(); errResponse != (structs.ErrorResponse{}) {
		return events.APIGatewayProxyResponse{
			Body:       errResponse.Message,
			StatusCode: errResponse.StatusCode,
		}, nil
	}

	requestBody, errResponse := utils.GetUserRequestBody(request)

	if errResponse != (structs.ErrorResponse{}) {
		return events.APIGatewayProxyResponse{
			Body:       errResponse.Message,
			StatusCode: errResponse.StatusCode,
		}, nil
	}

	var user structs.UserDBModel
	if err := requestHandler.Db.Table("users").Where("email = ?", requestBody.Email).First(&user).Error; err != nil {
		log.Printf("Got error when login/fetching user: %s", err)
		return events.APIGatewayProxyResponse{
			Body:       "No user with the given email address. Maybe try WindsOfWinterWillNeverBeFinished@WeAreAllSinnersBeforeTheSeven.com",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	//Compare passwords / Check correct password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(requestBody.Password)); err != nil {
		log.Printf("Got error when comparing passwords in login: %s", err)
		return events.APIGatewayProxyResponse{
			Body:       "Credentials not correct. Shame. Shame. Shame is the name of the game.",
			StatusCode: http.StatusUnauthorized,
		}, nil
	}
	out, _ := json.Marshal(structs.UserResponse{ApiKey: user.ApiKey})
	return events.APIGatewayProxyResponse{
		Body:       string(out),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(theReqHandler.handler)
}
