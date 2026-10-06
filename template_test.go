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

func TestTemplateMethods(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	q, cursor, display, description := "go", "next", "Go", "Go tools"
	limit := int32(3)
	for _, tc := range []struct {
		name, method, path, query, body, response string
		status                                    int
		call                                      func(Template) (any, error)
	}{
		{"list org", "GET", "/v1/orgs/acme/templates", "q=go&cursor=next&limit=3", "", `{"items":[]}`, 200, func(c Template) (any, error) {
			return c.ListForOrg(ctx, "", &ListTemplatesForOrgParams{Q: &q, Cursor: &cursor, Limit: &limit})
		}},
		{"create org", "POST", "/v1/orgs/other/templates", "", `{"metadata":{"name":"go","display_name":"Go","description":"Go tools"}}`, `{"id":"org-template"}`, 201, func(c Template) (any, error) {
			return c.CreateForOrg(ctx, CreateTemplateForOrgParams{Org: "other", Name: "go", DisplayName: &display, Description: &description})
		}},
		{"get org", "GET", "/v1/orgs/other/templates/go", "", "", `{"id":"org-template"}`, 200, func(c Template) (any, error) {
			return c.GetForOrg(ctx, "go", &GetTemplateForOrgParams{Org: "other"})
		}},
		{"update org", "PUT", "/v1/orgs/acme/templates/go", "", `{"metadata":{"name":"go","display_name":"Go","description":"Go tools"}}`, `{}`, 200, func(c Template) (any, error) {
			return c.UpdateForOrg(ctx, "go", UpdateTemplateForOrgParams{DisplayName: &display, Description: &description})
		}},
		{"delete org", "DELETE", "/v1/orgs/acme/templates/go", "", "", `{"builds":12}`, 200, func(c Template) (any, error) {
			return c.DeleteForOrg(ctx, "go", nil)
		}},
		{"activate org", "PUT", "/v1/orgs/acme/templates/go/activate", "", "", `{}`, 200, func(c Template) (any, error) {
			return c.ActivateForOrg(ctx, "go", nil)
		}},
		{"archive org", "PUT", "/v1/orgs/acme/templates/go/archive", "", "", `{}`, 200, func(c Template) (any, error) {
			return c.ArchiveForOrg(ctx, "go", nil)
		}},
		{"list project", "GET", "/v1/orgs/acme/projects/web/templates", "q=go&cursor=next&limit=3", "", `{"items":[]}`, 200, func(c Template) (any, error) {
			return c.ListForProject(ctx, "", "", &ListTemplatesForProjectParams{Q: &q, Cursor: &cursor, Limit: &limit})
		}},
		{"create project", "POST", "/v1/orgs/other/projects/another/templates", "", `{"metadata":{"name":"go","display_name":"Go","description":"Go tools"}}`, `{"id":"project-template"}`, 201, func(c Template) (any, error) {
			return c.CreateForProject(ctx, CreateTemplateForProjectParams{Org: "other", Project: "another", Name: "go", DisplayName: &display, Description: &description})
		}},
		{"get project", "GET", "/v1/orgs/other/projects/another/templates/go", "", "", `{"id":"project-template"}`, 200, func(c Template) (any, error) {
			return c.GetForProject(ctx, "go", &GetTemplateForProjectParams{Org: "other", Project: "another"})
		}},
		{"update project", "PUT", "/v1/orgs/acme/projects/web/templates/go", "", `{"metadata":{"name":"go","display_name":"Go","description":"Go tools"}}`, `{}`, 200, func(c Template) (any, error) {
			return c.UpdateForProject(ctx, "go", UpdateTemplateForProjectParams{DisplayName: &display, Description: &description})
		}},
		{"delete project", "DELETE", "/v1/orgs/acme/projects/web/templates/go", "", "", `{"builds":12}`, 200, func(c Template) (any, error) {
			return c.DeleteForProject(ctx, "go", nil)
		}},
		{"activate project", "PUT", "/v1/orgs/acme/projects/web/templates/go/activate", "", "", `{}`, 200, func(c Template) (any, error) {
			return c.ActivateForProject(ctx, "go", nil)
		}},
		{"archive project", "PUT", "/v1/orgs/acme/projects/web/templates/go/archive", "", "", `{}`, 200, func(c Template) (any, error) {
			return c.ArchiveForProject(ctx, "go", nil)
		}},
		{"get system", "GET", "/v1/system-templates/go", "", "", `{"id":"system-template"}`, 200, func(c Template) (any, error) {
			return c.GetSystem(ctx, "go")
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
					if err := json.Unmarshal([]byte(tc.body), &want); err != nil {
						t.Error(err)
					}
					if !reflect.DeepEqual(got, want) {
						t.Errorf("body = %v; want %v", got, want)
					}
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				_, _ = io.WriteString(w, tc.response)
			}, WithOrg("acme"), WithProject("web"))
			got, err := tc.call(client.Template())
			if err != nil || got == nil {
				t.Fatalf("result = %v, %v", got, err)
			}
			if tc.name == "delete org" || tc.name == "delete project" {
				if got.(*api.DeleteTemplateResponse).Builds != 12 {
					t.Errorf("delete result = %+v", got)
				}
			}
		})
	}
}

