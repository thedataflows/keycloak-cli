package admin

import (
	"errors"
	"fmt"
	"net/http"

	admininternal "github.com/thedataflows/keycloak-cli/pkg/admin/internal"
)

type ErrorKind string

const (
	ErrorUnauthorized ErrorKind = "unauthorized"
	ErrorNotFound     ErrorKind = "not-found"
	ErrorConflict     ErrorKind = "conflict"
	ErrorValidation   ErrorKind = "validation"
	ErrorTransport    ErrorKind = "transport"
)

type Error struct {
	Kind       ErrorKind
	Operation  string
	Resource   string
	StatusCode int
	Err        error
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	prefix := e.Operation
	if e.Resource != "" {
		if prefix != "" {
			prefix += " " + e.Resource
		} else {
			prefix = e.Resource
		}
	}
	message := e.message()
	if prefix == "" {
		return message
	}
	if message == "" {
		return prefix
	}
	return fmt.Sprintf("%s: %s", prefix, message)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func classifyError(err error, statusCode int, operation, resource string) error {
	if err == nil {
		return nil
	}

	var httpErr *admininternal.HTTPError
	if errors.As(err, &httpErr) {
		statusCode = httpErr.StatusCode
	}

	return &Error{
		Kind:       kindFromStatus(statusCode),
		Operation:  operation,
		Resource:   resource,
		StatusCode: statusCode,
		Err:        err,
	}
}

func (e *Error) message() string {
	var httpErr *admininternal.HTTPError
	if errors.As(e.Err, &httpErr) {
		return formatHTTPMessage(e.Kind, httpErr.StatusCode, httpErr.Body)
	}
	if e.Err == nil {
		return formatHTTPMessage(e.Kind, e.StatusCode, "")
	}
	return e.Err.Error()
}

func formatHTTPMessage(kind ErrorKind, statusCode int, body string) string {
	message := kindMessages[kind]
	if statusCode > 0 {
		message = fmt.Sprintf("%s (%d)", message, statusCode)
	}
	if body == "" {
		return message
	}
	return fmt.Sprintf("%s: %s", message, body)
}

var kindMessages = map[ErrorKind]string{
	ErrorUnauthorized: "unauthorized",
	ErrorNotFound:     "not found",
	ErrorConflict:     "conflict",
	ErrorValidation:   "validation failure",
	ErrorTransport:    "transport failure",
}

var statusToKind = map[int]ErrorKind{
	http.StatusUnauthorized:        ErrorUnauthorized,
	http.StatusForbidden:           ErrorUnauthorized,
	http.StatusNotFound:            ErrorNotFound,
	http.StatusConflict:            ErrorConflict,
	http.StatusBadRequest:          ErrorValidation,
	http.StatusUnprocessableEntity: ErrorValidation,
}

func kindFromStatus(statusCode int) ErrorKind {
	if kind, ok := statusToKind[statusCode]; ok {
		return kind
	}
	return ErrorTransport
}
