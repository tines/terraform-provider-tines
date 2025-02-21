package tines_cli

import (
	"fmt"
	"strings"
)

const (
	errEmptyApiKey         = "API Token must not be empty"
	errEmptyTenant         = "Tines Tenant must not be empty"
	errInternalServerError = "internal server error"
	errDoRequestError      = "error while attempting to make the HTTP request"
	errUnmarshalError      = "error unmarshalling the JSON response"
	errReadBodyError       = "error reading the HTTP response body bytes"
	errParseError          = "error parsing the input"
)

type ErrorType string

const (
	ErrorTypeRequest        ErrorType = "request"
	ErrorTypeAuthentication ErrorType = "authentication"
	ErrorTypeAuthorization  ErrorType = "authorization"
	ErrorTypeNotFound       ErrorType = "not_found"
	ErrorTypeRateLimit      ErrorType = "rate_limit"
	ErrorTypeServer         ErrorType = "server"
)

type Error struct {
	Type       ErrorType      `json:"type,omitempty"`
	StatusCode int            `json:"status_code,omitempty"`
	Errors     []ErrorMessage `json:"errors,omitempty"`
}

type ErrorMessage struct {
	Message string `json:"message,omitempty"`
	Details string `json:"details,omitempty"`
}

func (e Error) Error() string {
	var errString string
	errMessages := []string{}
	for _, err := range e.Errors {
		msg := fmt.Sprintf("%s: %s", err.Message, err.Details)
		errMessages = append(errMessages, msg)
	}
	errString = fmt.Sprintf("%d error(s) occurred: %s", len(e.Errors), strings.Join(errMessages, ", "))
	return errString
}
