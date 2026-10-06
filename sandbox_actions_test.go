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

type sandboxActionCase struct {
	name, method, path string
	body               map[string]any
	call               func(Sandbox, string) (any, error)
}

func sandboxActionCases() []sandboxActionCase {
	mode := api.PauseSandboxBodyModeHibernate
	forkMode := api.Full
	name, reason := "other", "handoff"
	tags := []string{"test"}
	cpu, memory, disk := int64(2), int64(4), int64(20)
	force := true
	return []sandboxActionCase{
		{"start desktop", http.MethodPut, "desktop-start", nil, func(s Sandbox, id string) (any, error) { return s.StartDesktop(context.Background(), id, nil) }},
		{"stop desktop", http.MethodPut, "desktop-stop", nil, func(s Sandbox, id string) (any, error) { return s.StopDesktop(context.Background(), id, nil) }},
		{"pause", http.MethodPut, "pause", map[string]any{"mode": "hibernate"}, func(s Sandbox, id string) (any, error) {
			return s.Pause(context.Background(), id, PauseSandboxParams{Mode: &mode})
		}},
		{"reboot", http.MethodPut, "reboot", nil, func(s Sandbox, id string) (any, error) { return s.Reboot(context.Background(), id, nil) }},
		{"reset", http.MethodPut, "reset", nil, func(s Sandbox, id string) (any, error) { return s.Reset(context.Background(), id, nil) }},
		{"resize", http.MethodPut, "resize", map[string]any{"cpu": float64(cpu), "memory": float64(memory), "disk": float64(disk), "force": force}, func(s Sandbox, id string) (any, error) {
			return s.Resize(context.Background(), id, ResizeSandboxParams{Cpu: &cpu, Memory: &memory, Disk: &disk, Force: &force})
		}},
		{"resume", http.MethodPut, "resume", nil, func(s Sandbox, id string) (any, error) { return s.Resume(context.Background(), id, nil) }},
		{"fork", http.MethodPost, "fork", map[string]any{"mode": "full", "name": name, "tags": []any{"test"}}, func(s Sandbox, id string) (any, error) {
			return s.Fork(context.Background(), id, ForkSandboxParams{Mode: &forkMode, Name: &name, Tags: &tags})
		}},
		{"rename", http.MethodPut, "rename", map[string]any{"name": name}, func(s Sandbox, id string) (any, error) {
			return s.Rename(context.Background(), id, RenameSandboxParams{Name: &name})
		}},
		{"transfer", http.MethodPut, "transfer", map[string]any{"to": "owner@example.com", "reason": reason}, func(s Sandbox, id string) (any, error) {
			return s.Transfer(context.Background(), id, TransferSandboxParams{To: "owner@example.com", Reason: &reason})
		}},
	}
}

func TestSandboxActions(t *testing.T) {
	for _, tc := range sandboxActionCases() {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tc.method || r.URL.Path != "/v1/sandboxes/api-dev/"+tc.path || r.URL.RawQuery != "org=acme&project=web" {
					t.Errorf("request = %s %s", r.Method, r.URL)
				}
				if r.Header.Get("Authorization") != "Bearer test-key" {
					t.Errorf("authorization = %q", r.Header.Get("Authorization"))
				}
				if tc.body != nil {
					var got map[string]any
					if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
						t.Error(err)
					}
					if !reflect.DeepEqual(got, tc.body) {
						t.Errorf("body = %#v; want %#v", got, tc.body)
					}
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = io.WriteString(w, `{"id":"`+sandboxID+`","status":"active"}`)
			}, WithOrg("acme"), WithProject("web"))
			got, err := tc.call(client.Sandbox(), "api-dev")
			if err != nil || got == nil {
				t.Errorf("response = %#v, %v", got, err)
			}
		})
	}
}

func TestSandboxActionsIDSkipsScope(t *testing.T) {
	for _, tc := range sandboxActionCases() {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.RawQuery != "" || r.URL.Path != "/v1/sandboxes/"+sandboxID+"/"+tc.path {
					t.Errorf("request = %s", r.URL)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = io.WriteString(w, `{}`)
			}, WithOrg("acme"), WithProject("web"))
			if _, err := tc.call(client.Sandbox(), sandboxID); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestSandboxActionsValidation(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	})
	for _, tc := range sandboxActionCases() {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := tc.call(client.Sandbox(), ""); err == nil {
				t.Error("empty ID or name accepted")
			}
			if _, err := tc.call(client.Sandbox(), "api-dev"); err == nil {
				t.Error("name without scope accepted")
			}
		})
	}
	if _, err := client.Sandbox().Transfer(context.Background(), sandboxID, TransferSandboxParams{}); err == nil {
		t.Error("transfer without owner accepted")
	}
}

