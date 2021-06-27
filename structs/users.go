package structs

type UserApiModel struct {
	ApiKey               string `json:"apiKey"`
	Id                   int    `json:"id"`
	Email                string `json:"email"`
	Name                 string `json:"name"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
	Tier                 string `json:"tier"`
}

type UserDBModel struct {
	Id           int    `json:"id"`
	ApiKey       string `json:"api_key"`
	Message      string `json:"message"`
	PasswordHash string `json:"password_hash"`
	Tier         string `json:"tier"`
	Name         string `json:"name"`
	Email        string `json:"email"`
}

type RequestEvent struct {
	Id          int    `json:"id"`
	UserId      int    `json:"user_id"`
	Route       string `json:"route"`
	RequestBody string `json:"request_body"`
	ApiKey      string `json:"api_key"`
}

type ErrorEvent struct {
	Id           int    `json:"id"`
	UserId       int    `json:"user_id"`
	Route        string `json:"route"`
	RequestBody  string `json:"request_body"`
	ErrorMessage string `json:"error_message"`
	ExtraInfo    string `json:"extra_info"`
}

type UserResponse struct {
	// The user id
	// example: 1
	Id int `json:"id"`
	// The api-key that the user should send to get access to the api
	// example: 1d8db1d2-6f5b-4254-8b74-44f5e5229add
	ApiKey string `json:"apiKey"`
}
