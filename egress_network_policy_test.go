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

func TestEgressNetworkPolicyMethods(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	domains := []string{"github.com"}
	allowed := []string{"10.0.0.0/8"}
	denied := []string{"169.254.0.0/16"}
	blocked := []string{"evil.example.com"}
	display, extends := "Build access", "tektona/dev"
	create := `{"name":"build","allowed_cidrs":["10.0.0.0/8"],"allowed_domains":["github.com"],"denied_cidrs":["169.254.0.0/16"],"denied_domains":["evil.example.com"],"display_name":"Build access","extends":"tektona/dev"}`
	update := `{"allowed_cidrs":["10.0.0.0/8"],"allowed_domains":["github.com"],"denied_cidrs":["169.254.0.0/16"],"denied_domains":["evil.example.com"],"display_name":"Build access","extends":"tektona/dev"}`
	for _, tc := range []struct {
		name, method, path, body, response string
		status                             int
		call                               func(EgressNetworkPolicy) (any, error)
	}{
		{"list system", "GET", "/v1/system-egress-network-policies", "", `{"items":[],"system_denies":{}}`, 200, func(p EgressNetworkPolicy) (any, error) { return p.ListSystem(ctx) }},
		{"get system", "GET", "/v1/system-egress-network-policies/open", "", `{"name":"open"}`, 200, func(p EgressNetworkPolicy) (any, error) { return p.GetSystem(ctx, "open") }},
		{"list org default", "GET", "/v1/orgs/acme/egress-network-policies", "", `{"items":[]}`, 200, func(p EgressNetworkPolicy) (any, error) { return p.ListForOrg(ctx, "") }},
		{"create org", "POST", "/v1/orgs/other/egress-network-policies", create, `{"name":"build"}`, 201, func(p EgressNetworkPolicy) (any, error) {
			return p.CreateForOrg(ctx, CreateEgressNetworkPolicyForOrgParams{Org: "other", Name: "build", AllowedCidrs: &allowed, AllowedDomains: &domains, DeniedCidrs: &denied, DeniedDomains: &blocked, DisplayName: &display, Extends: &extends})
		}},
		{"get org default", "GET", "/v1/orgs/acme/egress-network-policies/build", "", `{"name":"build"}`, 200, func(p EgressNetworkPolicy) (any, error) { return p.GetForOrg(ctx, "build", nil) }},
		{"get org override", "GET", "/v1/orgs/other/egress-network-policies/build", "", `{"name":"build"}`, 200, func(p EgressNetworkPolicy) (any, error) {
			return p.GetForOrg(ctx, "build", &GetEgressNetworkPolicyForOrgParams{Org: "other"})
		}},
		{"update org", "PUT", "/v1/orgs/other/egress-network-policies/build", update, `{"name":"build"}`, 200, func(p EgressNetworkPolicy) (any, error) {
			return p.UpdateForOrg(ctx, "build", UpdateEgressNetworkPolicyForOrgParams{Org: "other", AllowedCidrs: &allowed, AllowedDomains: &domains, DeniedCidrs: &denied, DeniedDomains: &blocked, DisplayName: &display, Extends: &extends})
		}},
		{"delete org", "DELETE", "/v1/orgs/other/egress-network-policies/build", "", "", 204, func(p EgressNetworkPolicy) (any, error) {
			return true, p.DeleteForOrg(ctx, "build", &DeleteEgressNetworkPolicyForOrgParams{Org: "other"})
		}},
		{"list project default", "GET", "/v1/orgs/acme/projects/web/egress-network-policies", "", `{"items":[]}`, 200, func(p EgressNetworkPolicy) (any, error) { return p.ListForProject(ctx, "", "") }},
		{"list project override", "GET", "/v1/orgs/other/projects/app/egress-network-policies", "", `{"items":[]}`, 200, func(p EgressNetworkPolicy) (any, error) { return p.ListForProject(ctx, "other", "app") }},
		{"create project", "POST", "/v1/orgs/other/projects/app/egress-network-policies", create, `{"name":"build"}`, 201, func(p EgressNetworkPolicy) (any, error) {
			return p.CreateForProject(ctx, CreateEgressNetworkPolicyForProjectParams{Org: "other", Project: "app", Name: "build", AllowedCidrs: &allowed, AllowedDomains: &domains, DeniedCidrs: &denied, DeniedDomains: &blocked, DisplayName: &display, Extends: &extends})
		}},
		{"get project default", "GET", "/v1/orgs/acme/projects/web/egress-network-policies/build", "", `{"name":"build"}`, 200, func(p EgressNetworkPolicy) (any, error) { return p.GetForProject(ctx, "build", nil) }},
		{"get project override", "GET", "/v1/orgs/other/projects/app/egress-network-policies/build", "", `{"name":"build"}`, 200, func(p EgressNetworkPolicy) (any, error) {
			return p.GetForProject(ctx, "build", &GetEgressNetworkPolicyForProjectParams{Org: "other", Project: "app"})
		}},
		{"update project", "PUT", "/v1/orgs/other/projects/app/egress-network-policies/build", update, `{"name":"build"}`, 200, func(p EgressNetworkPolicy) (any, error) {
			return p.UpdateForProject(ctx, "build", UpdateEgressNetworkPolicyForProjectParams{Org: "other", Project: "app", AllowedCidrs: &allowed, AllowedDomains: &domains, DeniedCidrs: &denied, DeniedDomains: &blocked, DisplayName: &display, Extends: &extends})
		}},
		{"delete project", "DELETE", "/v1/orgs/other/projects/app/egress-network-policies/build", "", "", 204, func(p EgressNetworkPolicy) (any, error) {
			return true, p.DeleteForProject(ctx, "build", &DeleteEgressNetworkPolicyForProjectParams{Org: "other", Project: "app"})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tc.method || r.URL.Path != tc.path || r.URL.RawQuery != "" {
					t.Errorf("request = %s %s; want %s %s", r.Method, r.URL, tc.method, tc.path)
				}
				if r.Header.Get("Authorization") != "Bearer test-key" {
					t.Error("missing bearer token")
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
			if client.EgressNetworkPolicy() == nil {
				t.Fatal("missing accessor")
			}
			result, err := tc.call(client.EgressNetworkPolicy())
			if err != nil || result == nil {
				t.Fatalf("result = %v, %v", result, err)
			}
		})
	}
}

func TestEgressNetworkPolicyValidation(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	})
	p := client.EgressNetworkPolicy()
	for _, tc := range []struct {
		name string
		call func() error
	}{
		{"system name", func() error { _, err := p.GetSystem(ctx, ""); return err }},
		{"org list", func() error { _, err := p.ListForOrg(ctx, ""); return err }},
		{"org create scope", func() error {
			_, err := p.CreateForOrg(ctx, CreateEgressNetworkPolicyForOrgParams{Name: "x"})
			return err
		}},
		{"org create name", func() error {
			_, err := p.CreateForOrg(ctx, CreateEgressNetworkPolicyForOrgParams{Org: "a"})
			return err
		}},
		{"org get scope", func() error { _, err := p.GetForOrg(ctx, "x", nil); return err }},
		{"org get name", func() error {
			_, err := p.GetForOrg(ctx, "", &GetEgressNetworkPolicyForOrgParams{Org: "a"})
			return err
		}},
		{"org update scope", func() error { _, err := p.UpdateForOrg(ctx, "x", UpdateEgressNetworkPolicyForOrgParams{}); return err }},
		{"org update name", func() error {
			_, err := p.UpdateForOrg(ctx, "", UpdateEgressNetworkPolicyForOrgParams{Org: "a"})
			return err
		}},
		{"org delete scope", func() error { return p.DeleteForOrg(ctx, "x", nil) }},
		{"org delete name", func() error { return p.DeleteForOrg(ctx, "", &DeleteEgressNetworkPolicyForOrgParams{Org: "a"}) }},
		{"project list", func() error { _, err := p.ListForProject(ctx, "", ""); return err }},
		{"project create scope", func() error {
			_, err := p.CreateForProject(ctx, CreateEgressNetworkPolicyForProjectParams{Org: "a", Name: "x"})
			return err
		}},
		{"project create name", func() error {
			_, err := p.CreateForProject(ctx, CreateEgressNetworkPolicyForProjectParams{Org: "a", Project: "b"})
			return err
		}},
		{"project get scope", func() error { _, err := p.GetForProject(ctx, "x", nil); return err }},
		{"project get name", func() error {
			_, err := p.GetForProject(ctx, "", &GetEgressNetworkPolicyForProjectParams{Org: "a", Project: "b"})
			return err
		}},
		{"project update scope", func() error {
			_, err := p.UpdateForProject(ctx, "x", UpdateEgressNetworkPolicyForProjectParams{Org: "a"})
			return err
		}},
		{"project update name", func() error {
			_, err := p.UpdateForProject(ctx, "", UpdateEgressNetworkPolicyForProjectParams{Org: "a", Project: "b"})
			return err
		}},
		{"project delete scope", func() error { return p.DeleteForProject(ctx, "x", nil) }},
		{"project delete name", func() error {
			return p.DeleteForProject(ctx, "", &DeleteEgressNetworkPolicyForProjectParams{Org: "a", Project: "b"})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if err := tc.call(); err == nil {
				t.Error("expected local validation error")
			}
		})
	}
}

