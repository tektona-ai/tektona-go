package tektona

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tektona-ai/tektona-go/api"
)

// ListGitCredentialsParams selects the credential scope.
type ListGitCredentialsParams struct {
	Org, Project string
	Scope        *api.ListOrgProjectGitCredentialsParamsScope
}

// CreateGitCredentialParams holds the scope and credential fields.
type CreateGitCredentialParams struct {
	Org, Project  string
	DisplayName   string
	Forge         api.CreateGitCredentialBodyForge
	Name          string
	RepositoryIds *[]string
	Scope         api.CreateGitCredentialBodyScope
	Token         string
}

// UpdateGitCredentialParams holds the scope, query, and credential fields.
type UpdateGitCredentialParams struct {
	Org, Project  string
	Scope         api.UpdateOrgProjectGitCredentialParamsScope
	DisplayName   string
	Forge         api.UpdateGitCredentialBodyForge
	RepositoryIds *[]string
	Token         *string
}

// DeleteGitCredentialParams selects the credential scope.
type DeleteGitCredentialParams struct {
	Org, Project string
	Scope        api.DeleteOrgProjectGitCredentialParamsScope
}

// GitCredential supports project git credential operations.
type GitCredential interface {
	List(ctx context.Context, params *ListGitCredentialsParams) (*api.GitCredentialListBody, error)
	Create(ctx context.Context, params CreateGitCredentialParams) (*api.GitCredentialCreateBody, error)
	Update(ctx context.Context, credentialID string, params UpdateGitCredentialParams) (*api.GitCredentialUpdateBody, error)
	Delete(ctx context.Context, credentialID string, params DeleteGitCredentialParams) error
}

type gitCredentialClient struct {
	raw          *api.ClientWithResponses
	org, project string
}

// GitCredential returns the git credential methods.
func (c *Client) GitCredential() GitCredential {
	return gitCredentialClient{raw: c.raw, org: c.org, project: c.project}
}

func (g gitCredentialClient) address(org, project string) (string, string, error) {
	if org == "" {
		org = g.org
	}
	if project == "" {
		project = g.project
	}
	if org == "" || project == "" {
		return "", "", fmt.Errorf("git credential requires org and project")
	}
	return org, project, nil
}

func (g gitCredentialClient) List(ctx context.Context, params *ListGitCredentialsParams) (*api.GitCredentialListBody, error) {
	var p ListGitCredentialsParams
	if params != nil {
		p = *params
	}
	org, project, err := g.address(p.Org, p.Project)
	if err != nil {
		return nil, err
	}
	return resourceJSON[api.GitCredentialListBody](func() (*http.Response, error) {
		return g.raw.ListOrgProjectGitCredentials(ctx, org, project, &api.ListOrgProjectGitCredentialsParams{Scope: p.Scope})
	}, http.StatusOK, "list git credentials")
}

func (g gitCredentialClient) Create(ctx context.Context, params CreateGitCredentialParams) (*api.GitCredentialCreateBody, error) {
	org, project, err := g.address(params.Org, params.Project)
	if err != nil {
		return nil, err
	}
	if params.Name == "" || params.DisplayName == "" || params.Forge == "" || params.Scope == "" || params.Token == "" {
		return nil, fmt.Errorf("create git credential requires name, display name, forge, scope, and token")
	}
	return resourceJSON[api.GitCredentialCreateBody](func() (*http.Response, error) {
		return g.raw.CreateOrgProjectGitCredential(ctx, org, project, api.CreateOrgProjectGitCredentialJSONRequestBody{
			Name: params.Name, DisplayName: params.DisplayName, Forge: params.Forge,
			Scope: params.Scope, Token: params.Token, RepositoryIds: params.RepositoryIds,
		})
	}, http.StatusCreated, "create git credential")
}

func (g gitCredentialClient) Update(ctx context.Context, credentialID string, params UpdateGitCredentialParams) (*api.GitCredentialUpdateBody, error) {
	org, project, err := g.address(params.Org, params.Project)
	if err != nil {
		return nil, err
	}
	if credentialID == "" || params.Scope == "" || params.DisplayName == "" || params.Forge == "" {
		return nil, fmt.Errorf("update git credential requires ID, scope, display name, and forge")
	}
	return resourceJSON[api.GitCredentialUpdateBody](func() (*http.Response, error) {
		return g.raw.UpdateOrgProjectGitCredential(ctx, org, project, credentialID,
			&api.UpdateOrgProjectGitCredentialParams{Scope: params.Scope}, api.UpdateOrgProjectGitCredentialJSONRequestBody{
				DisplayName: params.DisplayName, Forge: params.Forge,
				RepositoryIds: params.RepositoryIds, Token: params.Token,
			})
	}, http.StatusOK, "update git credential")
}

func (g gitCredentialClient) Delete(ctx context.Context, credentialID string, params DeleteGitCredentialParams) error {
	org, project, err := g.address(params.Org, params.Project)
	if err != nil {
		return err
	}
	if credentialID == "" || params.Scope == "" {
		return fmt.Errorf("delete git credential requires ID and scope")
	}
	return resourceNoContent(func() (*http.Response, error) {
		return g.raw.DeleteOrgProjectGitCredential(ctx, org, project, credentialID,
			&api.DeleteOrgProjectGitCredentialParams{Scope: params.Scope})
	})
}
