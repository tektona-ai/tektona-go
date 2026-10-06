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

func TestRegistryMethods(t *testing.T) {
	t.Parallel()
	cursor, limit, dryRun := "next", int32(2), true
	namespace, password, token, username := "team", "password", "token", "user"
	for _, tc := range []struct {
		name, method, path, query, body, response string
		status                                    int
		call                                      func(Registry) (any, error)
	}{
		{"list", "GET", "/v1/orgs/acme/projects/web/registries", "cursor=next&limit=2", "", `{"items":[]}`, 200, func(r Registry) (any, error) {
			return r.List(context.Background(), &ListRegistriesParams{Cursor: &cursor, Limit: &limit})
		}},
		{"create", "POST", "/v1/orgs/other/projects/app/registries", "dry_run=true", `{"auth_type":"basic","display_name":"Images","endpoint":"ghcr.io","name":"images","namespace":"team","password":"password","token":"token","username":"user"}`, `{"name":"images"}`, 201, func(r Registry) (any, error) {
			return r.Create(context.Background(), CreateRegistryParams{Org: "other", Project: "app", DryRun: &dryRun, AuthType: api.CreateProjectRegistryInputBodyAuthTypeBasic, DisplayName: "Images", Endpoint: "ghcr.io", Name: "images", Namespace: &namespace, Password: &password, Token: &token, Username: &username})
		}},
		{"get", "GET", "/v1/orgs/acme/projects/web/registries/images", "", "", `{"name":"images"}`, 200, func(r Registry) (any, error) {
			return r.Get(context.Background(), "images", nil)
		}},
		{"update", "PUT", "/v1/orgs/other/projects/app/registries/images", "dry_run=true", `{"auth_type":"bearer","display_name":"Images","endpoint":"ghcr.io","name":"images","namespace":"team","password":"password","token":"token","username":"user"}`, `{"name":"images"}`, 200, func(r Registry) (any, error) {
			return r.Update(context.Background(), "images", UpdateRegistryParams{Org: "other", Project: "app", DryRun: &dryRun, AuthType: api.UpdateProjectRegistryInputBodyAuthTypeBearer, DisplayName: "Images", Endpoint: "ghcr.io", Name: "images", Namespace: &namespace, Password: &password, Token: &token, Username: &username})
		}},
		{"delete", "DELETE", "/v1/orgs/acme/projects/web/registries/images", "", "", "", 204, func(r Registry) (any, error) {
			return true, r.Delete(context.Background(), "images", nil)
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
			if client.Registry() != client.Registry() {
				t.Error("accessor returned a different interface")
			}
			got, err := tc.call(client.Registry())
			if err != nil || got == nil {
				t.Fatalf("result = %v, %v", got, err)
			}
		})
	}
}

func TestRegistryValidation(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	for _, tc := range []struct {
		name string
		opts []Option
		call func(Registry) error
	}{
		{"missing scope", nil, func(r Registry) error { _, err := r.List(ctx, nil); return err }},
		{"missing create fields", []Option{WithOrg("acme"), WithProject("web")}, func(r Registry) error { _, err := r.Create(ctx, CreateRegistryParams{}); return err }},
		{"missing get name", []Option{WithOrg("acme"), WithProject("web")}, func(r Registry) error { _, err := r.Get(ctx, "", nil); return err }},
		{"missing update fields", []Option{WithOrg("acme"), WithProject("web")}, func(r Registry) error { _, err := r.Update(ctx, "images", UpdateRegistryParams{}); return err }},
		{"missing delete name", []Option{WithOrg("acme"), WithProject("web")}, func(r Registry) error { return r.Delete(ctx, "", nil) }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { t.Errorf("unexpected request: %s %s", r.Method, r.URL) }, tc.opts...)
			if err := tc.call(client.Registry()); err == nil {
				t.Error("expected local error")
			}
		})
	}
}

func TestRegistryScopeAndError(t *testing.T) {
	t.Setenv("TEKTONA_ORG", "env-org")
	t.Setenv("TEKTONA_PROJECT", "env-project")
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/orgs/option-org/projects/env-project/registries" {
			t.Errorf("path = %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/problem+json")
		w.Header().Set("Retry-After", "5")
		w.WriteHeader(429)
		_, _ = io.WriteString(w, `{"detail":"later"}`)
	}, WithOrg("option-org"))
	_, err := client.Registry().List(context.Background(), nil)
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 429 || apiErr.Header.Get("Retry-After") != "5" || apiErr.Problem == nil {
		t.Errorf("error = %v", err)
	}
	blocked := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { t.Errorf("unexpected request: %s %s", r.Method, r.URL) }, WithOrg(""), WithProject(""))
	if _, err := blocked.Registry().List(context.Background(), nil); err == nil {
		t.Error("explicit empty scope inherited environment")
	}
}
