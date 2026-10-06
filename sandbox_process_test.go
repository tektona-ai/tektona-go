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

func TestSandboxProcessMethods(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	org, project := "other", "service"
	resume, force, enabled, eof := true, true, true, true
	tail := int32(3)
	autostart := api.ListSandboxProcessesParamsAutostart("false")
	mode := api.StartProcessBodyOnHibernate("preserve")
	maxLog, timeout, waitTimeout := int64(1234), int64(10), int64(20)
	name, cwd, user := "job", "/srv", "root"
	env := map[string]string{"A": "B"}
	tty := api.TtySizeDTO{Cols: 80, Rows: 24}
	start := StartSandboxProcessOptions{
		Org: &org, Project: &project, Command: "echo hi", Autostart: &enabled,
		PreventAutoPause: &enabled, Wait: &enabled, Cwd: &cwd, Name: &name, User: &user,
		Env: &env, MaxLogBytes: &maxLog, TimeoutSeconds: &timeout,
		WaitTimeoutSeconds: &waitTimeout, OnHibernate: &mode, Tty: &tty,
	}
	for _, tc := range []struct {
		name, method, path, query, body, response string
		status                                    int
		call                                      func(SandboxProcess) (any, error)
	}{
		{"list", "GET", "/v1/sandboxes/box/processes", "autostart=false&org=other&project=service&resume=true", "", `{"processes":[]}`, 200, func(p SandboxProcess) (any, error) {
			return p.List(ctx, "box", &ListSandboxProcessesParams{Org: &org, Project: &project, Autostart: &autostart, Resume: &resume})
		}},
		{"start", "POST", "/v1/sandboxes/box/processes", "org=other&project=service", `{"command":"echo hi","autostart":true,"prevent_auto_pause":true,"wait":true,"cwd":"/srv","name":"job","user":"root","env":{"A":"B"},"max_log_bytes":1234,"timeout_seconds":10,"wait_timeout_seconds":20,"on_hibernate":"preserve","tty":{"cols":80,"rows":24,"enabled":false}}`, `{"id":"p"}`, 201, func(p SandboxProcess) (any, error) {
			return p.Start(ctx, "box", start)
		}},
		{"stop", "DELETE", "/v1/sandboxes/box/processes/job", "force=true&org=other&project=service&resume=true", "", `{"id":"p"}`, 202, func(p SandboxProcess) (any, error) {
			return p.Stop(ctx, "box", "job", &StopSandboxProcessParams{Org: &org, Project: &project, Force: &force, Resume: &resume})
		}},
		{"get", "GET", "/v1/sandboxes/box/processes/job", "org=other&project=service&resume=true", "", `{"id":"p"}`, 200, func(p SandboxProcess) (any, error) {
			return p.Get(ctx, "box", "job", &GetSandboxProcessParams{Org: &org, Project: &project, Resume: &resume})
		}},
		{"autostart", "PUT", "/v1/sandboxes/box/processes/job/autostart", "org=other&project=service&resume=true", `{"enabled":true}`, "", 204, func(p SandboxProcess) (any, error) {
			return true, p.SetAutostart(ctx, "box", "job", SetSandboxProcessAutostartOptions{Org: &org, Project: &project, Resume: &resume, Enabled: enabled})
		}},
		{"logs", "GET", "/v1/sandboxes/box/processes/job/logs", "org=other&project=service&resume=true&tail=3", "", `{"frames":[]}`, 200, func(p SandboxProcess) (any, error) {
			return p.GetLogs(ctx, "box", "job", &GetSandboxProcessLogsParams{Org: &org, Project: &project, Resume: &resume, Tail: &tail})
		}},
		{"rename", "PUT", "/v1/sandboxes/box/processes/job/name", "org=other&project=service&resume=true", `{"name":"next"}`, `{"id":"p"}`, 200, func(p SandboxProcess) (any, error) {
			return p.Rename(ctx, "box", "job", RenameSandboxProcessOptions{Org: &org, Project: &project, Resume: &resume, Name: "next"})
		}},
		{"resize", "PUT", "/v1/sandboxes/box/processes/job/resize", "org=other&project=service&resume=true", `{"cols":80,"rows":24}`, "", 204, func(p SandboxProcess) (any, error) {
			return true, p.Resize(ctx, "box", "job", ResizeSandboxProcessOptions{Org: &org, Project: &project, Resume: &resume, Cols: 80, Rows: 24})
		}},
		{"signal", "PUT", "/v1/sandboxes/box/processes/job/signal", "org=other&project=service&resume=true", `{"signal":"TERM"}`, "", 204, func(p SandboxProcess) (any, error) {
			return true, p.Signal(ctx, "box", "job", SignalSandboxProcessOptions{Org: &org, Project: &project, Resume: &resume, Signal: "TERM"})
		}},
		{"stdin", "PUT", "/v1/sandboxes/box/processes/job/stdin", "org=other&project=service&resume=true", `{"data_base64":"aGk=","eof":true}`, "", 204, func(p SandboxProcess) (any, error) {
			return true, p.WriteStdin(ctx, "box", "job", WriteSandboxProcessStdinOptions{Org: &org, Project: &project, Resume: &resume, DataBase64: "aGk=", Eof: &eof})
		}},
		{"stream access", "PUT", "/v1/sandboxes/box/processes/job/stream-access", "org=other&project=service&resume=true", `{"mode":"logs"}`, `{"token":"abc","url":"wss://example.test","expires_at":"2026-10-06T12:00:00Z"}`, 200, func(p SandboxProcess) (any, error) {
			return p.CreateStreamAccess(ctx, "box", "job", CreateSandboxProcessStreamAccessOptions{Org: &org, Project: &project, Resume: &resume, Mode: "logs"})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				wantQuery, _ := url.ParseQuery(tc.query)
				if r.Method != tc.method || r.URL.Path != tc.path || !reflect.DeepEqual(r.URL.Query(), wantQuery) {
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
			got, err := tc.call(client.SandboxProcess())
			if err != nil || got == nil {
				t.Fatalf("result = %v, %v", got, err)
			}
		})
	}
}

func TestSandboxProcessValidation(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	})
	p := client.SandboxProcess()
	for _, tc := range []struct {
		name string
		call func(string, string) error
	}{
		{"list", func(s, _ string) error { _, err := p.List(ctx, s, nil); return err }},
		{"start", func(s, _ string) error { _, err := p.Start(ctx, s, StartSandboxProcessOptions{}); return err }},
		{"stop", func(s, r string) error { _, err := p.Stop(ctx, s, r, nil); return err }},
		{"get", func(s, r string) error { _, err := p.Get(ctx, s, r, nil); return err }},
		{"autostart", func(s, r string) error { return p.SetAutostart(ctx, s, r, SetSandboxProcessAutostartOptions{}) }},
		{"logs", func(s, r string) error { _, err := p.GetLogs(ctx, s, r, nil); return err }},
		{"rename", func(s, r string) error { _, err := p.Rename(ctx, s, r, RenameSandboxProcessOptions{}); return err }},
		{"resize", func(s, r string) error { return p.Resize(ctx, s, r, ResizeSandboxProcessOptions{}) }},
		{"signal", func(s, r string) error { return p.Signal(ctx, s, r, SignalSandboxProcessOptions{}) }},
		{"stdin", func(s, r string) error { return p.WriteStdin(ctx, s, r, WriteSandboxProcessStdinOptions{}) }},
		{"stream access", func(s, r string) error {
			_, err := p.CreateStreamAccess(ctx, s, r, CreateSandboxProcessStreamAccessOptions{})
			return err
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			for _, target := range [][2]string{{"", "job"}, {"box", ""}, {"box", "job"}} {
				if tc.name == "list" || tc.name == "start" {
					if target[0] != "" && target[1] == "" {
						continue
					}
				}
				if err := tc.call(target[0], target[1]); err == nil {
					t.Errorf("expected local error for %q/%q", target[0], target[1])
				}
			}
		})
	}
}

