package domain

import "errors"

var (
	UserDoesntExistError = errors.New("the user doesn't exist")
	UserWithEmailAlreadyExists = errors.New("the email is already in use")
	InvalidParametersError = errors.New("Invalid parameter")
)