package tektona

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/tektona-ai/tektona-go/api"
)

func TestProjectMethods(t *testing.T) {
	t.Parallel()
	cursor, limit, description := "next", int32(2), "A project"
	pause, deleteAfter, resume := "15m", "7d", true
	mode := api.LifecycleDefaultsBodyAutoPauseMode("suspend")
	for _, tc := range []struct {
		name, method, path, query, body, response string
		status                                    int
		call                                      func(Project) (any, error)
	}{
		{"list all", "GET", "/v1/projects", "cursor=next&limit=2", "", `{"items":[]}`, 200, func(p Project) (any, error) {
			return p.List(context.Background(), &ListProjectsParams{Cursor: &cursor, Limit: &limit})
		}},
		{"list for org", "GET", "/v1/orgs/acme/projects", "cursor=next&limit=2", "", `{"items":[]}`, 200, func(p Project) (any, error) {
			return p.ListForOrg(context.Background(), "", &ListProjectsForOrgParams{Cursor: &cursor, Limit: &limit})
		}},
		{"create", "POST", "/v1/orgs/other/projects", "", `{"name":"new","display_name":"New","description":"A project"}`, `{"name":"new"}`, 201, func(p Project) (any, error) {
			return p.Create(context.Background(), CreateProjectParams{Org: "other", Name: "new", DisplayName: "New", Description: &description})
		}},
		{"get default", "GET", "/v1/orgs/acme/projects/web", "", "", `{"name":"web"}`, 200, func(p Project) (any, error) {
			return p.Get(context.Background(), "web", nil)
		}},
		{"get override", "GET", "/v1/orgs/other/projects/another", "", "", `{"name":"another"}`, 200, func(p Project) (any, error) {
			return p.Get(context.Background(), "another", &GetProjectParams{Org: "other"})
		}},
		{"update", "PUT", "/v1/orgs/other/projects/another", "", `{"display_name":"New","description":"A project"}`, `{"name":"another"}`, 200, func(p Project) (any, error) {
			return p.Update(context.Background(), "another", UpdateProjectParams{Org: "other", DisplayName: "New", Description: &description})
		}},
		{"delete", "DELETE", "/v1/orgs/other/projects/another", "", "", "", 204, func(p Project) (any, error) {
			err := p.Delete(context.Background(), "another", &DeleteProjectParams{Org: "other"})
			return true, err
		}},
		{"get lifecycle", "GET", "/v1/orgs/other/projects/another/lifecycle-defaults", "", "", `{}`, 200, func(p Project) (any, error) {
			return p.GetLifecycleDefaults(context.Background(), "another", &GetProjectLifecycleDefaultsParams{Org: "other"})
		}},
		{"update lifecycle", "PUT", "/v1/orgs/acme/projects/web/lifecycle-defaults", "", `{"auto_delete_after":"7d","auto_pause_after":"15m","auto_pause_mode":"suspend","auto_resume":true}`, `{}`, 200, func(p Project) (any, error) {
			return p.UpdateLifecycleDefaults(context.Background(), "web", UpdateProjectLifecycleDefaultsParams{
				AutoDeleteAfter: &deleteAfter, AutoPauseAfter: &pause, AutoPauseMode: &mode, AutoResume: &resume,
			})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tc.method || r.URL.Path != tc.path || r.URL.RawQuery != tc.query {
					t.Errorf("request = %s %s; want %s %s?%s", r.Method, r.URL, tc.method, tc.path, tc.query)
				}
				if tc.body != "" {
					var got, want any
					if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
						t.Error(err)
					}
					_ = json.Unmarshal([]byte(tc.body), &want)
					if !reflect.DeepEqual(got, want) {
						t.Errorf("body = %v; want %v", got, want)
					}
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				_, _ = io.WriteString(w, tc.response)
			}, WithOrg("acme"), WithProject("web"))
			if client.Project() != client.Project() {
				t.Error("accessor returned a different interface")
			}
			got, err := tc.call(client.Project())
			if err != nil || got == nil {
				t.Fatalf("result = %v, %v", got, err)
			}
		})
	}
}