func TestSandboxProcessIDScopeAndCopy(t *testing.T) {
	t.Parallel()
	const id = "01ARZ3NDEKTSV4RRFFQ69G5FAV"
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Has("org") || r.URL.Query().Has("project") {
			t.Errorf("unexpected scope: %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"processes":[]}`)
	}, WithOrg("acme"), WithProject("web"))
	params := &ListSandboxProcessesParams{}
	if _, err := client.SandboxProcess().List(context.Background(), id, params); err != nil {
		t.Fatal(err)
	}
	if params.Org != nil || params.Project != nil {
		t.Error("modified caller parameters")
	}
}

func TestSandboxProcessNameScopeAndCopy(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("org") != "acme" || r.URL.Query().Get("project") != "web" {
			t.Errorf("scope = %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"processes":[]}`)
	}, WithOrg("acme"), WithProject("web"))
	params := &ListSandboxProcessesParams{}
	if _, err := client.SandboxProcess().List(context.Background(), "box", params); err != nil {
		t.Fatal(err)
	}
	if params.Org != nil || params.Project != nil {
		t.Error("modified caller parameters")
	}
}

func TestSandboxProcessExplicitEmptyScope(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	}, WithOrg("acme"), WithProject("web"))
	empty := ""
	if _, err := client.SandboxProcess().List(context.Background(), "box", &ListSandboxProcessesParams{Org: &empty}); err == nil {
		t.Error("explicit empty organization inherited default")
	}
	if _, err := client.SandboxProcess().Start(context.Background(), "box", StartSandboxProcessOptions{Project: &empty}); err == nil {
		t.Error("explicit empty project inherited default")
	}
}

func TestSandboxProcessResponseStatus(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		status int
		call   func(SandboxProcess) error
	}{
		{"start requires 201", http.StatusOK, func(p SandboxProcess) error {
			_, err := p.Start(context.Background(), "box", StartSandboxProcessOptions{Command: "date"})
			return err
		}},
		{"stop requires 202", http.StatusOK, func(p SandboxProcess) error {
			_, err := p.Stop(context.Background(), "box", "job", nil)
			return err
		}},
		{"resize requires 204", http.StatusOK, func(p SandboxProcess) error {
			return p.Resize(context.Background(), "box", "job", ResizeSandboxProcessOptions{Rows: 24, Cols: 80})
		}},
		{"stream access requires 200", http.StatusCreated, func(p SandboxProcess) error {
			_, err := p.CreateStreamAccess(context.Background(), "box", "job", CreateSandboxProcessStreamAccessOptions{Mode: "logs"})
			return err
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				_, _ = io.WriteString(w, `{}`)
			}, WithOrg("acme"), WithProject("web"))
			var apiErr *APIError
			if err := tc.call(client.SandboxProcess()); !errors.As(err, &apiErr) || apiErr.StatusCode != tc.status {
				t.Errorf("error = %v; want status %d", err, tc.status)
			}
		})
	}
}

func TestSandboxProcessError(t *testing.T) {
	t.Parallel()
	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.Header().Set("Retry-After", "4")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = io.WriteString(w, `{"detail":"later"}`)
	}, WithOrg("acme"), WithProject("web"))
	_, err := client.SandboxProcess().CreateStreamAccess(context.Background(), "box", "job", CreateSandboxProcessStreamAccessOptions{Mode: "logs"})
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 429 || apiErr.Header.Get("Retry-After") != "4" || apiErr.Problem == nil {
		t.Errorf("error = %v", err)
	}
}
