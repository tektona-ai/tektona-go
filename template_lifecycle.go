package tektona

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tektona-ai/tektona-go/api"
)

// UpdateTemplateLifecycleForOrgParams holds organization lifecycle settings.
type UpdateTemplateLifecycleForOrgParams = api.UpdateOrgTemplateLifecycleSettingsJSONRequestBody

// PreviewTemplateLifecycleForOrgParams holds proposed organization lifecycle settings.
type PreviewTemplateLifecycleForOrgParams = api.PreviewOrgTemplateLifecycleSettingsJSONRequestBody

// GetTemplateLifecycleForProjectParams holds an optional organization override.
type GetTemplateLifecycleForProjectParams struct{ Org string }

// UpdateTemplateLifecycleForProjectParams holds an organization override and lifecycle settings.
type UpdateTemplateLifecycleForProjectParams struct {
	Org                     string
	ArchiveAfterDaysUnused  *int32
	DeleteAfterDaysArchived *int32
}

// PreviewTemplateLifecycleForProjectParams holds an organization override and proposed lifecycle settings.
type PreviewTemplateLifecycleForProjectParams struct {
	Org                     string
	ArchiveAfterDaysUnused  *int32
	DeleteAfterDaysArchived *int32
}

// TemplateLifecycle supports organization and project template lifecycle settings.
type TemplateLifecycle interface {
	GetForOrg(ctx context.Context, org string) (*api.TemplateLifecycleSettingsResponse, error)
	UpdateForOrg(ctx context.Context, org string, params UpdateTemplateLifecycleForOrgParams) (*api.TemplateLifecycleSettingsResponse, error)
	PreviewForOrg(ctx context.Context, org string, params PreviewTemplateLifecycleForOrgParams) (*api.TemplateLifecyclePreviewResponse, error)
	GetForProject(ctx context.Context, project string, params *GetTemplateLifecycleForProjectParams) (*api.TemplateLifecycleSettingsResponse, error)
	UpdateForProject(ctx context.Context, project string, params UpdateTemplateLifecycleForProjectParams) (*api.TemplateLifecycleSettingsResponse, error)
	PreviewForProject(ctx context.Context, project string, params PreviewTemplateLifecycleForProjectParams) (*api.TemplateLifecyclePreviewResponse, error)
}

type templateLifecycleClient struct {
	raw *api.ClientWithResponses
	org string
}

// TemplateLifecycle returns the template lifecycle methods.
func (c *Client) TemplateLifecycle() TemplateLifecycle {
	return &templateLifecycleClient{raw: c.raw, org: c.org}
}

func (t *templateLifecycleClient) address(org, project string) (string, string, error) {
	if org == "" {
		org = t.org
	}
	if org == "" || project == "" {
		return "", "", fmt.Errorf("template lifecycle requires org and project")
	}
	return org, project, nil
}

func (t *templateLifecycleClient) orgName(org string) (string, error) {
	if org == "" {
		org = t.org
	}
	if org == "" {
		return "", fmt.Errorf("template lifecycle requires org")
	}
	return org, nil
}

func (t *templateLifecycleClient) GetForOrg(ctx context.Context, org string) (*api.TemplateLifecycleSettingsResponse, error) {
	org, err := t.orgName(org)
	if err != nil {
		return nil, err
	}
	return resourceJSON[api.TemplateLifecycleSettingsResponse](func() (*http.Response, error) {
		return t.raw.GetOrgTemplateLifecycleSettings(ctx, org)
	}, http.StatusOK, "get organization template lifecycle")
}

func (t *templateLifecycleClient) UpdateForOrg(ctx context.Context, org string, params UpdateTemplateLifecycleForOrgParams) (*api.TemplateLifecycleSettingsResponse, error) {
	org, err := t.orgName(org)
	if err != nil {
		return nil, err
	}
	return resourceJSON[api.TemplateLifecycleSettingsResponse](func() (*http.Response, error) {
		return t.raw.UpdateOrgTemplateLifecycleSettings(ctx, org, params)
	}, http.StatusOK, "update organization template lifecycle")
}

func (t *templateLifecycleClient) PreviewForOrg(ctx context.Context, org string, params PreviewTemplateLifecycleForOrgParams) (*api.TemplateLifecyclePreviewResponse, error) {
	org, err := t.orgName(org)
	if err != nil {
		return nil, err
	}
	return resourceJSON[api.TemplateLifecyclePreviewResponse](func() (*http.Response, error) {
		return t.raw.PreviewOrgTemplateLifecycleSettings(ctx, org, params)
	}, http.StatusOK, "preview organization template lifecycle")
}

func (t *templateLifecycleClient) GetForProject(ctx context.Context, project string, params *GetTemplateLifecycleForProjectParams) (*api.TemplateLifecycleSettingsResponse, error) {
	var org string
	if params != nil {
		org = params.Org
	}
	org, project, err := t.address(org, project)
	if err != nil {
		return nil, err
	}
	return resourceJSON[api.TemplateLifecycleSettingsResponse](func() (*http.Response, error) {
		return t.raw.GetProjectTemplateLifecycleSettings(ctx, org, project)
	}, http.StatusOK, "get project template lifecycle")
}

func (t *templateLifecycleClient) UpdateForProject(ctx context.Context, project string, params UpdateTemplateLifecycleForProjectParams) (*api.TemplateLifecycleSettingsResponse, error) {
	org, project, err := t.address(params.Org, project)
	if err != nil {
		return nil, err
	}
	return resourceJSON[api.TemplateLifecycleSettingsResponse](func() (*http.Response, error) {
		return t.raw.UpdateProjectTemplateLifecycleSettings(ctx, org, project, api.UpdateProjectTemplateLifecycleSettingsJSONRequestBody{
			ArchiveAfterDaysUnused: params.ArchiveAfterDaysUnused, DeleteAfterDaysArchived: params.DeleteAfterDaysArchived,
		})
	}, http.StatusOK, "update project template lifecycle")
}

func (t *templateLifecycleClient) PreviewForProject(ctx context.Context, project string, params PreviewTemplateLifecycleForProjectParams) (*api.TemplateLifecyclePreviewResponse, error) {
	org, project, err := t.address(params.Org, project)
	if err != nil {
		return nil, err
	}
	return resourceJSON[api.TemplateLifecyclePreviewResponse](func() (*http.Response, error) {
		return t.raw.PreviewProjectTemplateLifecycleSettings(ctx, org, project, api.PreviewProjectTemplateLifecycleSettingsJSONRequestBody{
			ArchiveAfterDaysUnused: params.ArchiveAfterDaysUnused, DeleteAfterDaysArchived: params.DeleteAfterDaysArchived,
		})
	}, http.StatusOK, "preview project template lifecycle")
}
