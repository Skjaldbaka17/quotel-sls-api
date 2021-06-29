package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/structs"
	"github.com/Skjaldbaka17/quotel-sls-api/local-dependencies/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RequestHandler struct {
	utils.RequestHandler
}

var theReqHandler = RequestHandler{}

// swagger:route POST /users/signup USERS SignUp
// Create A user to get a free ApiKey
// responses:
//	200: userResponse
//  400: incorrectBodyStructureResponse
//  500: internalServerErrorResponse

// CreateUsers handles post requests to create a user and an accompanying ApiKey
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

	errResponse = utils.ValidateUserInformation(&requestBody)

	if errResponse != (structs.ErrorResponse{}) {
		return events.APIGatewayProxyResponse{
			Body:       errResponse.Message,
			StatusCode: errResponse.StatusCode,
		}, nil
	}

	uuid, _ := uuid.NewRandom()
	apiKey := uuid.String()
	passHash, _ := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
	requestBody.Tier = utils.TIERS[0]
	user := structs.UserDBModel{Name: requestBody.Name, ApiKey: apiKey, Tier: requestBody.Tier, Email: requestBody.Email, PasswordHash: string(passHash)}

	result := requestHandler.Db.Table("users").Select("name", "api_key", "tier", "email", "password_hash").Create(&user)

	//Error handle
	if result.Error != nil {
		m1 := regexp.MustCompile(`duplicate key value violates unique constraint "users_email_key"`)
		if m1.Match([]byte(result.Error.Error())) {
			log.Printf("Got error when creating user, constraint error: %s", result.Error)
			return events.APIGatewayProxyResponse{
				Body:       "This email is taken.",
				StatusCode: http.StatusBadRequest,
			}, nil
		}
		log.Printf("Got error when creating user: %s", result.Error)
		return events.APIGatewayProxyResponse{
			Body:       utils.InternalServerError,
			StatusCode: http.StatusInternalServerError,
		}, nil
	} else if user.Id <= 0 {
		log.Printf("Got no id when creating user: %s", result.Error)
		return events.APIGatewayProxyResponse{
			Body:       utils.InternalServerError,
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	out, _ := json.Marshal(structs.UserResponse{Id: user.Id, ApiKey: user.ApiKey})
	return events.APIGatewayProxyResponse{
		Body:       string(out),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(theReqHandler.handler)
}
