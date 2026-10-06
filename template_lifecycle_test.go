package tektona

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestTemplateLifecycleMethods(t *testing.T) {
	t.Parallel()
	archive, deleteAfter := int32(30), int32(0)
	settings := `{"archive_after_days_unused":30,"delete_after_days_archived":0}`
	for _, tc := range []struct {
		name, method, path, body, response string
		call                               func(TemplateLifecycle) (any, error)
	}{
		{"get org default", "GET", "/v1/orgs/acme/settings/template-lifecycle", "", `{"settings":{}}`, func(l TemplateLifecycle) (any, error) {
			return l.GetForOrg(context.Background(), "")
		}},
		{"get org override", "GET", "/v1/orgs/other/settings/template-lifecycle", "", `{"settings":{}}`, func(l TemplateLifecycle) (any, error) {
			return l.GetForOrg(context.Background(), "other")
		}},
		{"update org", "PUT", "/v1/orgs/acme/settings/template-lifecycle", settings, `{"settings":{}}`, func(l TemplateLifecycle) (any, error) {
			return l.UpdateForOrg(context.Background(), "", UpdateTemplateLifecycleForOrgParams{ArchiveAfterDaysUnused: &archive, DeleteAfterDaysArchived: &deleteAfter})
		}},
		{"preview org", "PUT", "/v1/orgs/other/settings/template-lifecycle/preview", settings, `{"versions":[]}`, func(l TemplateLifecycle) (any, error) {
			return l.PreviewForOrg(context.Background(), "other", PreviewTemplateLifecycleForOrgParams{ArchiveAfterDaysUnused: &archive, DeleteAfterDaysArchived: &deleteAfter})
		}},
		{"get project default", "GET", "/v1/orgs/acme/projects/web/settings/template-lifecycle", "", `{"settings":{}}`, func(l TemplateLifecycle) (any, error) {
			return l.GetForProject(context.Background(), "web", nil)
		}},
		{"get project override", "GET", "/v1/orgs/other/projects/web/settings/template-lifecycle", "", `{"settings":{}}`, func(l TemplateLifecycle) (any, error) {
			return l.GetForProject(context.Background(), "web", &GetTemplateLifecycleForProjectParams{Org: "other"})
		}},
		{"update project", "PUT", "/v1/orgs/other/projects/web/settings/template-lifecycle", settings, `{"settings":{}}`, func(l TemplateLifecycle) (any, error) {
			return l.UpdateForProject(context.Background(), "web", UpdateTemplateLifecycleForProjectParams{Org: "other", ArchiveAfterDaysUnused: &archive, DeleteAfterDaysArchived: &deleteAfter})
		}},
		{"preview project", "PUT", "/v1/orgs/acme/projects/web/settings/template-lifecycle/preview", settings, `{"versions":[]}`, func(l TemplateLifecycle) (any, error) {
			return l.PreviewForProject(context.Background(), "web", PreviewTemplateLifecycleForProjectParams{ArchiveAfterDaysUnused: &archive, DeleteAfterDaysArchived: &deleteAfter})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tc.method || r.URL.Path != tc.path || r.URL.RawQuery != "" {
					t.Errorf("request = %s %s; want %s %s", r.Method, r.URL, tc.method, tc.path)
				}
				if tc.body != "" {
					var got, want any
					if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
						t.Error(err)
					}
					if err := json.Unmarshal([]byte(tc.body), &want); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, want) {
						t.Errorf("body = %v; want %v", got, want)
					}
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = io.WriteString(w, tc.response)
			}, WithOrg("acme"), WithProject("ignored"))
			got, err := tc.call(client.TemplateLifecycle())
			if err != nil || got == nil {
				t.Fatalf("result = %v, %v", got, err)
			}
		})
	}
}

