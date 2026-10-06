package tektona

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tektona-ai/tektona-go/api"
)

// ListSandboxProcessesParams holds scope and process filters.
type ListSandboxProcessesParams = api.ListSandboxProcessesParams

// StopSandboxProcessParams holds scope and stop options.
type StopSandboxProcessParams = api.StopSandboxProcessParams

// GetSandboxProcessParams holds scope and the resume option.
type GetSandboxProcessParams = api.GetSandboxProcessParams

// GetSandboxProcessLogsParams holds scope and log options.
type GetSandboxProcessLogsParams = api.GetSandboxProcessLogsParams

// StartSandboxProcessOptions holds scope and process creation fields.
type StartSandboxProcessOptions struct {
	Org, Project                                    *string
	Command                                         string
	Autostart, PreventAutoPause, Wait               *bool
	Cwd, Name, User                                 *string
	Env                                             *map[string]string
	MaxLogBytes, TimeoutSeconds, WaitTimeoutSeconds *int64
	OnHibernate                                     *api.StartProcessBodyOnHibernate
	Tty                                             *api.TtySizeDTO
}

// SetSandboxProcessAutostartOptions holds scope and autostart options.
type SetSandboxProcessAutostartOptions struct {
	Org, Project *string
	Resume       *bool
	Enabled      bool
}

// RenameSandboxProcessOptions holds scope and the new name.
type RenameSandboxProcessOptions struct {
	Org, Project *string
	Resume       *bool
	Name         string
}

// ResizeSandboxProcessOptions holds scope and the terminal size.
type ResizeSandboxProcessOptions struct {
	Org, Project *string
	Resume       *bool
	Cols, Rows   int32
}

// SignalSandboxProcessOptions holds scope and the signal.
type SignalSandboxProcessOptions struct {
	Org, Project *string
	Resume       *bool
	Signal       api.SignalBodySignal
}

// WriteSandboxProcessStdinOptions holds scope and stdin data.
type WriteSandboxProcessStdinOptions struct {
	Org, Project *string
	Resume       *bool
	DataBase64   string
	Eof          *bool
}

// CreateSandboxProcessStreamAccessOptions holds scope and token mode.
type CreateSandboxProcessStreamAccessOptions struct {
	Org, Project *string
	Resume       *bool
	Mode         api.StreamAccessBodyMode
}

// SandboxProcess supports process operations within a sandbox.
type SandboxProcess interface {
	List(context.Context, string, *ListSandboxProcessesParams) (*api.ProcessListResponse, error)
	Start(context.Context, string, StartSandboxProcessOptions) (*api.ProcessDTO, error)
	Stop(context.Context, string, string, *StopSandboxProcessParams) (*api.ProcessDTO, error)
	Get(context.Context, string, string, *GetSandboxProcessParams) (*api.ProcessDTO, error)
	SetAutostart(context.Context, string, string, SetSandboxProcessAutostartOptions) error
	GetLogs(context.Context, string, string, *GetSandboxProcessLogsParams) (*api.LogsResponse, error)
	Rename(context.Context, string, string, RenameSandboxProcessOptions) (*api.ProcessDTO, error)
	Resize(context.Context, string, string, ResizeSandboxProcessOptions) error
	Signal(context.Context, string, string, SignalSandboxProcessOptions) error
	WriteStdin(context.Context, string, string, WriteSandboxProcessStdinOptions) error
	CreateStreamAccess(context.Context, string, string, CreateSandboxProcessStreamAccessOptions) (*api.StreamAccessResponseDTO, error)
}

type sandboxProcessClient struct {
	raw          *api.ClientWithResponses
	org, project string
}

// SandboxProcess returns the process methods for a sandbox.
func (c *Client) SandboxProcess() SandboxProcess {
	return &sandboxProcessClient{raw: c.raw, org: c.org, project: c.project}
}

func (s *sandboxProcessClient) address(sandbox, ref string, org, project **string) error {
	if sandbox == "" {
		return fmt.Errorf("process requires a sandbox ID or name")
	}
	if ref == "" {
		return fmt.Errorf("process requires a process ID or name")
	}
	return (&sandboxClient{org: s.org, project: s.project}).address(sandbox, org, project)
}

func (s *sandboxProcessClient) List(ctx context.Context, sandbox string, params *ListSandboxProcessesParams) (*api.ProcessListResponse, error) {
	var p ListSandboxProcessesParams
	if params != nil {
		p = *params
	}
	if err := s.address(sandbox, "list", &p.Org, &p.Project); err != nil {
		return nil, err
	}
	return resourceJSON[api.ProcessListResponse](func() (*http.Response, error) {
		return s.raw.ListSandboxProcesses(ctx, sandbox, &p)
	}, http.StatusOK, "list sandbox processes")
}

func (s *sandboxProcessClient) Start(ctx context.Context, sandbox string, opts StartSandboxProcessOptions) (*api.ProcessDTO, error) {
	if err := s.address(sandbox, "start", &opts.Org, &opts.Project); err != nil {
		return nil, err
	}
	return resourceJSON[api.ProcessDTO](func() (*http.Response, error) {
		return s.raw.StartSandboxProcess(ctx, sandbox, &api.StartSandboxProcessParams{Org: opts.Org, Project: opts.Project}, api.StartProcessBody{
			Command: opts.Command, Autostart: opts.Autostart, PreventAutoPause: opts.PreventAutoPause,
			Wait: opts.Wait, Cwd: opts.Cwd, Name: opts.Name, User: opts.User, Env: opts.Env,
			MaxLogBytes: opts.MaxLogBytes, TimeoutSeconds: opts.TimeoutSeconds,
			WaitTimeoutSeconds: opts.WaitTimeoutSeconds, OnHibernate: opts.OnHibernate, Tty: opts.Tty,
		})
	}, http.StatusCreated, "start sandbox process")
}

