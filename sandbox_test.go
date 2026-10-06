package tektona

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/tektona-ai/tektona-go/api"
)

const sandboxID = "01J5K3NDEKTSV4RRFFQ69G5FAV"

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return fn(req) }

func TestBearerTransport(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, target, want string }{
		{"configured URL", "https://api.example.com/v1", "Bearer secret"},
		{"other host", "https://other.example.com/v1", ""},
		{"other port", "https://api.example.com:8080/v1", ""},
		{"scheme downgrade", "http://api.example.com/v1", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			transport := bearerTransport{host: "api.example.com", scheme: "https", key: "secret", next: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if got := r.Header.Get("Authorization"); got != tc.want {
					t.Errorf("authorization = %q; want %q", got, tc.want)
				}
				return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("")), Header: make(http.Header)}, nil
			})}
			req, err := http.NewRequest(http.MethodGet, tc.target, nil)
			if err != nil {
				t.Fatal(err)
			}
			_, err = transport.RoundTrip(req)
			if err != nil {
				t.Fatal(err)
			}
			if req.Header.Get("Authorization") != "" {
				t.Error("transport mutated the caller request")
			}
		})
	}
}

func newTestClient(t *testing.T, handler http.HandlerFunc, opts ...Option) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	opts = append([]Option{WithAPIKey("test-key"), WithBaseURL(server.URL)}, opts...)
	client, err := NewClient(opts...)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestCreateSandbox(t *testing.T) {
	t.Parallel()
	t.Run("every create field", func(t *testing.T) {
		t.Parallel()
		name, location, user, workdir := "api-dev", "eu", "dev", "/workspace"
		policy, profile, pause, deleteAfter := "open", "models", "15m", "30d"
		mode := AutoPauseMode("suspend")
		public, resume := true, false
		tags := []string{"dev"}
		env := map[string]string{"NODE_ENV": "test"}
		cpu := int64(4)
		want := map[string]any{
			"org": "acme", "project": "web", "template": "go-dev:stable",
			"name": name, "location": location, "user": user, "workdir": workdir,
			"egress_network_policy": policy, "egress_proxy_profile": profile,
			"public": public, "tags": []any{"dev"}, "env": map[string]any{"NODE_ENV": "test"},
			"auto_pause_after": pause, "auto_pause_mode": "suspend", "auto_resume": resume,
			"auto_delete_after": deleteAfter, "resources": map[string]any{"cpu": float64(cpu)},
		}
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost || r.URL.Path != "/v1/sandboxes" || r.URL.Query().Get("wait") != "true" {
				t.Errorf("unexpected create request: %s %s", r.Method, r.URL)
			}
			if r.Header.Get("Authorization") != "Bearer test-key" {
				t.Errorf("unexpected authorization: %s", r.Header.Get("Authorization"))
			}
			var got map[string]any
			if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("create body:\ngot  %#v\nwant %#v", got, want)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{"id":"`+sandboxID+`"}`)
		}, WithOrg("acme"), WithProject("web"))
		created, err := client.Sandbox().Create(context.Background(), CreateSandboxParams{
			Template: "go-dev:stable", Wait: true, Name: &name, Location: &location,
			Resources: &SandboxResources{Cpu: &cpu}, Env: &env, User: &user, Workdir: &workdir,
			Tags: &tags, Public: &public, EgressNetworkPolicy: &policy, EgressProxyProfile: &profile,
			AutoPauseAfter: &pause, AutoPauseMode: &mode, AutoResume: &resume, AutoDeleteAfter: &deleteAfter,
		})
		if err != nil {
			t.Fatal(err)
		}
		if created.Id != sandboxID {
			t.Errorf("id = %q", created.Id)
		}
	})
	t.Run("explicit scope wins", func(t *testing.T) {
		t.Parallel()
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			var body api.CreateSandboxRequestBody
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Error(err)
			}
			if body.Org != "other" || body.Project != "new" {
				t.Errorf("scope = %s/%s", body.Org, body.Project)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{}`)
		}, WithOrg("acme"), WithProject("web"))
		_, err := client.Sandbox().Create(context.Background(), CreateSandboxParams{Org: "other", Project: "new", Template: "go-dev"})
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestSandboxList(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("org") != "acme" || q.Get("project") != "" || q.Get("cursor") != "next" || q.Get("limit") != "2" {
			t.Errorf("query = %v", q)
		}
		if _, exists := q["project"]; !exists {
			t.Error("explicit org-wide project was omitted")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"items":[],"pagination":{"next_cursor":"later"}}`)
	}, WithOrg("acme"), WithProject("web"))
	empty, cursor, limit := "", "next", int32(2)
	params := &ListSandboxesParams{Project: &empty, Cursor: &cursor, Limit: &limit}
	page, err := client.Sandbox().List(context.Background(), params)
	if err != nil {
		t.Fatal(err)
	}
	if page == nil || params.Org != "" {
		t.Errorf("page = %#v; input was changed: %#v", page, params)
	}
}

func TestSandboxAddress(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, path, query, method string }{
		{"name", "/v1/sandboxes/api-dev", "org=acme&project=web", http.MethodGet},
		{"ID", "/v1/sandboxes/" + sandboxID, "", http.MethodGet},
		{"delete name", "/v1/sandboxes/api-dev", "org=acme&project=web", http.MethodDelete},
		{"delete ID", "/v1/sandboxes/" + sandboxID, "", http.MethodDelete},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tc.method || r.URL.Path != tc.path || r.URL.RawQuery != tc.query {
					t.Errorf("request = %s %s", r.Method, r.URL)
				}
				if r.Method == http.MethodGet {
					w.Header().Set("Content-Type", "application/json")
					_, _ = io.WriteString(w, `{"id":"`+sandboxID+`"}`)
				} else {
					w.WriteHeader(http.StatusNoContent)
				}
			}, WithOrg("acme"), WithProject("web"))
			id := strings.TrimPrefix(tc.path, "/v1/sandboxes/")
			if tc.method == http.MethodGet {
				_, err := client.Sandbox().Get(context.Background(), id, nil)
				if err != nil {
					t.Fatal(err)
				}
			} else if err := client.Sandbox().Delete(context.Background(), id, nil); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestSandboxValidation(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	})
	for _, call := range []func() error{
		func() error {
			_, err := client.Sandbox().Create(context.Background(), CreateSandboxParams{})
			return err
		},
		func() error { _, err := client.Sandbox().Get(context.Background(), "", nil); return err },
		func() error { _, err := client.Sandbox().Get(context.Background(), "api-dev", nil); return err },
		func() error { _, err := client.Sandbox().List(context.Background(), nil); return err },
		func() error { return client.Sandbox().Delete(context.Background(), "", nil) },
		func() error { return client.Sandbox().Delete(context.Background(), "api-dev", nil) },
	} {
		if err := call(); err == nil {
			t.Error("expected local validation error")
		}
	}
}

func TestSandboxError(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.Header().Set("Retry-After", "5")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = io.WriteString(w, `{"detail":"try later","status":429}`)
	}, WithOrg("acme"))
	_, err := client.Sandbox().List(context.Background(), nil)
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 429 || apiErr.Header.Get("Retry-After") != "5" || apiErr.Problem == nil || *apiErr.Problem.Detail != "try later" {
		t.Errorf("error = %v", err)
	}
}

func TestClientOptions(t *testing.T) {
	t.Run("environment and option precedence", func(t *testing.T) {
		t.Setenv("TEKTONA_API_KEY", "env-key")
		t.Setenv("TEKTONA_ORG", "env-org")
		t.Setenv("TEKTONA_PROJECT", "env-project")
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "Bearer env-key" || r.URL.Query().Get("org") != "option-org" || r.URL.Query().Get("project") != "env-project" {
				t.Errorf("headers/query = %v %v", r.Header, r.URL.Query())
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{}`)
		}))
		t.Cleanup(server.Close)
		t.Setenv("TEKTONA_API_URL", server.URL)
		client, err := NewClient(WithOrg("option-org"))
		if err != nil {
			t.Fatal(err)
		}
		_, err = client.Sandbox().List(context.Background(), nil)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("invalid construction", func(t *testing.T) {
		t.Setenv("TEKTONA_API_KEY", "env-key")
		for _, opt := range []Option{WithAPIKey(""), WithBaseURL("/api"), WithBaseURL("ftp://example.com"), WithHTTPClient(nil)} {
			if _, err := NewClient(opt); err == nil {
				t.Error("expected construction error")
			}
		}
	})
	t.Run("explicit key and URL override environment", func(t *testing.T) {
		t.Setenv("TEKTONA_API_KEY", "env-key")
		t.Setenv("TEKTONA_API_URL", "/invalid")
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "Bearer override" {
				t.Error("API key did not override environment")
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{}`)
		}))
		t.Cleanup(server.Close)
		client, err := NewClient(WithAPIKey("override"), WithBaseURL(server.URL), WithOrg("acme"))
		if err != nil {
			t.Fatal(err)
		}
		_, err = client.Sandbox().List(context.Background(), nil)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("explicit empty scope blocks environment", func(t *testing.T) {
		t.Setenv("TEKTONA_ORG", "env-org")
		t.Setenv("TEKTONA_API_URL", "")
		client, err := NewClient(WithAPIKey("key"), WithOrg(""))
		if err != nil {
			t.Fatal(err)
		}
		_, err = client.Sandbox().List(context.Background(), nil)
		if err == nil {
			t.Error("expected missing org error")
		}
	})
	t.Run("raw calls retain auth and no scope defaults", func(t *testing.T) {
		t.Parallel()
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "Bearer test-key" {
				t.Error("raw auth override")
			}
			if r.URL.Query().Get("org") != "" {
				t.Error("raw call inherited scope")
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{}`)
		}, WithOrg("acme"))
		_, err := client.Raw().ListSandboxesWithResponse(context.Background(), &api.ListSandboxesParams{}, func(_ context.Context, r *http.Request) error {
			r.Header.Set("Authorization", "Bearer wrong")
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("path prefix", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/v1/sandboxes" {
				t.Errorf("path = %s", r.URL.Path)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{}`)
		}))
		t.Cleanup(server.Close)
		client, err := NewClient(WithAPIKey("key"), WithBaseURL(server.URL+"/api"), WithOrg("acme"))
		if err != nil {
			t.Fatal(err)
		}
		_, err = client.Sandbox().List(context.Background(), nil)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestCancelSandboxRequest(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}, WithOrg("acme"))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := client.Sandbox().List(ctx, nil)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("error = %v", err)
	}
}
