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

const testRepositoryID = "01234567-89ab-4cde-8f01-23456789abcd"

func TestRepositoryMethods(t *testing.T) {
	t.Parallel()
	defaultOnly, isDefault, refresh := true, false, true
	branch, query, cursor, limit := "main", "feat", "next", int32(3)
	for _, tc := range []struct {
		name, method, path, query, body, response string
		status                                    int
		call                                      func(Repository) (any, error)
	}{
		{"list", "GET", "/v1/orgs/acme/projects/web/repos", "default=true", "", `{"items":[]}`, 200, func(r Repository) (any, error) {
			return r.List(context.Background(), &ListRepositoriesParams{Default: &defaultOnly})
		}},
		{"create", "POST", "/v1/orgs/other/projects/app/repos", "", `{"name":"source","url":"https://example.com/source.git","default_branch":"main","is_default":false}`, `{"id":"` + testRepositoryID + `"}`, 201, func(r Repository) (any, error) {
			return r.Create(context.Background(), CreateRepositoryParams{Org: "other", Project: "app", Name: "source", Url: "https://example.com/source.git", DefaultBranch: &branch, IsDefault: &isDefault})
		}},
		{"update", "PUT", "/v1/orgs/other/projects/app/repos/" + testRepositoryID, "", `{"name":"source","url":"https://example.com/source.git","default_branch":"main","is_default":false}`, `{"id":"` + testRepositoryID + `"}`, 200, func(r Repository) (any, error) {
			return r.Update(context.Background(), testRepositoryID, UpdateRepositoryParams{Org: "other", Project: "app", Name: "source", Url: "https://example.com/source.git", DefaultBranch: &branch, IsDefault: &isDefault})
		}},
		{"delete", "DELETE", "/v1/orgs/acme/projects/web/repos/" + testRepositoryID, "", "", "", 204, func(r Repository) (any, error) {
			return true, r.Delete(context.Background(), testRepositoryID, nil)
		}},
		{"branches", "GET", "/v1/orgs/acme/projects/web/repos/" + testRepositoryID + "/branches", "refresh=true&q=feat&cursor=next&limit=3", "", `{"items":[{"name":"feature"}]}`, 200, func(r Repository) (any, error) {
			return r.ListBranches(context.Background(), testRepositoryID, &ListRepositoryBranchesParams{Refresh: &refresh, Q: &query, Cursor: &cursor, Limit: &limit})
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
			if client.Repository() != client.Repository() {
				t.Error("accessor returned a different interface")
			}
			got, err := tc.call(client.Repository())
			if err != nil || got == nil {
				t.Fatalf("result = %v, %v", got, err)
			}
		})
	}
}

func TestRepositoryValidation(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	for _, tc := range []struct {
		name string
		opts []Option
		call func(Repository) error
	}{
		{"missing scope", nil, func(r Repository) error { _, err := r.List(ctx, nil); return err }},
		{"missing create fields", []Option{WithOrg("acme"), WithProject("web")}, func(r Repository) error { _, err := r.Create(ctx, CreateRepositoryParams{}); return err }},
		{"invalid update id", []Option{WithOrg("acme"), WithProject("web")}, func(r Repository) error {
			_, err := r.Update(ctx, "invalid", UpdateRepositoryParams{Name: "name", Url: "https://example.com"})
			return err
		}},
		{"missing update fields", []Option{WithOrg("acme"), WithProject("web")}, func(r Repository) error {
			_, err := r.Update(ctx, testRepositoryID, UpdateRepositoryParams{})
			return err
		}},
		{"missing delete id", []Option{WithOrg("acme"), WithProject("web")}, func(r Repository) error { return r.Delete(ctx, "", nil) }},
		{"missing branches id", []Option{WithOrg("acme"), WithProject("web")}, func(r Repository) error { _, err := r.ListBranches(ctx, "", nil); return err }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { t.Errorf("unexpected request: %s %s", r.Method, r.URL) }, tc.opts...)
			if err := tc.call(client.Repository()); err == nil {
				t.Error("expected local error")
			}
		})
	}
}

func TestRepositoryScopeAndError(t *testing.T) {
	t.Setenv("TEKTONA_ORG", "env-org")
	t.Setenv("TEKTONA_PROJECT", "env-project")
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/orgs/option-org/projects/env-project/repos" {
			t.Errorf("path = %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/problem+json")
		w.Header().Set("Retry-After", "5")
		w.WriteHeader(429)
		_, _ = io.WriteString(w, `{"detail":"later"}`)
	}, WithOrg("option-org"))
	_, err := client.Repository().List(context.Background(), nil)
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 429 || apiErr.Header.Get("Retry-After") != "5" || apiErr.Problem == nil {
		t.Errorf("error = %v", err)
	}
	blocked := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { t.Errorf("unexpected request: %s %s", r.Method, r.URL) }, WithOrg(""), WithProject(""))
	if _, err := blocked.Repository().List(context.Background(), nil); err == nil {
		t.Error("explicit empty scope inherited environment")
	}
}
