package tektona_test

import (
	"context"
	"fmt"

	tektona "github.com/tektona-ai/tektona-go"
	"github.com/tektona-ai/tektona-go/api"
)

type fakeSandbox struct{}

func (fakeSandbox) Create(_ context.Context, params tektona.CreateSandboxParams) (*api.CreateSandboxResponseBody, error) {
	return &api.CreateSandboxResponseBody{Id: params.Template}, nil
}
func (fakeSandbox) Get(context.Context, string, *tektona.GetSandboxParams) (*api.SandboxObject, error) {
	return nil, nil
}
func (fakeSandbox) List(context.Context, *tektona.ListSandboxesParams) (*api.ListSandboxesBody, error) {
	return nil, nil
}
func (fakeSandbox) Delete(context.Context, string, *tektona.DeleteSandboxParams) error { return nil }
func (fakeSandbox) StartDesktop(context.Context, string, *tektona.StartDesktopParams) (*api.DesktopActionOutputBody, error) {
	return nil, nil
}
func (fakeSandbox) StopDesktop(context.Context, string, *tektona.StopDesktopParams) (*api.DesktopActionOutputBody, error) {
	return nil, nil
}
func (fakeSandbox) Pause(context.Context, string, tektona.PauseSandboxParams) (*api.PauseSandboxResponse, error) {
	return nil, nil
}
func (fakeSandbox) Reboot(context.Context, string, *tektona.RebootSandboxParams) (*api.RebootSandboxResponseBody, error) {
	return nil, nil
}
func (fakeSandbox) Reset(context.Context, string, *tektona.ResetSandboxParams) (*api.RebootSandboxResponseBody, error) {
	return nil, nil
}
func (fakeSandbox) Resize(context.Context, string, tektona.ResizeSandboxParams) (*api.ResizeResponse, error) {
	return nil, nil
}
func (fakeSandbox) Resume(context.Context, string, *tektona.ResumeSandboxParams) (*api.ResumeSandboxResponseBody, error) {
	return nil, nil
}
func (fakeSandbox) Fork(context.Context, string, tektona.ForkSandboxParams) (*api.ForkResponse, error) {
	return nil, nil
}
func (fakeSandbox) Rename(context.Context, string, tektona.RenameSandboxParams) (*api.RenameSandboxResponse, error) {
	return nil, nil
}
func (fakeSandbox) Transfer(context.Context, string, tektona.TransferSandboxParams) (*api.SandboxTransferResponse, error) {
	return nil, nil
}

func ExampleSandbox() {
	var sandbox tektona.Sandbox = fakeSandbox{}
	created, err := sandbox.Create(context.Background(), tektona.CreateSandboxParams{Template: "go-dev"})
	if err != nil {
		panic(err)
	}
	fmt.Println(created.Id)
	// Output: go-dev
}
