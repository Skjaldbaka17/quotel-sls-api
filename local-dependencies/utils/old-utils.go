package utils

// //ValidateUserRequestBody takes in the request and validates all the input fields, returns an error with reason for validation-failure
// //if validation fails.
// //TODO: Make validation better! i.e. make it "real"
// func GetUserRequestBody(rw http.ResponseWriter, r *http.Request, requestBody *structs.UserApiModel) error {
// 	//Save the state back into the body for later use (Especially useful for getting the AOD/QOD because if the AOD has not been set a random AOD is set and the function called again)
// 	buf, _ := ioutil.ReadAll(r.Body)
// 	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
// 	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

// 	//Save the state back into the body for later use (Especially useful for getting the AOD/QOD because if the AOD has not been set a random AOD is set and the function called again)
// 	r.Body = rdr2
// 	err := json.NewDecoder(rdr1).Decode(&requestBody)

// 	if err != nil {
// 		log.Printf("Got error when decoding: %s", err)
// 		err = errors.New("request body is not structured correctly. Please refer to the /docs page for information on how to structure the request body")
// 		rw.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
// 		return err
// 	}

// 	return nil
// }

// func ValidateUserInformation(rw http.ResponseWriter, r *http.Request, requestBody *structs.UserApiModel) error {
// 	//TODO: Add email validation
// 	if requestBody.Email == "" {
// 		err := errors.New("email should not be empty")
// 		rw.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
// 		return err
// 	}

// 	if requestBody.Name == "" {
// 		err := errors.New("name should not be empty")
// 		rw.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
// 		return err
// 	}

// 	if requestBody.Password == "" {
// 		err := errors.New("password should not be empty")
// 		rw.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
// 		return err
// 	}

// 	if len(requestBody.Password) < 8 {
// 		err := errors.New("password should be at least 8 characters long")
// 		rw.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
// 		return err
// 	}

// 	if requestBody.PasswordConfirmation == "" {
// 		err := errors.New("password confirmation should not be empty")
// 		rw.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
// 		return err
// 	}

// 	if requestBody.PasswordConfirmation != requestBody.Password {
// 		err := errors.New("passwords do not match")
// 		rw.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
// 		return err
// 	}
// 	return nil
// }

// // Check whether user has GOD-tier permissions
// func AuthorizeGODApiKey(rw http.ResponseWriter, r *http.Request) error {
// 	var requestBody structs.Request
// 	if err, _ := getBody(rw, r, &requestBody); err != nil {
// 		return err
// 	}

// 	var user structs.UserDBModel
// 	if err := Db.Table("users").Where("api_key = ?", requestBody.ApiKey).First(&user).Error; err != nil {
// 		log.Printf("error when searching for user with the given api key in AuthorIzeGOD (api key validation): %s", err)
// 		rw.WriteHeader(http.StatusUnauthorized)
// 		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: "You need special privileges to access this route."})
// 		return err
// 	}

// 	if user.Tier != TIERS[len(TIERS)-1] {
// 		err := errors.New("you do not have the authorization to perform this action. Is your name Bassi Maraj? This is not meant for you... Sorry for the inconvenience")
// 		rw.WriteHeader(http.StatusUnauthorized)
// 		json.NewEncoder(rw).Encode(structs.ErrorResponse{Message: err.Error()})
// 		return err
// 	}

// 	return nil
// }
