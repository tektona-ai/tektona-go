package tektona_test

import (
	"context"

	tektona "github.com/tektona-ai/tektona-go"
	"github.com/tektona-ai/tektona-go/api"
)

type fakeGitCredential struct{}

func (fakeGitCredential) List(context.Context, *tektona.ListGitCredentialsParams) (*api.GitCredentialListBody, error) {
	return &api.GitCredentialListBody{}, nil
}
func (fakeGitCredential) Create(context.Context, tektona.CreateGitCredentialParams) (*api.GitCredentialCreateBody, error) {
	return &api.GitCredentialCreateBody{}, nil
}
func (fakeGitCredential) Update(context.Context, string, tektona.UpdateGitCredentialParams) (*api.GitCredentialUpdateBody, error) {
	return &api.GitCredentialUpdateBody{}, nil
}
func (fakeGitCredential) Delete(context.Context, string, tektona.DeleteGitCredentialParams) error {
	return nil
}

type fakeRegistry struct{}

func (fakeRegistry) List(context.Context, *tektona.ListRegistriesParams) (*api.ListRegistriesBody, error) {
	return &api.ListRegistriesBody{}, nil
}
func (fakeRegistry) Create(context.Context, tektona.CreateRegistryParams) (*api.RegistryResponse, error) {
	return &api.RegistryResponse{}, nil
}
func (fakeRegistry) Get(context.Context, string, *tektona.GetRegistryParams) (*api.RegistryResponse, error) {
	return &api.RegistryResponse{}, nil
}
func (fakeRegistry) Update(context.Context, string, tektona.UpdateRegistryParams) (*api.RegistryResponse, error) {
	return &api.RegistryResponse{}, nil
}
func (fakeRegistry) Delete(context.Context, string, *tektona.DeleteRegistryParams) error {
	return nil
}

type fakeRepository struct{}

func (fakeRepository) List(context.Context, *tektona.ListRepositoriesParams) (*api.ListBody, error) {
	return &api.ListBody{}, nil
}
func (fakeRepository) Create(context.Context, tektona.CreateRepositoryParams) (*api.Response, error) {
	return &api.Response{}, nil
}
func (fakeRepository) Update(context.Context, string, tektona.UpdateRepositoryParams) (*api.Response, error) {
	return &api.Response{}, nil
}
func (fakeRepository) Delete(context.Context, string, *tektona.DeleteRepositoryParams) error {
	return nil
}
func (fakeRepository) ListBranches(context.Context, string, *tektona.ListRepositoryBranchesParams) (*api.BranchesBody, error) {
	return &api.BranchesBody{}, nil
}

var (
	_ tektona.GitCredential = fakeGitCredential{}
	_ tektona.Registry      = fakeRegistry{}
	_ tektona.Repository    = fakeRepository{}
)
