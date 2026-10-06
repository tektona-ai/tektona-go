package tektona

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tektona-ai/tektona-go/api"
)

// StartDesktopParams holds optional scope for a sandbox name.
type StartDesktopParams = api.StartDesktopParams

// StopDesktopParams holds optional scope for a sandbox name.
type StopDesktopParams = api.StopDesktopParams

// RebootSandboxParams holds optional scope for a sandbox name.
type RebootSandboxParams = api.RebootSandboxParams

// ResetSandboxParams holds optional scope for a sandbox name.
type ResetSandboxParams = api.ResetSandboxParams

// ResumeSandboxParams holds optional scope for a sandbox name.
type ResumeSandboxParams = api.ResumeSandboxParams

// PauseSandboxParams combines sandbox scope with the pause body.
type PauseSandboxParams struct {
	Org, Project *string
	Mode         *api.PauseSandboxBodyMode
}

// ForkSandboxParams combines sandbox scope with the fork body.
type ForkSandboxParams struct {
	Org, Project *string
	Mode         *api.ForkSandboxBodyMode
	Name         *string
	Tags         *[]string
}

// ResizeSandboxParams combines sandbox scope with the resize body.
type ResizeSandboxParams struct {
	Org, Project *string
	Cpu          *int64
	Memory       *int64
	Disk         *int64
	Force        *bool
}

// RenameSandboxParams combines sandbox scope with the new name. Nil removes the name.
type RenameSandboxParams struct {
	Org, Project *string
	Name         *string
}

// TransferSandboxParams combines sandbox scope with the new owner and optional reason.
type TransferSandboxParams struct {
	Org, Project *string
	To           string
	Reason       *string
}

func (s *sandboxClient) actionAddress(idOrName string, org, project **string) error {
	if idOrName == "" {
		return fmt.Errorf("sandbox action requires an ID or name")
	}
	return s.address(idOrName, org, project)
}

func (s *sandboxClient) StartDesktop(ctx context.Context, idOrName string, params *StartDesktopParams) (*api.DesktopActionOutputBody, error) {
	var p StartDesktopParams
	if params != nil {
		p = *params
	}
	if err := s.actionAddress(idOrName, &p.Org, &p.Project); err != nil {
		return nil, err
	}
	return resourceJSON[api.DesktopActionOutputBody](func() (*http.Response, error) {
		return s.raw.StartDesktop(ctx, idOrName, &p)
	}, http.StatusOK, "start desktop")
}

func (s *sandboxClient) StopDesktop(ctx context.Context, idOrName string, params *StopDesktopParams) (*api.DesktopActionOutputBody, error) {
	var p StopDesktopParams
	if params != nil {
		p = *params
	}
	if err := s.actionAddress(idOrName, &p.Org, &p.Project); err != nil {
		return nil, err
	}
	return resourceJSON[api.DesktopActionOutputBody](func() (*http.Response, error) {
		return s.raw.StopDesktop(ctx, idOrName, &p)
	}, http.StatusOK, "stop desktop")
}

func (s *sandboxClient) Pause(ctx context.Context, idOrName string, params PauseSandboxParams) (*api.PauseSandboxResponse, error) {
	if err := s.actionAddress(idOrName, &params.Org, &params.Project); err != nil {
		return nil, err
	}
	return resourceJSON[api.PauseSandboxResponse](func() (*http.Response, error) {
		return s.raw.PauseSandbox(ctx, idOrName, &api.PauseSandboxParams{Org: params.Org, Project: params.Project}, api.PauseSandboxJSONRequestBody{Mode: params.Mode})
	}, http.StatusOK, "pause sandbox")
}

func (s *sandboxClient) Reboot(ctx context.Context, idOrName string, params *RebootSandboxParams) (*api.RebootSandboxResponseBody, error) {
	var p RebootSandboxParams
	if params != nil {
		p = *params
	}
	if err := s.actionAddress(idOrName, &p.Org, &p.Project); err != nil {
		return nil, err
	}
	return resourceJSON[api.RebootSandboxResponseBody](func() (*http.Response, error) {
		return s.raw.RebootSandbox(ctx, idOrName, &p)
	}, http.StatusOK, "reboot sandbox")
}