func (s *sandboxProcessClient) Stop(ctx context.Context, sandbox, ref string, params *StopSandboxProcessParams) (*api.ProcessDTO, error) {
	var p StopSandboxProcessParams
	if params != nil {
		p = *params
	}
	if err := s.address(sandbox, ref, &p.Org, &p.Project); err != nil {
		return nil, err
	}
	return resourceJSON[api.ProcessDTO](func() (*http.Response, error) {
		return s.raw.StopSandboxProcess(ctx, sandbox, ref, &p)
	}, http.StatusAccepted, "stop sandbox process")
}

func (s *sandboxProcessClient) Get(ctx context.Context, sandbox, ref string, params *GetSandboxProcessParams) (*api.ProcessDTO, error) {
	var p GetSandboxProcessParams
	if params != nil {
		p = *params
	}
	if err := s.address(sandbox, ref, &p.Org, &p.Project); err != nil {
		return nil, err
	}
	return resourceJSON[api.ProcessDTO](func() (*http.Response, error) {
		return s.raw.GetSandboxProcess(ctx, sandbox, ref, &p)
	}, http.StatusOK, "get sandbox process")
}

func (s *sandboxProcessClient) SetAutostart(ctx context.Context, sandbox, ref string, opts SetSandboxProcessAutostartOptions) error {
	if err := s.address(sandbox, ref, &opts.Org, &opts.Project); err != nil {
		return err
	}
	return resourceNoContent(func() (*http.Response, error) {
		return s.raw.SetSandboxProcessAutostart(ctx, sandbox, ref, &api.SetSandboxProcessAutostartParams{Org: opts.Org, Project: opts.Project, Resume: opts.Resume}, api.AutostartBody{Enabled: opts.Enabled})
	})
}

func (s *sandboxProcessClient) GetLogs(ctx context.Context, sandbox, ref string, params *GetSandboxProcessLogsParams) (*api.LogsResponse, error) {
	var p GetSandboxProcessLogsParams
	if params != nil {
		p = *params
	}
	if err := s.address(sandbox, ref, &p.Org, &p.Project); err != nil {
		return nil, err
	}
	return resourceJSON[api.LogsResponse](func() (*http.Response, error) {
		return s.raw.GetSandboxProcessLogs(ctx, sandbox, ref, &p)
	}, http.StatusOK, "get sandbox process logs")
}

func (s *sandboxProcessClient) Rename(ctx context.Context, sandbox, ref string, opts RenameSandboxProcessOptions) (*api.ProcessDTO, error) {
	if err := s.address(sandbox, ref, &opts.Org, &opts.Project); err != nil {
		return nil, err
	}
	return resourceJSON[api.ProcessDTO](func() (*http.Response, error) {
		return s.raw.RenameSandboxProcess(ctx, sandbox, ref, &api.RenameSandboxProcessParams{Org: opts.Org, Project: opts.Project, Resume: opts.Resume}, api.RenameProcessBody{Name: opts.Name})
	}, http.StatusOK, "rename sandbox process")
}

func (s *sandboxProcessClient) Resize(ctx context.Context, sandbox, ref string, opts ResizeSandboxProcessOptions) error {
	if err := s.address(sandbox, ref, &opts.Org, &opts.Project); err != nil {
		return err
	}
	return resourceNoContent(func() (*http.Response, error) {
		return s.raw.ResizeSandboxProcess(ctx, sandbox, ref, &api.ResizeSandboxProcessParams{Org: opts.Org, Project: opts.Project, Resume: opts.Resume}, api.ProcessResizeBody{Cols: opts.Cols, Rows: opts.Rows})
	})
}

func (s *sandboxProcessClient) Signal(ctx context.Context, sandbox, ref string, opts SignalSandboxProcessOptions) error {
	if err := s.address(sandbox, ref, &opts.Org, &opts.Project); err != nil {
		return err
	}
	return resourceNoContent(func() (*http.Response, error) {
		return s.raw.SignalSandboxProcess(ctx, sandbox, ref, &api.SignalSandboxProcessParams{Org: opts.Org, Project: opts.Project, Resume: opts.Resume}, api.SignalBody{Signal: opts.Signal})
	})
}

func (s *sandboxProcessClient) WriteStdin(ctx context.Context, sandbox, ref string, opts WriteSandboxProcessStdinOptions) error {
	if err := s.address(sandbox, ref, &opts.Org, &opts.Project); err != nil {
		return err
	}
	return resourceNoContent(func() (*http.Response, error) {
		return s.raw.WriteSandboxProcessStdin(ctx, sandbox, ref, &api.WriteSandboxProcessStdinParams{Org: opts.Org, Project: opts.Project, Resume: opts.Resume}, api.StdinBody{DataBase64: opts.DataBase64, Eof: opts.Eof})
	})
}

func (s *sandboxProcessClient) CreateStreamAccess(ctx context.Context, sandbox, ref string, opts CreateSandboxProcessStreamAccessOptions) (*api.StreamAccessResponseDTO, error) {
	if err := s.address(sandbox, ref, &opts.Org, &opts.Project); err != nil {
		return nil, err
	}
	return resourceJSON[api.StreamAccessResponseDTO](func() (*http.Response, error) {
		return s.raw.CreateSandboxProcessStreamAccess(ctx, sandbox, ref, &api.CreateSandboxProcessStreamAccessParams{Org: opts.Org, Project: opts.Project, Resume: opts.Resume}, api.StreamAccessBody{Mode: opts.Mode})
	}, http.StatusOK, "create sandbox process stream access")
}
