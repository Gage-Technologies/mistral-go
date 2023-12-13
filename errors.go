package mistral

import (
	"fmt"
)

// MistralError is the base error type for all Mistral errors.
type MistralError struct {
	Message string
}

func (e *MistralError) Error() string {
	return e.Message
}

// MistralAPIError is returned when the API responds with an error message.
type MistralAPIError struct {
	MistralError
	HTTPStatus int
	Headers    map[string][]string
}

func NewMistralAPIError(message string, httpStatus int, headers map[string][]string) *MistralAPIError {
	return &MistralAPIError{
		MistralError: MistralError{Message: message},
		HTTPStatus:       httpStatus,
		Headers:          headers,
	}
}

func (e *MistralAPIError) Error() string {
	return fmt.Sprintf("%s (HTTP status: %d)", e.Message, e.HTTPStatus)
}

// MistralConnectionError is returned when the SDK cannot reach the API server for any reason.
type MistralConnectionError struct {
	MistralError
}

func NewMistralConnectionError(message string) *MistralConnectionError {
	return &MistralConnectionError{
		MistralError: MistralError{Message: message},
	}
}