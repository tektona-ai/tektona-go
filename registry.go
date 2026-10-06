package tektona

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tektona-ai/tektona-go/api"
)

// ListRegistriesParams holds scope and pagination options.
type ListRegistriesParams struct {
	Org, Project string
	Cursor       *string
	Limit        *int32
}

// CreateRegistryParams holds scope, query, and registry fields.
type CreateRegistryParams struct {
	Org, Project string
	DryRun       *bool
	AuthType     api.CreateProjectRegistryInputBodyAuthType
	DisplayName  string
	Endpoint     string
	Name         string
	Namespace    *string
	Password     *string
	Token        *string
	Username     *string
}

// GetRegistryParams holds an optional scope override.
type GetRegistryParams struct{ Org, Project string }

// UpdateRegistryParams holds scope, query, and registry fields.
type UpdateRegistryParams struct {
	Org, Project string
	DryRun       *bool
	AuthType     api.UpdateProjectRegistryInputBodyAuthType
	DisplayName  string
	Endpoint     string
	Name         string
	Namespace    *string
	Password     *string
	Token        *string
	Username     *string
}

// DeleteRegistryParams holds an optional scope override.
type DeleteRegistryParams struct{ Org, Project string }

// Registry supports project container registry operations.
type Registry interface {
	List(ctx context.Context, params *ListRegistriesParams) (*api.ListRegistriesBody, error)
	Create(ctx context.Context, params CreateRegistryParams) (*api.RegistryResponse, error)
	Get(ctx context.Context, name string, params *GetRegistryParams) (*api.RegistryResponse, error)
	Update(ctx context.Context, name string, params UpdateRegistryParams) (*api.RegistryResponse, error)
	Delete(ctx context.Context, name string, params *DeleteRegistryParams) error
}

type registryClient struct {
	raw          *api.ClientWithResponses
	org, project string
}

// Registry returns the container registry methods.
func (c *Client) Registry() Registry {
	return registryClient{raw: c.raw, org: c.org, project: c.project}
}

func (r registryClient) address(org, project string) (string, string, error) {
	if org == "" {
		org = r.org
	}
	if project == "" {
		project = r.project
	}
	if org == "" || project == "" {
		return "", "", fmt.Errorf("registry requires org and project")
	}
	return org, project, nil
}

func (r registryClient) List(ctx context.Context, params *ListRegistriesParams) (*api.ListRegistriesBody, error) {
	var p ListRegistriesParams
	if params != nil {
		p = *params
	}
	org, project, err := r.address(p.Org, p.Project)
	if err != nil {
		return nil, err
	}
	return resourceJSON[api.ListRegistriesBody](func() (*http.Response, error) {
		return r.raw.ListProjectRegistries(ctx, org, project, &api.ListProjectRegistriesParams{Cursor: p.Cursor, Limit: p.Limit})
	}, http.StatusOK, "list registries")
}

func (r registryClient) Create(ctx context.Context, params CreateRegistryParams) (*api.RegistryResponse, error) {
	org, project, err := r.address(params.Org, params.Project)
	if err != nil {
		return nil, err
	}
	if params.AuthType == "" || params.DisplayName == "" || params.Endpoint == "" || params.Name == "" {
		return nil, fmt.Errorf("create registry requires auth type, display name, endpoint, and name")
	}
	return resourceJSON[api.RegistryResponse](func() (*http.Response, error) {
		return r.raw.CreateProjectRegistry(ctx, org, project, &api.CreateProjectRegistryParams{DryRun: params.DryRun},
			api.CreateProjectRegistryJSONRequestBody{
				AuthType: params.AuthType, DisplayName: params.DisplayName, Endpoint: params.Endpoint,
				Name: params.Name, Namespace: params.Namespace, Password: params.Password,
				Token: params.Token, Username: params.Username,
			})
	}, http.StatusCreated, "create registry")
}

func (r registryClient) Get(ctx context.Context, name string, params *GetRegistryParams) (*api.RegistryResponse, error) {
	var p GetRegistryParams
	if params != nil {
		p = *params
	}
	org, project, err := r.address(p.Org, p.Project)
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, fmt.Errorf("get registry requires name")
	}
	return resourceJSON[api.RegistryResponse](func() (*http.Response, error) {
		return r.raw.GetProjectRegistry(ctx, org, project, name)
	}, http.StatusOK, "get registry")
}

func (r registryClient) Update(ctx context.Context, name string, params UpdateRegistryParams) (*api.RegistryResponse, error) {
	org, project, err := r.address(params.Org, params.Project)
	if err != nil {
		return nil, err
	}
	if name == "" || params.AuthType == "" || params.DisplayName == "" || params.Endpoint == "" || params.Name == "" {
		return nil, fmt.Errorf("update registry requires name, auth type, display name, and endpoint")
	}
	return resourceJSON[api.RegistryResponse](func() (*http.Response, error) {
		return r.raw.UpdateProjectRegistry(ctx, org, project, name, &api.UpdateProjectRegistryParams{DryRun: params.DryRun},
			api.UpdateProjectRegistryJSONRequestBody{
				AuthType: params.AuthType, DisplayName: params.DisplayName, Endpoint: params.Endpoint,
				Name: params.Name, Namespace: params.Namespace, Password: params.Password,
				Token: params.Token, Username: params.Username,
			})
	}, http.StatusOK, "update registry")
}

func (r registryClient) Delete(ctx context.Context, name string, params *DeleteRegistryParams) error {
	var p DeleteRegistryParams
	if params != nil {
		p = *params
	}
	org, project, err := r.address(p.Org, p.Project)
	if err != nil {
		return err
	}
	if name == "" {
		return fmt.Errorf("delete registry requires name")
	}
	return resourceNoContent(func() (*http.Response, error) {
		return r.raw.DeleteProjectRegistry(ctx, org, project, name)
	})
}
