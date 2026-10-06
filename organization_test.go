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

func TestOrganizationMethods(t *testing.T) {
	t.Parallel()
	cursor, limit, filter := "next", int32(2), "alice"
	location := "eu"
	role := api.UpdateOrgInputBodyDefaultProjectRole("reader")
	for _, tc := range []struct {
		name, method, path, query, body, response string
		status                                    int
		call                                      func(Organization) (any, error)
	}{
		{"list", "GET", "/v1/orgs", "cursor=next&limit=2", "", `{"items":[]}`, 200, func(o Organization) (any, error) {
			return o.List(context.Background(), &ListOrganizationsParams{Cursor: &cursor, Limit: &limit})
		}},
		{"create", "POST", "/v1/orgs", "", `{"display_name":"New Org","name":"new-org"}`, `{"name":"new-org"}`, 201, func(o Organization) (any, error) {
			return o.Create(context.Background(), CreateOrganizationParams{Name: "new-org", DisplayName: "New Org"})
		}},
		{"get", "GET", "/v1/orgs/acme", "", "", `{"name":"acme"}`, 200, func(o Organization) (any, error) {
			return o.Get(context.Background(), "acme")
		}},
		{"get override", "GET", "/v1/orgs/other", "", "", `{"name":"other"}`, 200, func(o Organization) (any, error) {
			return o.Get(context.Background(), "other")
		}},
		{"update", "PUT", "/v1/orgs/acme", "", `{"default_location_id":"eu","default_project_role":"reader","display_name":"Renamed"}`, `{"name":"acme"}`, 200, func(o Organization) (any, error) {
			return o.Update(context.Background(), "acme", UpdateOrganizationParams{DisplayName: "Renamed", DefaultLocationId: &location, DefaultProjectRole: &role})
		}},
		{"members", "GET", "/v1/orgs/acme/members", "q=alice&cursor=next&limit=2", "", `{"items":[]}`, 200, func(o Organization) (any, error) {
			return o.ListMembers(context.Background(), "", &ListOrganizationMembersParams{Q: &filter, Cursor: &cursor, Limit: &limit})
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
			if client.Organization() != client.Organization() {
				t.Error("accessor returned a different interface")
			}
			got, err := tc.call(client.Organization())
			if err != nil || got == nil {
				t.Fatalf("result = %v, %v", got, err)
			}
		})
	}
}

func TestOrganizationValidation(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	})
	for _, call := range []func() error{
		func() error {
			_, err := client.Organization().Create(context.Background(), CreateOrganizationParams{})
			return err
		},
		func() error { _, err := client.Organization().Get(context.Background(), ""); return err },
		func() error {
			_, err := client.Organization().Update(context.Background(), "", UpdateOrganizationParams{DisplayName: "name"})
			return err
		},
		func() error { _, err := client.Organization().ListMembers(context.Background(), "", nil); return err },
	} {
		if err := call(); err == nil {
			t.Error("expected local validation error")
		}
	}
}

func TestOrganizationResponses(t *testing.T) {
	t.Parallel()
	for _, call := range []struct {
		name string
		fn   func(Organization) error
	}{
		{"list", func(o Organization) error { _, err := o.List(context.Background(), nil); return err }},
		{"create", func(o Organization) error {
			_, err := o.Create(context.Background(), CreateOrganizationParams{Name: "a", DisplayName: "A"})
			return err
		}},
		{"get", func(o Organization) error { _, err := o.Get(context.Background(), "acme"); return err }},
		{"update", func(o Organization) error {
			_, err := o.Update(context.Background(), "acme", UpdateOrganizationParams{DisplayName: "A"})
			return err
		}},
		{"members", func(o Organization) error { _, err := o.ListMembers(context.Background(), "", nil); return err }},
	} {
		t.Run(call.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/problem+json")
				w.Header().Set("Retry-After", "5")
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = io.WriteString(w, `{"detail":"later"}`)
			}, WithOrg("acme"))
			err := call.fn(client.Organization())
			var apiErr *APIError
			if !errors.As(err, &apiErr) || apiErr.StatusCode != 429 || apiErr.Header.Get("Retry-After") != "5" || apiErr.Problem == nil {
				t.Errorf("error = %v", err)
			}
		})
	}
}

func TestOrganizationScopePrecedence(t *testing.T) {
	t.Setenv("TEKTONA_ORG", "env-org")
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/orgs/option-org/members" {
			t.Errorf("path = %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{}`)
	}, WithOrg("option-org"))
	if _, err := client.Organization().ListMembers(context.Background(), "", nil); err != nil {
		t.Fatal(err)
	}
	blocked := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	}, WithOrg(""))
	if _, err := blocked.Organization().ListMembers(context.Background(), "", nil); err == nil {
		t.Error("explicit empty option inherited environment")
	}
}

func TestOrganizationRejectsEmptyResourceName(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	}, WithOrg("production"))
	if _, err := client.Organization().Get(context.Background(), ""); err == nil {
		t.Error("get accepted an empty organization name")
	}
	if _, err := client.Organization().Update(context.Background(), "", UpdateOrganizationParams{DisplayName: "Production"}); err == nil {
		t.Error("update accepted an empty organization name")
	}
}