func TestTemplateLifecycleValidation(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		opts []Option
		call func(TemplateLifecycle) error
	}{
		{"get org", nil, func(l TemplateLifecycle) error { _, err := l.GetForOrg(context.Background(), ""); return err }},
		{"update org", nil, func(l TemplateLifecycle) error {
			_, err := l.UpdateForOrg(context.Background(), "", UpdateTemplateLifecycleForOrgParams{})
			return err
		}},
		{"preview org", nil, func(l TemplateLifecycle) error {
			_, err := l.PreviewForOrg(context.Background(), "", PreviewTemplateLifecycleForOrgParams{})
			return err
		}},
		{"get project missing org", nil, func(l TemplateLifecycle) error {
			_, err := l.GetForProject(context.Background(), "web", nil)
			return err
		}},
		{"get project missing project", []Option{WithOrg("acme")}, func(l TemplateLifecycle) error { _, err := l.GetForProject(context.Background(), "", nil); return err }},
		{"update project missing org", nil, func(l TemplateLifecycle) error {
			_, err := l.UpdateForProject(context.Background(), "web", UpdateTemplateLifecycleForProjectParams{})
			return err
		}},
		{"update project missing project", []Option{WithOrg("acme")}, func(l TemplateLifecycle) error {
			_, err := l.UpdateForProject(context.Background(), "", UpdateTemplateLifecycleForProjectParams{})
			return err
		}},
		{"preview project missing org", nil, func(l TemplateLifecycle) error {
			_, err := l.PreviewForProject(context.Background(), "web", PreviewTemplateLifecycleForProjectParams{})
			return err
		}},
		{"preview project missing project", []Option{WithOrg("acme")}, func(l TemplateLifecycle) error {
			_, err := l.PreviewForProject(context.Background(), "", PreviewTemplateLifecycleForProjectParams{})
			return err
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				t.Errorf("unexpected request: %s %s", r.Method, r.URL)
			}, append([]Option{WithOrg("")}, tc.opts...)...)
			if err := tc.call(client.TemplateLifecycle()); err == nil {
				t.Error("expected local validation error")
			}
		})
	}
}

func TestTemplateLifecycleResponses(t *testing.T) {
	t.Parallel()
	for _, call := range []struct {
		name string
		fn   func(TemplateLifecycle) error
	}{
		{"get org", func(l TemplateLifecycle) error { _, err := l.GetForOrg(context.Background(), ""); return err }},
		{"update org", func(l TemplateLifecycle) error {
			_, err := l.UpdateForOrg(context.Background(), "", UpdateTemplateLifecycleForOrgParams{})
			return err
		}},
		{"preview org", func(l TemplateLifecycle) error {
			_, err := l.PreviewForOrg(context.Background(), "", PreviewTemplateLifecycleForOrgParams{})
			return err
		}},
		{"get project", func(l TemplateLifecycle) error {
			_, err := l.GetForProject(context.Background(), "web", nil)
			return err
		}},
		{"update project", func(l TemplateLifecycle) error {
			_, err := l.UpdateForProject(context.Background(), "web", UpdateTemplateLifecycleForProjectParams{})
			return err
		}},
		{"preview project", func(l TemplateLifecycle) error {
			_, err := l.PreviewForProject(context.Background(), "web", PreviewTemplateLifecycleForProjectParams{})
			return err
		}},
	} {
		t.Run(call.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/problem+json")
				w.Header().Set("Retry-After", "5")
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = io.WriteString(w, `{"detail":"later"}`)
			}, WithOrg("acme"))
			err := call.fn(client.TemplateLifecycle())
			var apiErr *APIError
			if !errors.As(err, &apiErr) || apiErr.StatusCode != 429 || apiErr.Header.Get("Retry-After") != "5" || apiErr.Problem == nil {
				t.Errorf("error = %v", err)
			}
		})
	}
}

func TestTemplateLifecycleUnexpectedResponses(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, contentType, body string
		status                  int
	}{
		{"wrong status", "application/json", `{}`, http.StatusCreated},
		{"missing JSON", "text/plain", "ok", http.StatusOK},
		{"invalid JSON", "application/json", "{", http.StatusOK},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", tc.contentType)
				w.WriteHeader(tc.status)
				_, _ = io.WriteString(w, tc.body)
			}, WithOrg("acme"))
			_, err := client.TemplateLifecycle().GetForOrg(context.Background(), "")
			var apiErr *APIError
			if tc.status != http.StatusOK && (!errors.As(err, &apiErr) || apiErr.StatusCode != tc.status) {
				t.Errorf("error = %v", err)
			}
			if tc.status == http.StatusOK && (err == nil || errors.As(err, &apiErr)) {
				t.Errorf("decode error = %v", err)
			}
		})
	}
}

func TestTemplateLifecycleScopePrecedence(t *testing.T) {
	t.Setenv("TEKTONA_ORG", "env-org")
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/orgs/option-org/projects/web/settings/template-lifecycle" {
			t.Errorf("path = %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{}`)
	}, WithOrg("option-org"))
	if _, err := client.TemplateLifecycle().GetForProject(context.Background(), "web", nil); err != nil {
		t.Fatal(err)
	}

	blocked := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	}, WithOrg(""))
	if _, err := blocked.TemplateLifecycle().GetForOrg(context.Background(), ""); err == nil {
		t.Error("explicit empty option inherited environment")
	}
}
