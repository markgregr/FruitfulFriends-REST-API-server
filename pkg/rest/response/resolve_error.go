package response

import (
	"errors"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExist          = errors.New("user already exists")
)

func ResolveError(err error) Error {
	switch {
	case errors.Is(err, ErrUserExist):
		return NewUserExistError()
	case errors.Is(err, ErrInvalidCredentials):
		return NewInvalidCredentialsError()
	default:
		return NewInternalError()
	}
	return NewInternalError()
}