func TestSandboxActionsScopeOverrides(t *testing.T) {
	org, project := "other", "different"
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("org") != org || r.URL.Query().Get("project") != project {
			t.Errorf("scope = %v", r.URL.Query())
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{}`)
	}, WithOrg("acme"), WithProject("web"))
	ctx := context.Background()
	s := client.Sandbox()
	calls := []func() error{
		func() error {
			_, err := s.StartDesktop(ctx, "api-dev", &StartDesktopParams{Org: &org, Project: &project})
			return err
		},
		func() error {
			_, err := s.StopDesktop(ctx, "api-dev", &StopDesktopParams{Org: &org, Project: &project})
			return err
		},
		func() error {
			_, err := s.Pause(ctx, "api-dev", PauseSandboxParams{Org: &org, Project: &project})
			return err
		},
		func() error {
			_, err := s.Reboot(ctx, "api-dev", &RebootSandboxParams{Org: &org, Project: &project})
			return err
		},
		func() error {
			_, err := s.Reset(ctx, "api-dev", &ResetSandboxParams{Org: &org, Project: &project})
			return err
		},
		func() error {
			_, err := s.Resize(ctx, "api-dev", ResizeSandboxParams{Org: &org, Project: &project})
			return err
		},
		func() error {
			_, err := s.Resume(ctx, "api-dev", &ResumeSandboxParams{Org: &org, Project: &project})
			return err
		},
		func() error {
			_, err := s.Fork(ctx, "api-dev", ForkSandboxParams{Org: &org, Project: &project})
			return err
		},
		func() error {
			_, err := s.Rename(ctx, "api-dev", RenameSandboxParams{Org: &org, Project: &project})
			return err
		},
		func() error {
			_, err := s.Transfer(ctx, "api-dev", TransferSandboxParams{Org: &org, Project: &project, To: "owner"})
			return err
		},
	}
	for _, call := range calls {
		if err := call(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestSandboxActionParamsUnchanged(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{}`)
	}, WithOrg("acme"), WithProject("web"))
	query := &StartDesktopParams{}
	if _, err := client.Sandbox().StartDesktop(context.Background(), "api-dev", query); err != nil {
		t.Fatal(err)
	}
	if query.Org != nil || query.Project != nil {
		t.Errorf("query parameters changed: %#v", query)
	}
	flat := PauseSandboxParams{}
	if _, err := client.Sandbox().Pause(context.Background(), "api-dev", flat); err != nil {
		t.Fatal(err)
	}
	if flat.Org != nil || flat.Project != nil {
		t.Errorf("flat parameters changed: %#v", flat)
	}
}

func TestRenameSandboxRemovesName(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if value, ok := body["name"]; !ok || value != nil {
			t.Errorf("rename body = %#v; want name: null", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{}`)
	})
	if _, err := client.Sandbox().Rename(context.Background(), sandboxID, RenameSandboxParams{}); err != nil {
		t.Fatal(err)
	}
}

func TestSandboxActionEmptyScopeOverride(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	}, WithOrg("acme"), WithProject("web"))
	empty := ""
	if _, err := client.Sandbox().StartDesktop(context.Background(), "api-dev", &StartDesktopParams{Project: &empty}); err == nil {
		t.Error("empty project override accepted")
	}
	if _, err := client.Sandbox().Fork(context.Background(), "api-dev", ForkSandboxParams{Org: &empty}); err == nil {
		t.Error("empty org override accepted")
	}
}

func TestSandboxActionResponseStatus(t *testing.T) {
	for _, tc := range sandboxActionCases() {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/problem+json")
				w.Header().Set("Retry-After", "10")
				w.WriteHeader(http.StatusAccepted)
				_, _ = io.WriteString(w, `{"detail":"not ready"}`)
			})
			_, err := tc.call(client.Sandbox(), sandboxID)
			var apiErr *APIError
			if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusAccepted || apiErr.Header.Get("Retry-After") != "10" {
				t.Errorf("error = %v", err)
			}
		})
	}
}
