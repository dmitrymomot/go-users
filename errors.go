package users

import (
	"errors"
)

// Predefined users package errors
var (
	ErrNotFound              = errors.New("user not found")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrTakenEmail            = errors.New("email is already taken")
	ErrEmailMissed           = errors.New("email is missed")
	ErrNotExistedUser        = errors.New("could not update not existed user")
	ErrCouldNotSetPassword   = errors.New("could not set password")
	ErrInvalidToken          = errors.New("invalid token string")
	ErrCouldNotGenerateToken = errors.New("could not generate new token")
)
