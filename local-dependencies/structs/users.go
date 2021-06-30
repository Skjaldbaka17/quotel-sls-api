package structs

type UserApiModel struct {
	ApiKey               string `json:"apiKey,omitempty"`
	Id                   int    `json:"id,omitempty"`
	Email                string `json:"email,omitempty"`
	Name                 string `json:"name,omitempty"`
	Password             string `json:"password,omitempty"`
	PasswordConfirmation string `json:"passwordConfirmation,omitempty"`
	Tier                 string `json:"tier,omitempty"`
}

type UserDBModel struct {
	Id           int    `json:"id,omitempty"`
	ApiKey       string `json:"api_key,omitempty"`
	Message      string `json:"message,omitempty"`
	PasswordHash string `json:"password_hash,omitempty"`
	Tier         string `json:"tier,omitempty"`
	Name         string `json:"name,omitempty"`
	Email        string `json:"email,omitempty"`
}

type RequestEvent struct {
	Id          int    `json:"id,omitempty"`
	UserId      int    `json:"user_id,omitempty"`
	Route       string `json:"route,omitempty"`
	RequestBody string `json:"request_body,omitempty"`
	ApiKey      string `json:"api_key,omitempty"`
	Request     string `json:"request,omitempty"`
}

type ErrorEvent struct {
	Id           int    `json:"id,omitempty"`
	UserId       int    `json:"user_id,omitempty"`
	Route        string `json:"route,omitempty"`
	RequestBody  string `json:"request_body,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	ExtraInfo    string `json:"extra_info,omitempty"`
}

type UserResponse struct {
	// The user id
	// example: 1
	Id int `json:"id,omitempty"`
	// The api-key that the user should send to get access to the api
	// example: 1d8db1d2-6f5b-4254-8b74-44f5e5229add
	ApiKey string `json:"apiKey,omitempty"`
}
