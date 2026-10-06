package tektona

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tektona-ai/tektona-go/api"
)

// ListProjectsParams holds pagination options for projects across organizations.
type ListProjectsParams = api.ListProjectsParams

// ListProjectsForOrgParams holds pagination options for one organization.
type ListProjectsForOrgParams = api.ListOrgProjectsParams

// CreateProjectParams holds the organization and project fields.
type CreateProjectParams struct {
	Org         string
	Name        string
	DisplayName string
	Description *string
}

// GetProjectParams holds an optional organization override.
type GetProjectParams struct{ Org string }

// UpdateProjectParams holds an organization override and project fields.
type UpdateProjectParams struct {
	Org         string
	DisplayName string
	Description *string
}

// DeleteProjectParams holds an optional organization override.
type DeleteProjectParams struct{ Org string }

// GetProjectLifecycleDefaultsParams holds an optional organization override.
type GetProjectLifecycleDefaultsParams struct{ Org string }

// UpdateProjectLifecycleDefaultsParams holds the organization and lifecycle fields.
type UpdateProjectLifecycleDefaultsParams struct {
	Org             string
	AutoDeleteAfter *string
	AutoPauseAfter  *string
	AutoPauseMode   *api.LifecycleDefaultsBodyAutoPauseMode
	AutoResume      *bool
}

// Project supports project and project lifecycle default operations.
type Project interface {
	List(ctx context.Context, params *ListProjectsParams) (*api.ListAccessibleProjectsBody, error)
	ListForOrg(ctx context.Context, org string, params *ListProjectsForOrgParams) (*api.ListProjectsBody, error)
	Create(ctx context.Context, params CreateProjectParams) (*api.ProjectResponse, error)
	Get(ctx context.Context, project string, params *GetProjectParams) (*api.ProjectResponse, error)
	Update(ctx context.Context, project string, params UpdateProjectParams) (*api.ProjectResponse, error)
	Delete(ctx context.Context, project string, params *DeleteProjectParams) error
	GetLifecycleDefaults(ctx context.Context, project string, params *GetProjectLifecycleDefaultsParams) (*api.LifecycleDefaultsBody, error)
	UpdateLifecycleDefaults(ctx context.Context, project string, params UpdateProjectLifecycleDefaultsParams) (*api.LifecycleDefaultsBody, error)
}

type projectClient struct {
	raw          *api.ClientWithResponses
	org, project string
}

// Project returns the project methods.
func (c *Client) Project() Project { return c.projectAPI }

func (p *projectClient) address(org, project string) (string, string, error) {
	if org == "" {
		org = p.org
	}
	if project == "" {
		project = p.project
	}
	if org == "" || project == "" {
		return "", "", fmt.Errorf("project requires org and project")
	}
	return org, project, nil
}

func (p *projectClient) List(ctx context.Context, params *ListProjectsParams) (*api.ListAccessibleProjectsBody, error) {
	resp, err := p.raw.ListProjectsWithResponse(ctx, params)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, responseError(resp.StatusCode(), resp.HTTPResponse.Header, resp.Body, resp.ApplicationproblemJSONDefault)
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("list projects: HTTP 200 lacks a JSON response")
	}
	return resp.JSON200, nil
}

func (p *projectClient) ListForOrg(ctx context.Context, org string, params *ListProjectsForOrgParams) (*api.ListProjectsBody, error) {
	if org == "" {
		org = p.org
	}
	if org == "" {
		return nil, fmt.Errorf("list organization projects requires org")
	}
	resp, err := p.raw.ListOrgProjectsWithResponse(ctx, org, params)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, responseError(resp.StatusCode(), resp.HTTPResponse.Header, resp.Body, resp.ApplicationproblemJSONDefault)
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("list organization projects: HTTP 200 lacks a JSON response")
	}
	return resp.JSON200, nil
}

