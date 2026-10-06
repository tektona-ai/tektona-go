package tektona_test

import (
	"context"

	tektona "github.com/tektona-ai/tektona-go"
	"github.com/tektona-ai/tektona-go/api"
)

type fakeSandboxProcess struct{}

func (fakeSandboxProcess) List(context.Context, string, *tektona.ListSandboxProcessesParams) (*api.ProcessListResponse, error) {
	return &api.ProcessListResponse{}, nil
}
func (fakeSandboxProcess) Start(context.Context, string, tektona.StartSandboxProcessOptions) (*api.ProcessDTO, error) {
	return &api.ProcessDTO{}, nil
}
func (fakeSandboxProcess) Stop(context.Context, string, string, *tektona.StopSandboxProcessParams) (*api.ProcessDTO, error) {
	return &api.ProcessDTO{}, nil
}
func (fakeSandboxProcess) Get(context.Context, string, string, *tektona.GetSandboxProcessParams) (*api.ProcessDTO, error) {
	return &api.ProcessDTO{}, nil
}
func (fakeSandboxProcess) SetAutostart(context.Context, string, string, tektona.SetSandboxProcessAutostartOptions) error {
	return nil
}
func (fakeSandboxProcess) GetLogs(context.Context, string, string, *tektona.GetSandboxProcessLogsParams) (*api.LogsResponse, error) {
	return &api.LogsResponse{}, nil
}
func (fakeSandboxProcess) Rename(context.Context, string, string, tektona.RenameSandboxProcessOptions) (*api.ProcessDTO, error) {
	return &api.ProcessDTO{}, nil
}
func (fakeSandboxProcess) Resize(context.Context, string, string, tektona.ResizeSandboxProcessOptions) error {
	return nil
}
func (fakeSandboxProcess) Signal(context.Context, string, string, tektona.SignalSandboxProcessOptions) error {
	return nil
}
func (fakeSandboxProcess) WriteStdin(context.Context, string, string, tektona.WriteSandboxProcessStdinOptions) error {
	return nil
}
func (fakeSandboxProcess) CreateStreamAccess(context.Context, string, string, tektona.CreateSandboxProcessStreamAccessOptions) (*api.StreamAccessResponseDTO, error) {
	return &api.StreamAccessResponseDTO{}, nil
}

var _ tektona.SandboxProcess = fakeSandboxProcess{}
