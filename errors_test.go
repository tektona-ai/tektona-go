package tektona

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
)

func TestResourceMalformedProblemPreservesHTTPError(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		call func(*Client) error
	}{
		{"location", func(c *Client) error { _, err := c.Location().List(context.Background()); return err }},
		{"meta", func(c *Client) error { _, err := c.Meta().Get(context.Background()); return err }},
		{"organization", func(c *Client) error { _, err := c.Organization().Get(context.Background(), "org"); return err }},
		{"project", func(c *Client) error {
			_, err := c.Project().Get(context.Background(), "web", &GetProjectParams{Org: "org"})
			return err
		}},
		{"sandbox", func(c *Client) error {
			_, err := c.Sandbox().List(context.Background(), &ListSandboxesParams{Org: "org"})
			return err
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			body := []byte(`{"detail":`)
			client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/problem+json")
				w.Header().Set("Retry-After", "7")
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write(body)
			})
			err := tc.call(client)
			var apiErr *APIError
			if !errors.As(err, &apiErr) || apiErr.StatusCode != 429 || apiErr.Header.Get("Retry-After") != "7" || !bytes.Equal(apiErr.Body, body) || apiErr.Problem != nil {
				t.Errorf("error = %#v", err)
			}
		})
	}
}

func TestResourceErrorBodyVariants(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, contentType, body string
		wantProblem             bool
	}{
		{"plain text", "text/plain", "slow down", false},
		{"valid problem", "application/problem+json", `{"detail":"slow down"}`, true},
		{"JSON error", "application/json", `{"detail":"slow down"}`, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", tc.contentType)
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = io.WriteString(w, tc.body)
			})
			_, err := client.Location().List(context.Background())
			var apiErr *APIError
			if !errors.As(err, &apiErr) || apiErr.StatusCode != 429 || string(apiErr.Body) != tc.body || (apiErr.Problem != nil) != tc.wantProblem {
				t.Errorf("error = %#v", err)
			}
		})
	}
}

func TestResourceSuccessAndFailureModes(t *testing.T) {
	t.Parallel()
	t.Run("expected 204", func(t *testing.T) {
		t.Parallel()
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Errorf("method = %s", r.Method)
			}
			w.WriteHeader(http.StatusNoContent)
		}, WithOrg("org"))
		if err := client.Project().Delete(context.Background(), "web", nil); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("success decode failure", func(t *testing.T) {
		t.Parallel()
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"items":`)
		})
		_, err := client.Location().List(context.Background())
		var apiErr *APIError
		if err == nil || errors.As(err, &apiErr) {
			t.Errorf("decode error = %v", err)
		}
	})
	t.Run("transport failure", func(t *testing.T) {
		t.Parallel()
		broken := errors.New("transport failed")
		client, err := NewClient(WithAPIKey("key"), WithBaseURL("https://example.com"), WithHTTPClient(&http.Client{
			Transport: roundTripFunc(func(*http.Request) (*http.Response, error) { return nil, broken }),
		}))
		if err != nil {
			t.Fatal(err)
		}
		_, err = client.Meta().Get(context.Background())
		if !errors.Is(err, broken) {
			t.Errorf("transport error = %v", err)
		}
	})
}

func TestRawMalformedProblemKeepsGeneratedSemantics(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = io.WriteString(w, `{"detail":`)
	})
	_, err := client.Raw().ListLocationsWithResponse(context.Background())
	var apiErr *APIError
	if err == nil || errors.As(err, &apiErr) {
		t.Errorf("raw parse error = %v", err)
	}
}
