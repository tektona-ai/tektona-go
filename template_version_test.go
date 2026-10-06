package tektona

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/tektona-ai/tektona-go/api"
)

func TestTemplateVersionRoutes(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	stateOrg := api.ListOrgTemplateVersionsParamsStateArchived
	stateProject := api.ListProjectTemplateVersionsParamsStateActive
	stateSystem := api.ListSystemTemplateVersionsParamsStateActive
	tagged, dryRun := false, true
	q, cursor, limit := "ubuntu", "next", int32(3)
	for _, tc := range []struct {
		name, method, path, query, body, response string
		status                                    int
		call                                      func(TemplateVersion) (any, error)
	}{
		{"org list", "GET", "/v1/orgs/acme/templates/base/versions", "cursor=next&limit=3&q=ubuntu&state=archived&tagged=false", "", `{"items":[]}`, 200, func(v TemplateVersion) (any, error) {
			return v.ListForOrg(ctx, "base", &ListTemplateVersionsForOrgParams{State: &stateOrg, Tagged: &tagged, Q: &q, Cursor: &cursor, Limit: &limit})
		}},
		{"org get", "GET", "/v1/orgs/other/templates/base/versions/v1", "", "", `{"id":"v1"}`, 200, func(v TemplateVersion) (any, error) {
			return v.GetForOrg(ctx, "base", "v1", &TemplateVersionForOrgParams{Org: "other"})
		}},
		{"org delete", "DELETE", "/v1/orgs/acme/templates/base/versions/v1", "", "", "", 204, func(v TemplateVersion) (any, error) {
			return true, v.DeleteForOrg(ctx, "base", "v1", nil)
		}},
		{"org prune", "PUT", "/v1/orgs/other/templates/base/versions/prune", "", `{"dry_run":true}`, `{"dry_run":true}`, 200, func(v TemplateVersion) (any, error) {
			return v.PruneForOrg(ctx, "base", PruneTemplateVersionsForOrgParams{Org: "other", DryRun: &dryRun})
		}},
		{"org activate", "PUT", "/v1/orgs/acme/templates/base/versions/v1/activate", "", "", `{"id":"v1"}`, 200, func(v TemplateVersion) (any, error) {
			return v.ActivateForOrg(ctx, "base", "v1", nil)
		}},
		{"org archive", "PUT", "/v1/orgs/acme/templates/base/versions/v1/archive", "", "", `{"id":"v1"}`, 200, func(v TemplateVersion) (any, error) {
			return v.ArchiveForOrg(ctx, "base", "v1", nil)
		}},
		{"project list", "GET", "/v1/orgs/other/projects/else/templates/base/versions", "cursor=next&limit=3&q=ubuntu&state=active&tagged=false", "", `{"items":[]}`, 200, func(v TemplateVersion) (any, error) {
			return v.ListForProject(ctx, "base", &ListTemplateVersionsForProjectParams{Org: "other", Project: "else", State: &stateProject, Tagged: &tagged, Q: &q, Cursor: &cursor, Limit: &limit})
		}},
		{"project get", "GET", "/v1/orgs/acme/projects/web/templates/base/versions/v1", "", "", `{"id":"v1"}`, 200, func(v TemplateVersion) (any, error) {
			return v.GetForProject(ctx, "base", "v1", nil)
		}},
		{"project delete", "DELETE", "/v1/orgs/other/projects/else/templates/base/versions/v1", "", "", "", 204, func(v TemplateVersion) (any, error) {
			return true, v.DeleteForProject(ctx, "base", "v1", &TemplateVersionForProjectParams{Org: "other", Project: "else"})
		}},
		{"project prune", "PUT", "/v1/orgs/acme/projects/else/templates/base/versions/prune", "", `{"dry_run":true}`, `{"dry_run":true}`, 200, func(v TemplateVersion) (any, error) {
			return v.PruneForProject(ctx, "base", PruneTemplateVersionsForProjectParams{Project: "else", DryRun: &dryRun})
		}},
		{"project activate", "PUT", "/v1/orgs/acme/projects/web/templates/base/versions/v1/activate", "", "", `{"id":"v1"}`, 200, func(v TemplateVersion) (any, error) {
			return v.ActivateForProject(ctx, "base", "v1", nil)
		}},
		{"project archive", "PUT", "/v1/orgs/acme/projects/web/templates/base/versions/v1/archive", "", "", `{"id":"v1"}`, 200, func(v TemplateVersion) (any, error) {
			return v.ArchiveForProject(ctx, "base", "v1", nil)
		}},
		{"system list", "GET", "/v1/system-templates/base/versions", "cursor=next&limit=3&q=ubuntu&state=active&tagged=false", "", `{"items":[]}`, 200, func(v TemplateVersion) (any, error) {
			return v.ListSystem(ctx, "base", &ListTemplateVersionsSystemParams{State: &stateSystem, Tagged: &tagged, Q: &q, Cursor: &cursor, Limit: &limit})
		}},
		{"system get", "GET", "/v1/system-templates/base/versions/v1", "", "", `{"id":"v1"}`, 200, func(v TemplateVersion) (any, error) {
			return v.GetSystem(ctx, "base", "v1")
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				wantQuery, err := url.ParseQuery(tc.query)
				if err != nil {
					t.Fatal(err)
				}
				if r.Method != tc.method || r.URL.Path != tc.path || !reflect.DeepEqual(r.URL.Query(), wantQuery) {
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
			got, err := tc.call(client.TemplateVersion())
			if err != nil || got == nil {
				t.Fatalf("result = %v, %v", got, err)
			}
		})
	}
}

func TestTemplateVersionValidation(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	for _, tc := range []struct {
		name string
		opts []Option
		call func(TemplateVersion) error
	}{
		{"org list no scope", nil, func(v TemplateVersion) error { _, err := v.ListForOrg(ctx, "base", nil); return err }},
		{"org list no name", []Option{WithOrg("acme")}, func(v TemplateVersion) error { _, err := v.ListForOrg(ctx, "", nil); return err }},
		{"org get no id", []Option{WithOrg("acme")}, func(v TemplateVersion) error { _, err := v.GetForOrg(ctx, "base", "", nil); return err }},
		{"org delete no id", []Option{WithOrg("acme")}, func(v TemplateVersion) error { return v.DeleteForOrg(ctx, "base", "", nil) }},
		{"org prune no name", []Option{WithOrg("acme")}, func(v TemplateVersion) error {
			_, err := v.PruneForOrg(ctx, "", PruneTemplateVersionsForOrgParams{})
			return err
		}},
		{"org activate no id", []Option{WithOrg("acme")}, func(v TemplateVersion) error { _, err := v.ActivateForOrg(ctx, "base", "", nil); return err }},
		{"org archive no scope", nil, func(v TemplateVersion) error { _, err := v.ArchiveForOrg(ctx, "base", "id", nil); return err }},
		{"project list no project", []Option{WithOrg("acme")}, func(v TemplateVersion) error { _, err := v.ListForProject(ctx, "base", nil); return err }},
		{"project get no name", []Option{WithOrg("acme"), WithProject("web")}, func(v TemplateVersion) error { _, err := v.GetForProject(ctx, "", "id", nil); return err }},
		{"project delete no id", []Option{WithOrg("acme"), WithProject("web")}, func(v TemplateVersion) error { return v.DeleteForProject(ctx, "base", "", nil) }},
		{"project prune no org", []Option{WithProject("web")}, func(v TemplateVersion) error {
			_, err := v.PruneForProject(ctx, "base", PruneTemplateVersionsForProjectParams{})
			return err
		}},
		{"project activate no id", []Option{WithOrg("acme"), WithProject("web")}, func(v TemplateVersion) error { _, err := v.ActivateForProject(ctx, "base", "", nil); return err }},
		{"project archive no scope", nil, func(v TemplateVersion) error { _, err := v.ArchiveForProject(ctx, "base", "id", nil); return err }},
		{"system list no name", nil, func(v TemplateVersion) error { _, err := v.ListSystem(ctx, "", nil); return err }},
		{"system get no id", nil, func(v TemplateVersion) error { _, err := v.GetSystem(ctx, "base", ""); return err }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				t.Errorf("unexpected request: %s %s", r.Method, r.URL)
			}, append([]Option{WithOrg(""), WithProject("")}, tc.opts...)...)
			if err := tc.call(client.TemplateVersion()); err == nil {
				t.Error("expected local validation error")
			}
		})
	}
}

