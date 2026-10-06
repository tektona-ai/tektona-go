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

func ExampleSandbox() {
	var sandbox tektona.Sandbox = fakeSandbox{}
	created, err := sandbox.Create(context.Background(), tektona.CreateSandboxParams{Template: "go-dev"})
	if err != nil {
		panic(err)
	}
	fmt.Println(created.Id)
	// Output: go-dev
}
