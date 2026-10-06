package tektona

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/tektona-ai/tektona-go/api"
)

// ListRepositoriesParams holds scope and a default-set filter.
type ListRepositoriesParams struct {
	Org, Project string
	Default      *bool
}

// CreateRepositoryParams holds scope and repository fields.
type CreateRepositoryParams struct {
	Org, Project  string
	Name          string
	Url           string
	DefaultBranch *string
	IsDefault     *bool
}

// UpdateRepositoryParams holds scope and repository fields.
type UpdateRepositoryParams = CreateRepositoryParams

// DeleteRepositoryParams holds an optional scope override.
type DeleteRepositoryParams struct{ Org, Project string }

// ListRepositoryBranchesParams holds scope, filters, and pagination options.
type ListRepositoryBranchesParams struct {
	Org, Project string
	Refresh      *bool
	Q            *string
	Cursor       *string
	Limit        *int32
}

// Repository supports project repository and branch operations.
type Repository interface {
	List(ctx context.Context, params *ListRepositoriesParams) (*api.ListBody, error)
	Create(ctx context.Context, params CreateRepositoryParams) (*api.Response, error)
	Update(ctx context.Context, repositoryID string, params UpdateRepositoryParams) (*api.Response, error)
	Delete(ctx context.Context, repositoryID string, params *DeleteRepositoryParams) error
	ListBranches(ctx context.Context, repositoryID string, params *ListRepositoryBranchesParams) (*api.BranchesBody, error)
}

type repositoryClient struct {
	raw          *api.ClientWithResponses
	org, project string
}

// Repository returns the project repository methods.
func (c *Client) Repository() Repository {
	return repositoryClient{raw: c.raw, org: c.org, project: c.project}
}

func (r repositoryClient) address(org, project string) (string, string, error) {
	if org == "" {
		org = r.org
	}
	if project == "" {
		project = r.project
	}
	if org == "" || project == "" {
		return "", "", fmt.Errorf("repository requires org and project")
	}
	return org, project, nil
}

func (r repositoryClient) List(ctx context.Context, params *ListRepositoriesParams) (*api.ListBody, error) {
	var p ListRepositoriesParams
	if params != nil {
		p = *params
	}
	org, project, err := r.address(p.Org, p.Project)
	if err != nil {
		return nil, err
	}
	return resourceJSON[api.ListBody](func() (*http.Response, error) {
		return r.raw.ListOrgProjectRepositories(ctx, org, project, &api.ListOrgProjectRepositoriesParams{Default: p.Default})
	}, http.StatusOK, "list repositories")
}

func (r repositoryClient) Create(ctx context.Context, params CreateRepositoryParams) (*api.Response, error) {
	org, project, err := r.address(params.Org, params.Project)
	if err != nil {
		return nil, err
	}
	if params.Name == "" || params.Url == "" {
		return nil, fmt.Errorf("create repository requires name and URL")
	}
	return resourceJSON[api.Response](func() (*http.Response, error) {
		return r.raw.CreateOrgProjectRepository(ctx, org, project, api.CreateOrgProjectRepositoryJSONRequestBody{
			Name: params.Name, Url: params.Url, DefaultBranch: params.DefaultBranch, IsDefault: params.IsDefault,
		})
	}, http.StatusCreated, "create repository")
}

func (r repositoryClient) Update(ctx context.Context, repositoryID string, params UpdateRepositoryParams) (*api.Response, error) {
	org, project, err := r.address(params.Org, params.Project)
	if err != nil {
		return nil, err
	}
	id, err := uuid.Parse(repositoryID)
	if err != nil || params.Name == "" || params.Url == "" {
		return nil, fmt.Errorf("update repository requires a valid ID, name, and URL")
	}
	return resourceJSON[api.Response](func() (*http.Response, error) {
		return r.raw.UpdateOrgProjectRepository(ctx, org, project, id, api.UpdateOrgProjectRepositoryJSONRequestBody{
			Name: params.Name, Url: params.Url, DefaultBranch: params.DefaultBranch, IsDefault: params.IsDefault,
		})
	}, http.StatusOK, "update repository")
}

func (r repositoryClient) Delete(ctx context.Context, repositoryID string, params *DeleteRepositoryParams) error {
	var p DeleteRepositoryParams
	if params != nil {
		p = *params
	}
	org, project, err := r.address(p.Org, p.Project)
	if err != nil {
		return err
	}
	id, err := uuid.Parse(repositoryID)
	if err != nil {
		return fmt.Errorf("delete repository requires a valid ID: %w", err)
	}
	return resourceNoContent(func() (*http.Response, error) {
		return r.raw.DeleteOrgProjectRepository(ctx, org, project, id)
	})
}

func (r repositoryClient) ListBranches(ctx context.Context, repositoryID string, params *ListRepositoryBranchesParams) (*api.BranchesBody, error) {
	var p ListRepositoryBranchesParams
	if params != nil {
		p = *params
	}
	org, project, err := r.address(p.Org, p.Project)
	if err != nil {
		return nil, err
	}
	id, err := uuid.Parse(repositoryID)
	if err != nil {
		return nil, fmt.Errorf("list repository branches requires a valid ID: %w", err)
	}
	return resourceJSON[api.BranchesBody](func() (*http.Response, error) {
		return r.raw.ListOrgProjectRepositoryBranches(ctx, org, project, id, &api.ListOrgProjectRepositoryBranchesParams{
			Refresh: p.Refresh, Q: p.Q, Cursor: p.Cursor, Limit: p.Limit,
		})
	}, http.StatusOK, "list repository branches")
}