func (s *sandboxClient) Reset(ctx context.Context, idOrName string, params *ResetSandboxParams) (*api.RebootSandboxResponseBody, error) {
	var p ResetSandboxParams
	if params != nil {
		p = *params
	}
	if err := s.actionAddress(idOrName, &p.Org, &p.Project); err != nil {
		return nil, err
	}
	return resourceJSON[api.RebootSandboxResponseBody](func() (*http.Response, error) {
		return s.raw.ResetSandbox(ctx, idOrName, &p)
	}, http.StatusOK, "reset sandbox")
}

func (s *sandboxClient) Resize(ctx context.Context, idOrName string, params ResizeSandboxParams) (*api.ResizeResponse, error) {
	if err := s.actionAddress(idOrName, &params.Org, &params.Project); err != nil {
		return nil, err
	}
	return resourceJSON[api.ResizeResponse](func() (*http.Response, error) {
		return s.raw.ResizeSandbox(ctx, idOrName, &api.ResizeSandboxParams{Org: params.Org, Project: params.Project}, api.ResizeSandboxJSONRequestBody{
			Cpu: params.Cpu, Memory: params.Memory, Disk: params.Disk, Force: params.Force,
		})
	}, http.StatusOK, "resize sandbox")
}

func (s *sandboxClient) Resume(ctx context.Context, idOrName string, params *ResumeSandboxParams) (*api.ResumeSandboxResponseBody, error) {
	var p ResumeSandboxParams
	if params != nil {
		p = *params
	}
	if err := s.actionAddress(idOrName, &p.Org, &p.Project); err != nil {
		return nil, err
	}
	return resourceJSON[api.ResumeSandboxResponseBody](func() (*http.Response, error) {
		return s.raw.ResumeSandbox(ctx, idOrName, &p)
	}, http.StatusOK, "resume sandbox")
}

func (s *sandboxClient) Fork(ctx context.Context, idOrName string, params ForkSandboxParams) (*api.ForkResponse, error) {
	if err := s.actionAddress(idOrName, &params.Org, &params.Project); err != nil {
		return nil, err
	}
	return resourceJSON[api.ForkResponse](func() (*http.Response, error) {
		return s.raw.ForkSandbox(ctx, idOrName, &api.ForkSandboxParams{Org: params.Org, Project: params.Project}, api.ForkSandboxJSONRequestBody{
			Mode: params.Mode, Name: params.Name, Tags: params.Tags,
		})
	}, http.StatusOK, "fork sandbox")
}

func (s *sandboxClient) Rename(ctx context.Context, idOrName string, params RenameSandboxParams) (*api.RenameSandboxResponse, error) {
	if err := s.actionAddress(idOrName, &params.Org, &params.Project); err != nil {
		return nil, err
	}
	return resourceJSON[api.RenameSandboxResponse](func() (*http.Response, error) {
		return s.raw.RenameSandbox(ctx, idOrName, &api.RenameSandboxParams{Org: params.Org, Project: params.Project}, api.RenameSandboxJSONRequestBody{Name: params.Name})
	}, http.StatusOK, "rename sandbox")
}

func (s *sandboxClient) Transfer(ctx context.Context, idOrName string, params TransferSandboxParams) (*api.SandboxTransferResponse, error) {
	if err := s.actionAddress(idOrName, &params.Org, &params.Project); err != nil {
		return nil, err
	}
	if params.To == "" {
		return nil, fmt.Errorf("transfer sandbox requires to")
	}
	return resourceJSON[api.SandboxTransferResponse](func() (*http.Response, error) {
		return s.raw.TransferSandbox(ctx, idOrName, &api.TransferSandboxParams{Org: params.Org, Project: params.Project}, api.TransferSandboxJSONRequestBody{
			To: params.To, Reason: params.Reason,
		})
	}, http.StatusOK, "transfer sandbox")
}