func TestTemplateVersionHTTPError(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.Header().Set("Retry-After", "7")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = io.WriteString(w, `{"detail":"later"}`)
	}, WithOrg("acme"), WithProject("web"))
	_, err := client.TemplateVersion().PruneForProject(context.Background(), "base", PruneTemplateVersionsForProjectParams{})
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 429 || apiErr.Header.Get("Retry-After") != "7" || apiErr.Problem == nil {
		t.Errorf("error = %v", err)
	}
}

func TestTemplateVersionDeleteRequiresNoContent(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{}`)
	}, WithOrg("acme"))
	err := client.TemplateVersion().DeleteForOrg(context.Background(), "base", "v1", nil)
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusOK {
		t.Errorf("error = %v; want HTTP 200 status error", err)
	}
}

func TestTemplateVersionScopeDefaults(t *testing.T) {
	t.Setenv("TEKTONA_ORG", "env-org")
	t.Setenv("TEKTONA_PROJECT", "env-project")
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/orgs/option-org/projects/override/templates/base/versions" {
			t.Errorf("path = %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{}`)
	}, WithOrg("option-org"))
	params := &ListTemplateVersionsForProjectParams{Project: "override"}
	if _, err := client.TemplateVersion().ListForProject(context.Background(), "base", params); err != nil {
		t.Fatal(err)
	}
	if params.Org != "" || params.Project != "override" {
		t.Errorf("input changed: %+v", params)
	}

	blocked := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	}, WithOrg(""), WithProject(""))
	if _, err := blocked.TemplateVersion().ListForProject(context.Background(), "base", nil); err == nil {
		t.Error("empty options inherited environment scope")
	}
}
