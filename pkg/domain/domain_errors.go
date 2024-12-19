package domain

import "errors"

var (
	ErrUserDoesntExist = errors.New("the user doesn't exist")
	ErrUserWithEmailAlreadyExists = errors.New("the email is already in use")
	ErrInvalidParametersError = errors.New("invalid parameter")
	ErrInternal	= errors.New("internal error processing request")
	ErrIncorrectPassword = errors.New("incorrect password")
)