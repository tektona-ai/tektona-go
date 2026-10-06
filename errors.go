package tektona

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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

func resourceBody(call func() (*http.Response, error), expected int) ([]byte, string, error) {
	resp, err := call()
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	contentType := resp.Header.Get("Content-Type")
	if resp.StatusCode != expected {
		var problem *api.ErrorModel
		if strings.Contains(contentType, "json") {
			var parsed api.ErrorModel
			if json.Unmarshal(body, &parsed) == nil {
				problem = &parsed
			}
		}
		return nil, "", responseError(resp.StatusCode, resp.Header, body, problem)
	}
	return body, contentType, nil
}

func resourceJSON[T any](call func() (*http.Response, error), expected int, name string) (*T, error) {
	body, contentType, err := resourceBody(call, expected)
	if err != nil {
		return nil, err
	}
	if !strings.Contains(contentType, "json") {
		return nil, fmt.Errorf("%s: HTTP %d lacks a JSON response", name, expected)
	}
	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func resourceNoContent(call func() (*http.Response, error)) error {
	_, _, err := resourceBody(call, http.StatusNoContent)
	return err
}
