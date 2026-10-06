package tektona

import (
	"fmt"
	"net/http"

	"github.com/tektona-ai/tektona-go/api"
)

// APIError reports an HTTP error response. Use errors.As to read its status and problem.
type APIError struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Problem    *api.ErrorModel
}

func (e *APIError) Error() string {
	if e.Problem != nil && e.Problem.Detail != nil {
		return fmt.Sprintf("Tektona API %d: %s", e.StatusCode, *e.Problem.Detail)
	}
	return fmt.Sprintf("Tektona API: HTTP %d", e.StatusCode)
}

func responseError(status int, header http.Header, body []byte, problem *api.ErrorModel) error {
	return &APIError{StatusCode: status, Header: header.Clone(), Body: body, Problem: problem}
}
