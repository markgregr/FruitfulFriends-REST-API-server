package response

import (
	"fmt"
	"net/http"
)

const (
	GeneralErrorKey = "_"
)

type InternalError struct {
	BaseError
}

func (e *InternalError) PublicMessage() string {
	return "internal error"
}

func (e *InternalError) GetHTTPStatus() int {
	return http.StatusInternalServerError
}

func NewInternalError() *InternalError {
	err := &InternalError{}
	return err
}

type UserExistError struct {
	BaseError
}

func (e *UserExistError) PublicMessage() string {
	return "user already exists"
}

func (e *UserExistError) GetHTTPStatus() int {
	return http.StatusConflict
}

func NewUserExistError() *UserExistError {
	return &UserExistError{}
}

type InvalidCredentialsError struct {
	BaseError
}

func (e *InvalidCredentialsError) PublicMessage() string {
	return "invalid credentials"
}

func (e *InvalidCredentialsError) GetHTTPStatus() int {
	return http.StatusForbidden
}

func NewInvalidCredentialsError() *InvalidCredentialsError {
	return &InvalidCredentialsError{}
}

type ErrorMessage struct {
	Code    ErrCode `json:"code"`
	Message string  `json:"message"`
}

func (e ErrorMessage) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

type ValidationError struct {
	BaseError
	details map[string]ErrorMessage
}

func (e *ValidationError) PublicMessage() string {
	return "validation error"
}

func (e *ValidationError) GetHTTPStatus() int {
	return http.StatusBadRequest
}

func (e *ValidationError) Errors() map[string]ErrorMessage {
	return e.details
}

func (e *ValidationError) SetError(key string, code ErrCode, message string) {
	e.details[key] = ErrorMessage{
		Code:    code,
		Message: message,
	}
}

func NewValidationError(details ...map[string]ErrorMessage) *ValidationError {
	err := &ValidationError{
		details: make(map[string]ErrorMessage),
	}

	if len(details) > 0 {
		for key, description := range details[0] {
			err.SetError(key, description.Code, description.Message)
		}
	}

	return err
}
