package response

import (
	"errors"
	"google.golang.org/grpc"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExist          = errors.New("user already exists")
)

func ResolveError(err error) Error {
	switch {
	case errors.Is(errors.New(grpc.ErrorDesc(err)), ErrUserExist):
		return NewUserExistError()
	case errors.Is(errors.New(grpc.ErrorDesc(err)), ErrInvalidCredentials):
		return NewInvalidCredentialsError()
	default:
		return NewInternalError()
	}
	return NewInternalError()
}
