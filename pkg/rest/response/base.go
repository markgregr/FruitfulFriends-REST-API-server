package response

import "net/http"

type Error interface {
	PublicMessage() string
	GetHTTPStatus() int
}

type BaseError struct{}

func (e *BaseError) GetHTTPStatus() int {
	return http.StatusInternalServerError
}
