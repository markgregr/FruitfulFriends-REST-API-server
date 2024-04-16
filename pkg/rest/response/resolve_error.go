package response

import (
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidCredentials = "invalid credentials"
	ErrUserExist          = "user already exists"
)

func ResolveError(err error) Error {
	st, ok := status.FromError(err)
	if !ok {
		return NewInternalError()
	}

	switch st.Message() {
	case ErrUserExist:
		return NewUserExistError()
	case ErrInvalidCredentials:
		return NewInvalidCredentialsError()
	default:
		return NewInternalError()
	}
}
