package tektona

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/tektona-ai/tektona-go/api"
)

func TestTemplateBuildMethods(t *testing.T) {
	t.Parallel()
	cursor, name, key, tag := "next", "go", "retry-1", "stable"
	limit := int32(3)
	level := api.GetTemplateBuildLogsParamsLevelInfo
	direction := api.Forward
	for _, tc := range []struct {
		name, method, path, query, body, response string
		status                                    int
		call                                      func(TemplateBuild) (any, error)
	}{
		{"list org", "GET", "/v1/orgs/acme/template-builds", "project=web&template=go&cursor=next&limit=3", "", `{"items":[]}`, 200, func(b TemplateBuild) (any, error) {
			return b.ListForOrg(context.Background(), "", &ListTemplateBuildsForOrgParams{Cursor: &cursor, Limit: &limit, Project: strPtr("web"), Template: &name})
		}},
		{"create org", "POST", "/v1/orgs/other/template-builds", "", `{"template":"go","build":{"image":"ubuntu:24.04"},"sandbox":{},"tag":"stable"}`, `{"id":"build-1"}`, 202, func(b TemplateBuild) (any, error) {
			return b.CreateForOrg(context.Background(), CreateTemplateBuildParams{Org: "other", Project: "ignored", Template: name, Build: api.TemplateBuildSpec{Image: "ubuntu:24.04"}, Sandbox: &api.SandboxSpec{}, Tag: &tag, IdempotencyKey: &key})
		}},
		{"list project", "GET", "/v1/orgs/acme/projects/web/template-builds", "q=go&cursor=next&limit=3", "", `{"items":[]}`, 200, func(b TemplateBuild) (any, error) {
			return b.ListForProject(context.Background(), "", "", &ListTemplateBuildsForProjectParams{Cursor: &cursor, Limit: &limit, Q: &name})
		}},
		{"create project", "POST", "/v1/orgs/other/projects/site/template-builds", "", `{"template":"go","build":{"image":"ubuntu:24.04"},"tag":"stable"}`, `{"id":"build-1"}`, 202, func(b TemplateBuild) (any, error) {
			return b.CreateForProject(context.Background(), CreateTemplateBuildParams{Org: "other", Project: "site", Template: name, Build: api.TemplateBuildSpec{Image: "ubuntu:24.04"}, Tag: &tag, IdempotencyKey: &key})
		}},
		{"get", "GET", "/v1/template-builds/build-1", "", "", `{"id":"build-1"}`, 200, func(b TemplateBuild) (any, error) {
			return b.Get(context.Background(), "build-1")
		}},
		{"cancel", "PUT", "/v1/template-builds/build-1/cancel", "", "", `{"id":"build-1"}`, 202, func(b TemplateBuild) (any, error) {
			return b.Cancel(context.Background(), "build-1")
		}},
		{"logs", "GET", "/v1/template-builds/build-1/logs", "cursor=next&direction=forward&level=info&limit=3", "", `{"items":[]}`, 200, func(b TemplateBuild) (any, error) {
			return b.GetLogs(context.Background(), "build-1", &GetTemplateBuildLogsParams{Cursor: &cursor, Direction: &direction, Level: &level, Limit: &limit})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tc.method || r.URL.Path != tc.path || r.URL.RawQuery != tc.query {
					t.Errorf("request = %s %s; want %s %s?%s", r.Method, r.URL, tc.method, tc.path, tc.query)
				}
				if r.Header.Get("Authorization") != "Bearer test-key" {
					t.Errorf("authorization = %q", r.Header.Get("Authorization"))
				}
				if strings.HasPrefix(tc.name, "create") && r.Header.Get("Idempotency-Key") != key {
					t.Errorf("idempotency key = %q", r.Header.Get("Idempotency-Key"))
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
			got, err := tc.call(client.TemplateBuild())
			if err != nil || got == nil {
				t.Fatalf("result = %v, %v", got, err)
			}
		})
	}
}

func strPtr(value string) *string { return &value }

func TestTemplateBuildValidation(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	}, WithOrg(""), WithProject(""))
	for _, tc := range []struct {
		name string
		call func(TemplateBuild) error
	}{
		{"list org", func(b TemplateBuild) error { _, err := b.ListForOrg(context.Background(), "", nil); return err }},
		{"create org scope", func(b TemplateBuild) error {
			_, err := b.CreateForOrg(context.Background(), CreateTemplateBuildParams{Template: "go"})
			return err
		}},
		{"create org template", func(b TemplateBuild) error {
			_, err := b.CreateForOrg(context.Background(), CreateTemplateBuildParams{Org: "acme"})
			return err
		}},
		{"list project", func(b TemplateBuild) error {
			_, err := b.ListForProject(context.Background(), "acme", "", nil)
			return err
		}},
		{"create project scope", func(b TemplateBuild) error {
			_, err := b.CreateForProject(context.Background(), CreateTemplateBuildParams{Org: "acme", Template: "go"})
			return err
		}},
		{"create project template", func(b TemplateBuild) error {
			_, err := b.CreateForProject(context.Background(), CreateTemplateBuildParams{Org: "acme", Project: "web"})
			return err
		}},
		{"get", func(b TemplateBuild) error { _, err := b.Get(context.Background(), ""); return err }},
		{"cancel", func(b TemplateBuild) error { _, err := b.Cancel(context.Background(), ""); return err }},
		{"logs", func(b TemplateBuild) error { _, err := b.GetLogs(context.Background(), "", nil); return err }},
		{"stream", func(b TemplateBuild) error { _, err := b.StreamLogs(context.Background(), "", nil); return err }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if err := tc.call(client.TemplateBuild()); err == nil {
				t.Error("expected local validation error")
			}
		})
	}
}

