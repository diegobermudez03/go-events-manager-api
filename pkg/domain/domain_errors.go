package domain

import "errors"

var (
	//	AUTH ERRORS
	ErrUserDoesntExist 				= errors.New("the user doesn't exist")
	ErrUserWithEmailAlreadyExists 	= errors.New("the email is already in use")
	ErrInvalidParametersError 		= errors.New("invalid parameter")
	ErrInternal						= errors.New("internal error processing request")
	ErrIncorrectPassword			= errors.New("incorrect password")
	ErrExpiredSession				= errors.New("the session has expired")
	ErrUnathorized					= errors.New("unathorized")

	//	PARTICIPATIONS ERRORS
	ErrNoParticipationFound 		= errors.New("the user isn't a participant of the event")

	// INVITATION ERRORS
	ErrAlreadyInvited 				= errors.New("user already invited")

	//	EVENTS ERRORS
	ErrInvalidImage 				= errors.New("invalid image file")
	ErrEventDoesntExist 			= errors.New("event doesnt exist")

	// ROLES ERRORS
	ErrRoleDoesntExist 				= errors.New("role doesnt exist")

)