func TestEgressNetworkPolicyScopePrecedence(t *testing.T) {
	t.Setenv("TEKTONA_ORG", "env-org")
	t.Setenv("TEKTONA_PROJECT", "env-project")
	ctx := context.Background()
	envClient := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/orgs/env-org/projects/env-project/egress-network-policies" {
			t.Errorf("environment scope path = %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{}`)
	})
	if _, err := envClient.EgressNetworkPolicy().ListForProject(ctx, "", ""); err != nil {
		t.Fatal(err)
	}
	var paths []string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{}`)
	}, WithOrg("option-org"), WithProject("option-project"))
	p := client.EgressNetworkPolicy()
	for _, call := range []func() error{
		func() error { _, err := p.GetForOrg(ctx, "x", nil); return err },
		func() error {
			_, err := p.GetForProject(ctx, "x", &GetEgressNetworkPolicyForProjectParams{Org: "call-org"})
			return err
		},
		func() error { _, err := p.ListForProject(ctx, "", "call-project"); return err },
		func() error { _, err := p.GetSystem(ctx, "x"); return err },
	} {
		if err := call(); err != nil {
			t.Fatal(err)
		}
	}
	want := []string{
		"/v1/orgs/option-org/egress-network-policies/x",
		"/v1/orgs/call-org/projects/option-project/egress-network-policies/x",
		"/v1/orgs/option-org/projects/call-project/egress-network-policies",
		"/v1/system-egress-network-policies/x",
	}
	if !reflect.DeepEqual(paths, want) {
		t.Errorf("paths = %v; want %v", paths, want)
	}
	blocked := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	}, WithOrg(""), WithProject(""))
	if _, err := blocked.EgressNetworkPolicy().ListForProject(ctx, "", ""); err == nil {
		t.Error("explicit empty options inherited environment")
	}
	if _, err := blocked.EgressNetworkPolicy().ListForOrg(ctx, ""); err == nil {
		t.Error("explicit empty org inherited environment")
	}
}

func TestEgressNetworkPolicyResponses(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	for _, tc := range []struct {
		name string
		call func(EgressNetworkPolicy) error
	}{
		{"list system", func(p EgressNetworkPolicy) error { _, err := p.ListSystem(ctx); return err }},
		{"get system", func(p EgressNetworkPolicy) error { _, err := p.GetSystem(ctx, "x"); return err }},
		{"list org", func(p EgressNetworkPolicy) error { _, err := p.ListForOrg(ctx, ""); return err }},
		{"create org", func(p EgressNetworkPolicy) error {
			_, err := p.CreateForOrg(ctx, CreateEgressNetworkPolicyForOrgParams{Name: "x"})
			return err
		}},
		{"get org", func(p EgressNetworkPolicy) error { _, err := p.GetForOrg(ctx, "x", nil); return err }},
		{"update org", func(p EgressNetworkPolicy) error {
			_, err := p.UpdateForOrg(ctx, "x", UpdateEgressNetworkPolicyForOrgParams{})
			return err
		}},
		{"delete org", func(p EgressNetworkPolicy) error { return p.DeleteForOrg(ctx, "x", nil) }},
		{"list project", func(p EgressNetworkPolicy) error { _, err := p.ListForProject(ctx, "", ""); return err }},
		{"create project", func(p EgressNetworkPolicy) error {
			_, err := p.CreateForProject(ctx, CreateEgressNetworkPolicyForProjectParams{Name: "x"})
			return err
		}},
		{"get project", func(p EgressNetworkPolicy) error { _, err := p.GetForProject(ctx, "x", nil); return err }},
		{"update project", func(p EgressNetworkPolicy) error {
			_, err := p.UpdateForProject(ctx, "x", UpdateEgressNetworkPolicyForProjectParams{})
			return err
		}},
		{"delete project", func(p EgressNetworkPolicy) error { return p.DeleteForProject(ctx, "x", nil) }},
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
			if err := tc.call(client.EgressNetworkPolicy()); !errors.As(err, &apiErr) || apiErr.StatusCode != 429 || apiErr.Header.Get("Retry-After") != "5" || apiErr.Problem == nil || string(apiErr.Body) != `{"detail":"later"}` {
				t.Errorf("error = %v", err)
			}
		})
	}
}

func TestEgressNetworkPolicyResponseFailures(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name        string
		status      int
		contentType string
		body        string
		call        func(EgressNetworkPolicy) error
	}{
		{"delete requires 204", 200, "application/json", `{}`, func(p EgressNetworkPolicy) error { return p.DeleteForOrg(context.Background(), "x", nil) }},
		{"json required", 200, "text/plain", "ok", func(p EgressNetworkPolicy) error { _, err := p.GetForOrg(context.Background(), "x", nil); return err }},
		{"valid json required", 200, "application/json", "{", func(p EgressNetworkPolicy) error { _, err := p.GetSystem(context.Background(), "x"); return err }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", tc.contentType)
				w.WriteHeader(tc.status)
				_, _ = io.WriteString(w, tc.body)
			}, WithOrg("acme"))
			if err := tc.call(client.EgressNetworkPolicy()); err == nil {
				t.Error("expected response failure")
			}
		})
	}
}