func TestTemplateBuildErrors(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		call func(TemplateBuild) error
	}{
		{"list org", func(b TemplateBuild) error { _, err := b.ListForOrg(context.Background(), "", nil); return err }},
		{"create org", func(b TemplateBuild) error {
			_, err := b.CreateForOrg(context.Background(), CreateTemplateBuildParams{Template: "go"})
			return err
		}},
		{"list project", func(b TemplateBuild) error { _, err := b.ListForProject(context.Background(), "", "", nil); return err }},
		{"create project", func(b TemplateBuild) error {
			_, err := b.CreateForProject(context.Background(), CreateTemplateBuildParams{Template: "go"})
			return err
		}},
		{"get", func(b TemplateBuild) error { _, err := b.Get(context.Background(), "build-1"); return err }},
		{"cancel", func(b TemplateBuild) error { _, err := b.Cancel(context.Background(), "build-1"); return err }},
		{"logs", func(b TemplateBuild) error { _, err := b.GetLogs(context.Background(), "build-1", nil); return err }},
		{"stream", func(b TemplateBuild) error { _, err := b.StreamLogs(context.Background(), "build-1", nil); return err }},
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
			if err := tc.call(client.TemplateBuild()); !errors.As(err, &apiErr) || apiErr.StatusCode != 429 || apiErr.Header.Get("Retry-After") != "5" || apiErr.Problem == nil {
				t.Errorf("error = %v", err)
			}
		})
	}
}

func TestTemplateBuildStream(t *testing.T) {
	t.Parallel()
	for _, action := range []string{"close", "cancel"} {
		t.Run(action, func(t *testing.T) {
			t.Parallel()
			done := make(chan struct{})
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				defer close(done)
				if r.URL.Path != "/v1/template-builds/build-1/logs/stream" || r.URL.Query().Get("cursor") != "next" || r.Header.Get("Last-Event-ID") != "previous" {
					t.Errorf("stream request = %s, headers = %v", r.URL, r.Header)
				}
				w.Header().Set("Content-Type", "text/event-stream")
				_, _ = io.WriteString(w, "id: next\ndata: one\n\n")
				w.(http.Flusher).Flush()
				<-r.Context().Done()
			}, WithOrg("acme"), WithProject("web"))
			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)
			cursor, previous := "next", "previous"
			stream, err := client.TemplateBuild().StreamLogs(ctx, "build-1", &StreamTemplateBuildLogsParams{Cursor: &cursor, LastEventID: &previous})
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = stream.Close() })
			frame := make([]byte, len("id: next\ndata: one\n\n"))
			if _, err := io.ReadFull(stream, frame); err != nil || string(frame) != "id: next\ndata: one\n\n" {
				t.Fatalf("frame = %q, %v", frame, err)
			}
			if action == "close" {
				if err := stream.Close(); err != nil {
					t.Fatal(err)
				}
			} else {
				cancel()
				if _, err := stream.Read(frame); err == nil {
					t.Error("read continued after cancellation")
				}
			}
			select {
			case <-done:
			case <-time.After(2 * time.Second):
				t.Error("server did not observe stream closure")
			}
		})
	}
}
