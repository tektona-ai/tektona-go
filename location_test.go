package tektona

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
)

func TestLocationList(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/locations" || r.URL.RawQuery != "" {
			t.Errorf("request = %s %s", r.Method, r.URL)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"items":[]}`)
	}, WithOrg("acme"), WithProject("web"))
	if client.Location() != client.Location() {
		t.Error("location accessor returned a different interface")
	}
	got, err := client.Location().List(context.Background())
	if err != nil || got == nil || got.Items == nil {
		t.Fatalf("list = %#v, %v", got, err)
	}
}

func TestLocationListResponses(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, contentType, body string
		status                  int
	}{
		{"error", "application/problem+json", `{"detail":"denied"}`, http.StatusForbidden},
		{"unexpected success", "application/json", `{}`, http.StatusCreated},
		{"missing JSON", "text/plain", `ok`, http.StatusOK},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", tc.contentType)
				w.WriteHeader(tc.status)
				_, _ = io.WriteString(w, tc.body)
			})
			_, err := client.Location().List(context.Background())
			var apiErr *APIError
			if tc.status != http.StatusOK && (!errors.As(err, &apiErr) || apiErr.StatusCode != tc.status) {
				t.Errorf("error = %v", err)
			}
			if tc.status == http.StatusOK && (err == nil || errors.As(err, &apiErr)) {
				t.Errorf("missing JSON error = %v", err)
			}
		})
	}
}
