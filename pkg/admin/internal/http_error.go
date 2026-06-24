package internal

import "fmt"

type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("HTTP %d", e.StatusCode)
	}
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Body)
}

func newHTTPError(statusCode int, body []byte) error {
	return &HTTPError{StatusCode: statusCode, Body: string(body)}
}
