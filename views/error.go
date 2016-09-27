package views

import "net/http"

// ErrorMessage represents a HTTP error with status code and message.
type ErrorMessage struct {
	// Code is HTTP status code.
	Code    int
	Message string
}

// Error is for implementation of error interface.
func (e *ErrorMessage) Error() string {
	return e.Message
}

// ErrorMapper maps error to http error.
type ErrorMapper interface {
	MapError(http.ResponseWriter, *http.Request, error)
}

// errorMapper is a default implementation of ErrorMapper interface.
type errorMapper struct {
}

func newErrorMapper() *errorMapper {
	return &errorMapper{}
}

func (h *errorMapper) MapError(w http.ResponseWriter, _ *http.Request, err error) {
	switch v := err.(type) {
	case *ErrorMessage:
		http.Error(w, v.Message, v.Code)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
