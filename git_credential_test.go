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

func TestGitCredentialMethods(t *testing.T) {
	t.Parallel()
	scope := api.ListOrgProjectGitCredentialsParamsScopePersonal
	repos := []string{"repository-id"}
	token := "new-token"
	for _, tc := range []struct {
		name, method, path, query, body, response string
		status                                    int
		call                                      func(GitCredential) (any, error)
	}{
		{"list", "GET", "/v1/orgs/acme/projects/web/git-credentials", "scope=personal", "", `{"items":[]}`, 200, func(g GitCredential) (any, error) {
			return g.List(context.Background(), &ListGitCredentialsParams{Scope: &scope})
		}},
		{"create", "POST", "/v1/orgs/other/projects/app/git-credentials", "", `{"name":"source","display_name":"Source","forge":"github","scope":"project","token":"secret","repository_ids":["repository-id"]}`, `{"credential":{"name":"source"}}`, 201, func(g GitCredential) (any, error) {
			return g.Create(context.Background(), CreateGitCredentialParams{Org: "other", Project: "app", Name: "source", DisplayName: "Source", Forge: api.CreateGitCredentialBodyForgeGithub, Scope: api.CreateGitCredentialBodyScopeProject, Token: "secret", RepositoryIds: &repos})
		}},
		{"update", "PUT", "/v1/orgs/other/projects/app/git-credentials/cred-1", "scope=project", `{"display_name":"Source","forge":"gitlab","repository_ids":["repository-id"],"token":"new-token"}`, `{"credential":{"name":"source"}}`, 200, func(g GitCredential) (any, error) {
			return g.Update(context.Background(), "cred-1", UpdateGitCredentialParams{Org: "other", Project: "app", Scope: api.UpdateOrgProjectGitCredentialParamsScopeProject, DisplayName: "Source", Forge: api.UpdateGitCredentialBodyForgeGitlab, RepositoryIds: &repos, Token: &token})
		}},
		{"delete", "DELETE", "/v1/orgs/acme/projects/web/git-credentials/cred-1", "scope=personal", "", "", 204, func(g GitCredential) (any, error) {
			return true, g.Delete(context.Background(), "cred-1", DeleteGitCredentialParams{Scope: api.DeleteOrgProjectGitCredentialParamsScopePersonal})
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
			if client.GitCredential() != client.GitCredential() {
				t.Error("accessor returned a different interface")
			}
			got, err := tc.call(client.GitCredential())
			if err != nil || got == nil {
				t.Fatalf("result = %v, %v", got, err)
			}
		})
	}
}

func TestGitCredentialValidation(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	for _, tc := range []struct {
		name string
		opts []Option
		call func(GitCredential) error
	}{
		{"missing scope", nil, func(g GitCredential) error { _, err := g.List(ctx, nil); return err }},
		{"missing create fields", []Option{WithOrg("acme"), WithProject("web")}, func(g GitCredential) error { _, err := g.Create(ctx, CreateGitCredentialParams{}); return err }},
		{"missing update id", []Option{WithOrg("acme"), WithProject("web")}, func(g GitCredential) error { _, err := g.Update(ctx, "", UpdateGitCredentialParams{}); return err }},
		{"missing update scope", []Option{WithOrg("acme"), WithProject("web")}, func(g GitCredential) error {
			_, err := g.Update(ctx, "id", UpdateGitCredentialParams{DisplayName: "Name", Forge: api.UpdateGitCredentialBodyForgeGithub})
			return err
		}},
		{"missing delete scope", []Option{WithOrg("acme"), WithProject("web")}, func(g GitCredential) error { return g.Delete(ctx, "id", DeleteGitCredentialParams{}) }},
		{"missing delete id", []Option{WithOrg("acme"), WithProject("web")}, func(g GitCredential) error {
			return g.Delete(ctx, "", DeleteGitCredentialParams{Scope: api.DeleteOrgProjectGitCredentialParamsScopeProject})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { t.Errorf("unexpected request: %s %s", r.Method, r.URL) }, tc.opts...)
			if err := tc.call(client.GitCredential()); err == nil {
				t.Error("expected local error")
			}
		})
	}
}

func TestGitCredentialScopeAndError(t *testing.T) {
	t.Setenv("TEKTONA_ORG", "env-org")
	t.Setenv("TEKTONA_PROJECT", "env-project")
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/orgs/option-org/projects/env-project/git-credentials" {
			t.Errorf("path = %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/problem+json")
		w.Header().Set("Retry-After", "5")
		w.WriteHeader(429)
		_, _ = io.WriteString(w, `{"detail":"later"}`)
	}, WithOrg("option-org"))
	_, err := client.GitCredential().List(context.Background(), nil)
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 429 || apiErr.Header.Get("Retry-After") != "5" || apiErr.Problem == nil {
		t.Errorf("error = %v", err)
	}
	blocked := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { t.Errorf("unexpected request: %s %s", r.Method, r.URL) }, WithOrg(""), WithProject(""))
	if _, err := blocked.GitCredential().List(context.Background(), nil); err == nil {
		t.Error("explicit empty scope inherited environment")
	}
}
