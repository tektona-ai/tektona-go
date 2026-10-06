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

func TestEgressProxyProfileMethods(t *testing.T) {
	t.Parallel()
	cursor, limit, label, pathMatch := "next", int32(3), "Team", "/v1"
	recipe := []api.InjectOp{{Target: api.InjectOpTarget("header"), Header: &api.HeaderOp{Name: "x-api-key", Template: "${secret:key}"}}}
	for _, tc := range []struct {
		name, method, path, query, body, response string
		status                                    int
		call                                      func(EgressProxyProfile) (any, error)
	}{
		{"list org", "GET", "/v1/orgs/acme/egress-proxy-profiles", "cursor=next&limit=3", "", `{"items":[],"pagination":{}}`, 200, func(p EgressProxyProfile) (any, error) {
			return p.ListForOrg(context.Background(), "", &ListEgressProxyProfilesForOrgParams{Cursor: &cursor, Limit: &limit})
		}},
		{"create org", "POST", "/v1/orgs/other/egress-proxy-profiles", "", `{"name":"team","display_name":"Team","default":false,"scope":"org"}`, `{"profile":{}}`, 201, func(p EgressProxyProfile) (any, error) {
			return p.CreateForOrg(context.Background(), CreateEgressProxyProfileForOrgParams{Org: "other", Name: "team", DisplayName: &label})
		}},
		{"list project", "GET", "/v1/orgs/other/projects/web/egress-proxy-profiles", "cursor=next&limit=3", "", `{"items":[],"pagination":{}}`, 200, func(p EgressProxyProfile) (any, error) {
			return p.ListForProject(context.Background(), "", &ListEgressProxyProfilesForProjectParams{Org: "other", Cursor: &cursor, Limit: &limit})
		}},
		{"create project", "POST", "/v1/orgs/other/projects/dev/egress-proxy-profiles", "", `{"name":"team","display_name":"Team","default":true,"scope":"project"}`, `{"profile":{}}`, 201, func(p EgressProxyProfile) (any, error) {
			return p.CreateForProject(context.Background(), CreateEgressProxyProfileForProjectParams{Org: "other", Project: "dev", Name: "team", DisplayName: &label, Default: true})
		}},
		{"add rule", "POST", "/v1/orgs/acme/projects/web/egress-proxy-profiles/team/rules", "", `{"domain_pattern":"api.example.com","path_match":"/v1","recipe":[{"target":"header","header":{"name":"x-api-key","template":"${secret:key}"}}]}`, `{"rule_id":"rule-id"}`, 201, func(p EgressProxyProfile) (any, error) {
			return p.AddRuleForProject(context.Background(), "team", AddEgressProxyProfileRuleForProjectParams{DomainPattern: "api.example.com", PathMatch: &pathMatch, Recipe: &recipe})
		}},
		{"delete profile", "DELETE", "/v1/orgs/acme/projects/web/egress-proxy-profiles/profile-id", "", "", "", 204, func(p EgressProxyProfile) (any, error) {
			return true, p.DeleteForProject(context.Background(), "profile-id", nil)
		}},
		{"set default", "PUT", "/v1/orgs/acme/projects/web/egress-proxy-profiles/profile-id/default", "", `{"is_default":true}`, "", 204, func(p EgressProxyProfile) (any, error) {
			return true, p.SetDefaultForProject(context.Background(), "profile-id", SetEgressProxyProfileDefaultForProjectParams{IsDefault: true})
		}},
		{"delete rule", "DELETE", "/v1/orgs/acme/projects/web/egress-proxy-profiles/profile-id/rules/rule-id", "", "", "", 204, func(p EgressProxyProfile) (any, error) {
			return true, p.DeleteRuleForProject(context.Background(), "profile-id", "rule-id", nil)
		}},
		{"update rule", "PUT", "/v1/orgs/acme/projects/web/egress-proxy-profiles/profile-id/rules/rule-id", "", `{"domain_pattern":"api.example.com","path_match":"/v1","recipe":[{"target":"header","header":{"name":"x-api-key","template":"${secret:key}"}}]}`, "", 204, func(p EgressProxyProfile) (any, error) {
			return true, p.UpdateRuleForProject(context.Background(), "profile-id", "rule-id", UpdateEgressProxyProfileRuleForProjectParams{DomainPattern: "api.example.com", PathMatch: &pathMatch, Recipe: &recipe})
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
			got, err := tc.call(client.EgressProxyProfile())
			if err != nil || got == nil {
				t.Fatalf("result = %v, %v", got, err)
			}
		})
	}
}

