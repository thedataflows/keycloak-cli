package kcapi

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorKind classifies a Keycloak API failure. The kinds are the unified
// error taxonomy for the kcapi library: everything the runtime client and
// the generic client surface is a *Error carrying one of these kinds.
type ErrorKind string

const (
	KindNotFound   ErrorKind = "not_found"
	KindConflict   ErrorKind = "conflict"
	KindValidation ErrorKind = "validation"
	KindAuth       ErrorKind = "auth"
	KindServer     ErrorKind = "server"
	KindNetwork    ErrorKind = "network"
)

// Sentinel errors matched via errors.Is through (*Error).Is: a *Error whose
// Kind maps to the sentinel satisfies it without inspecting the chain.
var (
	ErrNotFound   = errors.New("kcapi: not found")
	ErrConflict   = errors.New("kcapi: conflict")
	ErrValidation = errors.New("kcapi: validation failed")
	ErrAuth       = errors.New("kcapi: authentication failed")
)

// Error is the typed error returned for failed Keycloak API calls. Status and
// Body carry the HTTP response when the failure came from one; Err wraps the
// underlying cause.
type Error struct {
	Kind   ErrorKind
	Op     string // e.g. "invoke getUser", "fetch users"
	Status int    // HTTP status, 0 for non-HTTP
	Body   string // truncated response body
	Err    error  // wrapped cause
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	message := e.message()
	if e.Op == "" {
		return message
	}
	if message == "" {
		return e.Op
	}
	return fmt.Sprintf("%s: %s", e.Op, message)
}

func (e *Error) Unwrap() error {
	return e.Err
}

// Is maps Kind to the sentinel errors so callers can use errors.Is instead of
// matching on Kind directly; any other target falls through to the cause chain.
func (e *Error) Is(target error) bool {
	switch target {
	case ErrNotFound:
		return e.Kind == KindNotFound
	case ErrConflict:
		return e.Kind == KindConflict
	case ErrValidation:
		return e.Kind == KindValidation
	case ErrAuth:
		return e.Kind == KindAuth
	}
	return errors.Is(e.Err, target)
}

// message renders the error text: a classified HTTP failure renders as
// "kind (status): body"; a failure with no body falls back to the wrapped
// cause's own text so transport errors keep their detail.
func (e *Error) message() string {
	if e.Body != "" {
		return formatHTTPMessage(e.Kind, e.Status, e.Body)
	}
	if e.Err == nil {
		return formatHTTPMessage(e.Kind, e.Status, "")
	}
	var httpErr *Error
	if errors.As(e.Err, &httpErr) {
		return formatHTTPMessage(e.Kind, httpErr.Status, httpErr.Body)
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
	KindAuth:       "unauthorized",
	KindNotFound:   "not found",
	KindConflict:   "conflict",
	KindValidation: "validation failure",
	KindServer:     "server error",
	KindNetwork:    "transport failure",
}

var statusToKind = map[int]ErrorKind{
	http.StatusUnauthorized:        KindAuth,
	http.StatusForbidden:           KindAuth,
	http.StatusNotFound:            KindNotFound,
	http.StatusConflict:            KindConflict,
	http.StatusBadRequest:          KindValidation,
	http.StatusUnprocessableEntity: KindValidation,
}

// kindFromStatus maps an HTTP status to its ErrorKind. 5xx are server
// failures; anything unrecognized (including 0, i.e. no response at all) is a
// network/transport failure.
func kindFromStatus(statusCode int) ErrorKind {
	if statusCode >= http.StatusInternalServerError {
		return KindServer
	}
	if kind, ok := statusToKind[statusCode]; ok {
		return kind
	}
	return KindNetwork
}

// classifyError stamps the operation context onto a failed call. When the
// failure already carries an HTTP response (a *Error from newHTTPError), its
// status and body win and the kind is re-derived from that status.
func classifyError(err error, statusCode int, op string) *Error {
	if err == nil {
		return nil
	}
	var httpErr *Error
	if errors.As(err, &httpErr) {
		statusCode = httpErr.Status
		return &Error{
			Kind:   kindFromStatus(statusCode),
			Op:     op,
			Status: statusCode,
			Body:   httpErr.Body,
			Err:    err,
		}
	}
	return &Error{
		Kind:   kindFromStatus(statusCode),
		Op:     op,
		Status: statusCode,
		Err:    err,
	}
}

// newHTTPError builds the typed error for a non-2xx HTTP response. It replaces
// the old admin/internal HTTPError carrier: status and truncated body live on
// *Error so classification and rendering need no extra type in the chain.
func newHTTPError(statusCode int, body []byte) error {
	return &Error{
		Kind:   kindFromStatus(statusCode),
		Status: statusCode,
		Body:   string(body),
	}
}