func TestTemplateValidation(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	for _, tc := range []struct {
		name string
		opts []Option
		call func(Template) error
	}{
		{"list org scope", nil, func(c Template) error { _, err := c.ListForOrg(ctx, "", nil); return err }},
		{"create org scope", nil, func(c Template) error {
			_, err := c.CreateForOrg(ctx, CreateTemplateForOrgParams{Name: "go"})
			return err
		}},
		{"create org name", []Option{WithOrg("acme")}, func(c Template) error { _, err := c.CreateForOrg(ctx, CreateTemplateForOrgParams{}); return err }},
		{"get org scope", nil, func(c Template) error { _, err := c.GetForOrg(ctx, "go", nil); return err }},
		{"get org name", []Option{WithOrg("acme")}, func(c Template) error { _, err := c.GetForOrg(ctx, "", nil); return err }},
		{"update org name", []Option{WithOrg("acme")}, func(c Template) error { _, err := c.UpdateForOrg(ctx, "", UpdateTemplateForOrgParams{}); return err }},
		{"delete org name", []Option{WithOrg("acme")}, func(c Template) error { _, err := c.DeleteForOrg(ctx, "", nil); return err }},
		{"activate org name", []Option{WithOrg("acme")}, func(c Template) error { _, err := c.ActivateForOrg(ctx, "", nil); return err }},
		{"archive org name", []Option{WithOrg("acme")}, func(c Template) error { _, err := c.ArchiveForOrg(ctx, "", nil); return err }},
		{"list project scope", []Option{WithOrg("acme")}, func(c Template) error { _, err := c.ListForProject(ctx, "", "", nil); return err }},
		{"create project scope", []Option{WithOrg("acme")}, func(c Template) error {
			_, err := c.CreateForProject(ctx, CreateTemplateForProjectParams{Name: "go"})
			return err
		}},
		{"create project name", []Option{WithOrg("acme"), WithProject("web")}, func(c Template) error {
			_, err := c.CreateForProject(ctx, CreateTemplateForProjectParams{})
			return err
		}},
		{"get project name", []Option{WithOrg("acme"), WithProject("web")}, func(c Template) error { _, err := c.GetForProject(ctx, "", nil); return err }},
		{"update project name", []Option{WithOrg("acme"), WithProject("web")}, func(c Template) error {
			_, err := c.UpdateForProject(ctx, "", UpdateTemplateForProjectParams{})
			return err
		}},
		{"delete project name", []Option{WithOrg("acme"), WithProject("web")}, func(c Template) error { _, err := c.DeleteForProject(ctx, "", nil); return err }},
		{"activate project name", []Option{WithOrg("acme"), WithProject("web")}, func(c Template) error { _, err := c.ActivateForProject(ctx, "", nil); return err }},
		{"archive project name", []Option{WithOrg("acme"), WithProject("web")}, func(c Template) error { _, err := c.ArchiveForProject(ctx, "", nil); return err }},
		{"get system name", nil, func(c Template) error { _, err := c.GetSystem(ctx, ""); return err }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { t.Errorf("unexpected request: %s %s", r.Method, r.URL) }, tc.opts...)
			if err := tc.call(client.Template()); err == nil {
				t.Error("expected local validation error")
			}
		})
	}
}