func (p *projectClient) Create(ctx context.Context, params CreateProjectParams) (*api.ProjectResponse, error) {
	org := params.Org
	if org == "" {
		org = p.org
	}
	if org == "" || params.Name == "" || params.DisplayName == "" {
		return nil, fmt.Errorf("create project requires org, name, and display name")
	}
	resp, err := p.raw.CreateOrgProjectWithResponse(ctx, org, api.CreateOrgProjectJSONRequestBody{
		Name: params.Name, DisplayName: params.DisplayName, Description: params.Description,
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusCreated {
		return nil, responseError(resp.StatusCode(), resp.HTTPResponse.Header, resp.Body, resp.ApplicationproblemJSONDefault)
	}
	if resp.JSON201 == nil {
		return nil, fmt.Errorf("create project: HTTP 201 lacks a JSON response")
	}
	return resp.JSON201, nil
}

func (p *projectClient) Get(ctx context.Context, project string, params *GetProjectParams) (*api.ProjectResponse, error) {
	var org string
	if params != nil {
		org = params.Org
	}
	org, project, err := p.address(org, project)
	if err != nil {
		return nil, err
	}
	resp, err := p.raw.GetOrgProjectWithResponse(ctx, org, project)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, responseError(resp.StatusCode(), resp.HTTPResponse.Header, resp.Body, resp.ApplicationproblemJSONDefault)
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("get project: HTTP 200 lacks a JSON response")
	}
	return resp.JSON200, nil
}

func (p *projectClient) Update(ctx context.Context, project string, params UpdateProjectParams) (*api.ProjectResponse, error) {
	org, project, err := p.address(params.Org, project)
	if err != nil {
		return nil, err
	}
	if params.DisplayName == "" {
		return nil, fmt.Errorf("update project requires display name")
	}
	resp, err := p.raw.UpdateOrgProjectWithResponse(ctx, org, project, api.UpdateOrgProjectJSONRequestBody{
		DisplayName: params.DisplayName, Description: params.Description,
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, responseError(resp.StatusCode(), resp.HTTPResponse.Header, resp.Body, resp.ApplicationproblemJSONDefault)
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("update project: HTTP 200 lacks a JSON response")
	}
	return resp.JSON200, nil
}

func (p *projectClient) Delete(ctx context.Context, project string, params *DeleteProjectParams) error {
	var org string
	if params != nil {
		org = params.Org
	}
	org, project, err := p.address(org, project)
	if err != nil {
		return err
	}
	resp, err := p.raw.DeleteOrgProjectWithResponse(ctx, org, project)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusNoContent {
		return responseError(resp.StatusCode(), resp.HTTPResponse.Header, resp.Body, resp.ApplicationproblemJSONDefault)
	}
	return nil
}

func (p *projectClient) GetLifecycleDefaults(ctx context.Context, project string, params *GetProjectLifecycleDefaultsParams) (*api.LifecycleDefaultsBody, error) {
	var org string
	if params != nil {
		org = params.Org
	}
	org, project, err := p.address(org, project)
	if err != nil {
		return nil, err
	}
	resp, err := p.raw.GetProjectLifecycleDefaultsWithResponse(ctx, org, project)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, responseError(resp.StatusCode(), resp.HTTPResponse.Header, resp.Body, resp.ApplicationproblemJSONDefault)
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("get project lifecycle defaults: HTTP 200 lacks a JSON response")
	}
	return resp.JSON200, nil
}

func (p *projectClient) UpdateLifecycleDefaults(ctx context.Context, project string, params UpdateProjectLifecycleDefaultsParams) (*api.LifecycleDefaultsBody, error) {
	org, project, err := p.address(params.Org, project)
	if err != nil {
		return nil, err
	}
	resp, err := p.raw.UpdateProjectLifecycleDefaultsWithResponse(ctx, org, project, api.UpdateProjectLifecycleDefaultsJSONRequestBody{
		AutoDeleteAfter: params.AutoDeleteAfter, AutoPauseAfter: params.AutoPauseAfter,
		AutoPauseMode: params.AutoPauseMode, AutoResume: params.AutoResume,
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, responseError(resp.StatusCode(), resp.HTTPResponse.Header, resp.Body, resp.ApplicationproblemJSONDefault)
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("update project lifecycle defaults: HTTP 200 lacks a JSON response")
	}
	return resp.JSON200, nil
}