func TestProjectValidation(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	})
	for _, call := range []func() error{
		func() error { _, err := client.Project().ListForOrg(context.Background(), "", nil); return err },
		func() error {
			_, err := client.Project().Create(context.Background(), CreateProjectParams{})
			return err
		},
		func() error { _, err := client.Project().Get(context.Background(), "", nil); return err },
		func() error {
			_, err := client.Project().Update(context.Background(), "", UpdateProjectParams{})
			return err
		},
		func() error { return client.Project().Delete(context.Background(), "", nil) },
		func() error {
			_, err := client.Project().GetLifecycleDefaults(context.Background(), "", nil)
			return err
		},
		func() error {
			_, err := client.Project().UpdateLifecycleDefaults(context.Background(), "", UpdateProjectLifecycleDefaultsParams{})
			return err
		},
	} {
		if err := call(); err == nil {
			t.Error("expected local validation error")
		}
	}
}

func TestProjectResponses(t *testing.T) {
	t.Parallel()
	for _, call := range []struct {
		name string
		fn   func(Project) error
	}{
		{"list", func(p Project) error { _, err := p.List(context.Background(), nil); return err }},
		{"list for org", func(p Project) error { _, err := p.ListForOrg(context.Background(), "", nil); return err }},
		{"create", func(p Project) error {
			_, err := p.Create(context.Background(), CreateProjectParams{Name: "web", DisplayName: "Web"})
			return err
		}},
		{"get", func(p Project) error { _, err := p.Get(context.Background(), "web", nil); return err }},
		{"update", func(p Project) error {
			_, err := p.Update(context.Background(), "web", UpdateProjectParams{DisplayName: "Web"})
			return err
		}},
		{"delete", func(p Project) error { return p.Delete(context.Background(), "web", nil) }},
		{"get lifecycle", func(p Project) error { _, err := p.GetLifecycleDefaults(context.Background(), "web", nil); return err }},
		{"update lifecycle", func(p Project) error {
			_, err := p.UpdateLifecycleDefaults(context.Background(), "web", UpdateProjectLifecycleDefaultsParams{})
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
			}, WithOrg("acme"), WithProject("web"))
			err := call.fn(client.Project())
			var apiErr *APIError
			if !errors.As(err, &apiErr) || apiErr.StatusCode != 429 || apiErr.Header.Get("Retry-After") != "5" || apiErr.Problem == nil {
				t.Errorf("error = %v", err)
			}
		})
	}
}

func TestProjectMissingJSONAndDeleteStatus(t *testing.T) {
	t.Parallel()
	for _, status := range []int{http.StatusOK, http.StatusCreated} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(status)
				_, _ = io.WriteString(w, "ok")
			}, WithOrg("acme"), WithProject("web"))
			var err error
			if status == http.StatusOK {
				_, err = client.Project().Get(context.Background(), "web", nil)
			} else {
				err = client.Project().Delete(context.Background(), "web", nil)
			}
			if err == nil {
				t.Error("expected status or JSON error")
			}
		})
	}
}

func TestProjectScopePrecedence(t *testing.T) {
	t.Setenv("TEKTONA_ORG", "env-org")
	t.Setenv("TEKTONA_PROJECT", "env-project")
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/orgs/option-org/projects/web" {
			t.Errorf("path = %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{}`)
	}, WithOrg("option-org"))
	if client.org != "option-org" || client.project != "env-project" {
		t.Errorf("stored scope = %s/%s", client.org, client.project)
	}
	if _, err := client.Project().Get(context.Background(), "web", nil); err != nil {
		t.Fatal(err)
	}

	blocked := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	}, WithOrg(""), WithProject(""))
	if _, err := blocked.Project().Get(context.Background(), "web", nil); err == nil {
		t.Error("explicit empty options inherited environment")
	}
}

func TestProjectRejectsEmptyResourceName(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	}, WithOrg("production"), WithProject("production"))
	for _, tc := range []struct {
		name string
		call func() error
	}{
		{"get", func() error { _, err := client.Project().Get(context.Background(), "", nil); return err }},
		{"update", func() error {
			_, err := client.Project().Update(context.Background(), "", UpdateProjectParams{DisplayName: "Production"})
			return err
		}},
		{"delete", func() error { return client.Project().Delete(context.Background(), "", nil) }},
		{"get lifecycle", func() error {
			_, err := client.Project().GetLifecycleDefaults(context.Background(), "", nil)
			return err
		}},
		{"update lifecycle", func() error {
			_, err := client.Project().UpdateLifecycleDefaults(context.Background(), "", UpdateProjectLifecycleDefaultsParams{})
			return err
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.call(); err == nil {
				t.Error("expected a local error for an empty project name")
			}
		})
	}
}