func TestTemplateResponseErrors(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		call func(Template) error
	}{
		{"list", func(c Template) error { _, err := c.ListForOrg(context.Background(), "", nil); return err }},
		{"create", func(c Template) error {
			_, err := c.CreateForProject(context.Background(), CreateTemplateForProjectParams{Name: "go"})
			return err
		}},
		{"delete", func(c Template) error { _, err := c.DeleteForOrg(context.Background(), "go", nil); return err }},
		{"system", func(c Template) error { _, err := c.GetSystem(context.Background(), "go"); return err }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/problem+json")
				w.Header().Set("Retry-After", "5")
				w.WriteHeader(429)
				_, _ = io.WriteString(w, `{"detail":"later"}`)
			}, WithOrg("acme"), WithProject("web"))
			var apiErr *APIError
			if err := tc.call(client.Template()); !errors.As(err, &apiErr) || apiErr.StatusCode != 429 || apiErr.Header.Get("Retry-After") != "5" || apiErr.Problem == nil {
				t.Errorf("error = %v", err)
			}
		})
	}
}

func TestTemplateScopePrecedence(t *testing.T) {
	t.Setenv("TEKTONA_ORG", "env-org")
	t.Setenv("TEKTONA_PROJECT", "env-project")
	paths := []string{
		"/v1/orgs/option-org/templates/go",
		"/v1/orgs/option-org/projects/env-project/templates/go",
		"/v1/orgs/call-org/projects/call-project/templates/go",
		"/v1/system-templates/go",
	}
	index := 0
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if index >= len(paths) || r.URL.Path != paths[index] {
			t.Errorf("request %d path = %s", index, r.URL.Path)
		}
		index++
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{}`)
	}, WithOrg("option-org"))
	ctx := context.Background()
	if _, err := client.Template().GetForOrg(ctx, "go", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Template().GetForProject(ctx, "go", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Template().GetForProject(ctx, "go", &GetTemplateForProjectParams{Org: "call-org", Project: "call-project"}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Template().GetSystem(ctx, "go"); err != nil {
		t.Fatal(err)
	}
	if index != len(paths) {
		t.Errorf("requests = %d; want %d", index, len(paths))
	}

	blocked := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	}, WithOrg(""), WithProject(""))
	if _, err := blocked.Template().GetForOrg(ctx, "go", nil); err == nil {
		t.Error("explicit empty organization inherited environment")
	}
	if _, err := blocked.Template().GetForProject(ctx, "go", nil); err == nil {
		t.Error("explicit empty project inherited environment")
	}
}

func TestTemplateSuccessStatusAndJSON(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, contentType string
		status            int
	}{
		{"delete requires 200", "application/json", 204},
		{"delete requires JSON", "text/plain", 200},
		{"create requires 201", "application/json", 200},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", tc.contentType)
				w.WriteHeader(tc.status)
				_, _ = io.WriteString(w, `{}`)
			}, WithOrg("acme"))
			var err error
			if tc.name == "create requires 201" {
				_, err = client.Template().CreateForOrg(context.Background(), CreateTemplateForOrgParams{Name: "go"})
			} else {
				_, err = client.Template().DeleteForOrg(context.Background(), "go", nil)
			}
			if err == nil {
				t.Error("expected response error")
			}
		})
	}
}