func TestEgressProxyProfileValidation(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		opts []Option
		call func(EgressProxyProfile) error
	}{
		{"list org", nil, func(p EgressProxyProfile) error { _, err := p.ListForOrg(context.Background(), "", nil); return err }},
		{"create org", nil, func(p EgressProxyProfile) error {
			_, err := p.CreateForOrg(context.Background(), CreateEgressProxyProfileForOrgParams{Name: "x"})
			return err
		}},
		{"create org name", []Option{WithOrg("a")}, func(p EgressProxyProfile) error {
			_, err := p.CreateForOrg(context.Background(), CreateEgressProxyProfileForOrgParams{})
			return err
		}},
		{"list project", []Option{WithOrg("a")}, func(p EgressProxyProfile) error {
			_, err := p.ListForProject(context.Background(), "", nil)
			return err
		}},
		{"create project", []Option{WithOrg("a"), WithProject("b")}, func(p EgressProxyProfile) error {
			_, err := p.CreateForProject(context.Background(), CreateEgressProxyProfileForProjectParams{})
			return err
		}},
		{"add rule name", []Option{WithOrg("a"), WithProject("b")}, func(p EgressProxyProfile) error {
			_, err := p.AddRuleForProject(context.Background(), "", AddEgressProxyProfileRuleForProjectParams{})
			return err
		}},
		{"delete profile ID", []Option{WithOrg("a"), WithProject("b")}, func(p EgressProxyProfile) error { return p.DeleteForProject(context.Background(), "", nil) }},
		{"default profile ID", []Option{WithOrg("a"), WithProject("b")}, func(p EgressProxyProfile) error {
			return p.SetDefaultForProject(context.Background(), "", SetEgressProxyProfileDefaultForProjectParams{})
		}},
		{"delete rule ID", []Option{WithOrg("a"), WithProject("b")}, func(p EgressProxyProfile) error {
			return p.DeleteRuleForProject(context.Background(), "profile", "", nil)
		}},
		{"update rule ID", []Option{WithOrg("a"), WithProject("b")}, func(p EgressProxyProfile) error {
			return p.UpdateRuleForProject(context.Background(), "profile", "", UpdateEgressProxyProfileRuleForProjectParams{})
		}},
		{"project scope", []Option{WithOrg("a")}, func(p EgressProxyProfile) error { return p.DeleteForProject(context.Background(), "profile", nil) }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				t.Errorf("unexpected request: %s %s", r.Method, r.URL)
			}, tc.opts...)
			if err := tc.call(client.EgressProxyProfile()); err == nil {
				t.Error("expected a local error")
			}
		})
	}
}

func TestEgressProxyProfileResponses(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		call func(EgressProxyProfile) error
	}{
		{"list org", func(p EgressProxyProfile) error { _, err := p.ListForOrg(context.Background(), "", nil); return err }},
		{"create org", func(p EgressProxyProfile) error {
			_, err := p.CreateForOrg(context.Background(), CreateEgressProxyProfileForOrgParams{Name: "team"})
			return err
		}},
		{"list project", func(p EgressProxyProfile) error {
			_, err := p.ListForProject(context.Background(), "", nil)
			return err
		}},
		{"create project", func(p EgressProxyProfile) error {
			_, err := p.CreateForProject(context.Background(), CreateEgressProxyProfileForProjectParams{Name: "team"})
			return err
		}},
		{"add rule", func(p EgressProxyProfile) error {
			_, err := p.AddRuleForProject(context.Background(), "team", AddEgressProxyProfileRuleForProjectParams{})
			return err
		}},
		{"delete profile", func(p EgressProxyProfile) error { return p.DeleteForProject(context.Background(), "id", nil) }},
		{"set default", func(p EgressProxyProfile) error {
			return p.SetDefaultForProject(context.Background(), "id", SetEgressProxyProfileDefaultForProjectParams{})
		}},
		{"delete rule", func(p EgressProxyProfile) error {
			return p.DeleteRuleForProject(context.Background(), "id", "rule", nil)
		}},
		{"update rule", func(p EgressProxyProfile) error {
			return p.UpdateRuleForProject(context.Background(), "id", "rule", UpdateEgressProxyProfileRuleForProjectParams{})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/problem+json")
				w.Header().Set("Retry-After", "5")
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = io.WriteString(w, `{"detail":"later"}`)
			}, WithOrg("acme"), WithProject("web"))
			var apiErr *APIError
			if err := tc.call(client.EgressProxyProfile()); !errors.As(err, &apiErr) || apiErr.StatusCode != 429 || apiErr.Header.Get("Retry-After") != "5" || apiErr.Problem == nil {
				t.Errorf("error = %v", err)
			}
		})
	}
}

func TestEgressProxyProfileScopePrecedence(t *testing.T) {
	t.Setenv("TEKTONA_ORG", "environment")
	t.Setenv("TEKTONA_PROJECT", "environment")
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/orgs/call/projects/call/egress-proxy-profiles" {
			t.Errorf("path = %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"items":[],"pagination":{}}`)
	}, WithOrg("option"), WithProject("option"))
	if _, err := client.EgressProxyProfile().ListForProject(context.Background(), "call", &ListEgressProxyProfilesForProjectParams{Org: "call"}); err != nil {
		t.Fatal(err)
	}
	blocked := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	}, WithOrg(""), WithProject(""))
	if _, err := blocked.EgressProxyProfile().ListForProject(context.Background(), "", nil); err == nil {
		t.Error("explicit empty options inherited environment")
	}
}

func TestEgressProxyProfileUnexpectedStatus(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{}`)
	}, WithOrg("acme"), WithProject("web"))
	var apiErr *APIError
	if err := client.EgressProxyProfile().DeleteForProject(context.Background(), "id", nil); !errors.As(err, &apiErr) || apiErr.StatusCode != 200 {
		t.Errorf("error = %v", err)
	}
}
