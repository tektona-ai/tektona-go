package tektona

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tektona-ai/tektona-go/api"
)

// ListTemplatesForOrgParams holds filters and pagination for organization templates.
type ListTemplatesForOrgParams = api.ListOrgTemplatesParams

// ListTemplatesForProjectParams holds filters and pagination for project templates.
type ListTemplatesForProjectParams = api.ListProjectTemplatesParams

// CreateTemplateForOrgParams holds the organization and template metadata.
type CreateTemplateForOrgParams struct {
	Org         string
	Name        string
	DisplayName *string
	Description *string
}

// CreateTemplateForProjectParams holds the project and template metadata.
type CreateTemplateForProjectParams struct {
	Org         string
	Project     string
	Name        string
	DisplayName *string
	Description *string
}

// UpdateTemplateForOrgParams holds the organization and replacement metadata.
type UpdateTemplateForOrgParams struct {
	Org         string
	DisplayName *string
	Description *string
}

// UpdateTemplateForProjectParams holds the project and replacement metadata.
type UpdateTemplateForProjectParams struct {
	Org         string
	Project     string
	DisplayName *string
	Description *string
}

// GetTemplateForOrgParams holds an optional organization override.
type GetTemplateForOrgParams struct{ Org string }

// DeleteTemplateForOrgParams holds an optional organization override.
type DeleteTemplateForOrgParams struct{ Org string }

// ActivateTemplateForOrgParams holds an optional organization override.
type ActivateTemplateForOrgParams struct{ Org string }

// ArchiveTemplateForOrgParams holds an optional organization override.
type ArchiveTemplateForOrgParams struct{ Org string }

// GetTemplateForProjectParams holds optional project scope overrides.
type GetTemplateForProjectParams struct{ Org, Project string }

// DeleteTemplateForProjectParams holds optional project scope overrides.
type DeleteTemplateForProjectParams struct{ Org, Project string }

// ActivateTemplateForProjectParams holds optional project scope overrides.
type ActivateTemplateForProjectParams struct{ Org, Project string }

// ArchiveTemplateForProjectParams holds optional project scope overrides.
type ArchiveTemplateForProjectParams struct{ Org, Project string }

// Template supports organization, project, and system template operations.
type Template interface {
	ListForOrg(ctx context.Context, org string, params *ListTemplatesForOrgParams) (*api.ListTemplatesBody, error)
	CreateForOrg(ctx context.Context, params CreateTemplateForOrgParams) (*api.TemplateResponse, error)
	GetForOrg(ctx context.Context, name string, params *GetTemplateForOrgParams) (*api.TemplateResponse, error)
	UpdateForOrg(ctx context.Context, name string, params UpdateTemplateForOrgParams) (*api.TemplateResponse, error)
	DeleteForOrg(ctx context.Context, name string, params *DeleteTemplateForOrgParams) (*api.DeleteTemplateResponse, error)
	ActivateForOrg(ctx context.Context, name string, params *ActivateTemplateForOrgParams) (*api.TemplateResponse, error)
	ArchiveForOrg(ctx context.Context, name string, params *ArchiveTemplateForOrgParams) (*api.TemplateResponse, error)
	ListForProject(ctx context.Context, org, project string, params *ListTemplatesForProjectParams) (*api.ListTemplatesBody, error)
	CreateForProject(ctx context.Context, params CreateTemplateForProjectParams) (*api.TemplateResponse, error)
	GetForProject(ctx context.Context, name string, params *GetTemplateForProjectParams) (*api.TemplateResponse, error)
	UpdateForProject(ctx context.Context, name string, params UpdateTemplateForProjectParams) (*api.TemplateResponse, error)
	DeleteForProject(ctx context.Context, name string, params *DeleteTemplateForProjectParams) (*api.DeleteTemplateResponse, error)
	ActivateForProject(ctx context.Context, name string, params *ActivateTemplateForProjectParams) (*api.TemplateResponse, error)
	ArchiveForProject(ctx context.Context, name string, params *ArchiveTemplateForProjectParams) (*api.TemplateResponse, error)
	GetSystem(ctx context.Context, name string) (*api.TemplateResponse, error)
}

type templateClient struct {
	raw          *api.ClientWithResponses
	org, project string
}

// Template returns the template methods with the client's scope defaults.
func (c *Client) Template() Template {
	return &templateClient{raw: c.raw, org: c.org, project: c.project}
}

func (t *templateClient) orgScope(org string) (string, error) {
	if org == "" {
		org = t.org
	}
	if org == "" {
		return "", fmt.Errorf("template requires org")
	}
	return org, nil
}

func (t *templateClient) projectScope(org, project string) (string, string, error) {
	var err error
	org, err = t.orgScope(org)
	if project == "" {
		project = t.project
	}
	if err != nil || project == "" {
		return "", "", fmt.Errorf("template requires org and project")
	}
	return org, project, nil
}

func templateName(name string) error {
	if name == "" {
		return fmt.Errorf("template requires name")
	}
	return nil
}

func (t *templateClient) ListForOrg(ctx context.Context, org string, params *ListTemplatesForOrgParams) (*api.ListTemplatesBody, error) {
	org, err := t.orgScope(org)
	if err != nil {
		return nil, err
	}
	return resourceJSON[api.ListTemplatesBody](func() (*http.Response, error) {
		return t.raw.ListOrgTemplates(ctx, org, params)
	}, http.StatusOK, "list organization templates")
}

func (t *templateClient) CreateForOrg(ctx context.Context, params CreateTemplateForOrgParams) (*api.TemplateResponse, error) {
	org, err := t.orgScope(params.Org)
	if err != nil {
		return nil, err
	}
	if err := templateName(params.Name); err != nil {
		return nil, err
	}
	return resourceJSON[api.TemplateResponse](func() (*http.Response, error) {
		return t.raw.CreateOrgTemplate(ctx, org, api.CreateOrgTemplateJSONRequestBody{Metadata: api.TemplateMetadata{
			Name: params.Name, DisplayName: params.DisplayName, Description: params.Description,
		}})
	}, http.StatusCreated, "create organization template")
}

func (t *templateClient) GetForOrg(ctx context.Context, name string, params *GetTemplateForOrgParams) (*api.TemplateResponse, error) {
	var org string
	if params != nil {
		org = params.Org
	}
	org, err := t.orgScope(org)
	if err != nil {
		return nil, err
	}
	if err := templateName(name); err != nil {
		return nil, err
	}
	return resourceJSON[api.TemplateResponse](func() (*http.Response, error) {
		return t.raw.GetOrgTemplate(ctx, org, name)
	}, http.StatusOK, "get organization template")
}

func (t *templateClient) UpdateForOrg(ctx context.Context, name string, params UpdateTemplateForOrgParams) (*api.TemplateResponse, error) {
	org, err := t.orgScope(params.Org)
	if err != nil {
		return nil, err
	}
	if err := templateName(name); err != nil {
		return nil, err
	}
	return resourceJSON[api.TemplateResponse](func() (*http.Response, error) {
		return t.raw.UpdateOrgTemplate(ctx, org, name, api.UpdateOrgTemplateJSONRequestBody{Metadata: api.TemplateMetadata{
			Name: name, DisplayName: params.DisplayName, Description: params.Description,
		}})
	}, http.StatusOK, "update organization template")
}

func (t *templateClient) DeleteForOrg(ctx context.Context, name string, params *DeleteTemplateForOrgParams) (*api.DeleteTemplateResponse, error) {
	var org string
	if params != nil {
		org = params.Org
	}
	org, err := t.orgScope(org)
	if err != nil {
		return nil, err
	}
	if err := templateName(name); err != nil {
		return nil, err
	}
	return resourceJSON[api.DeleteTemplateResponse](func() (*http.Response, error) {
		return t.raw.DeleteOrgTemplate(ctx, org, name)
	}, http.StatusOK, "delete organization template")
}

func (t *templateClient) ActivateForOrg(ctx context.Context, name string, params *ActivateTemplateForOrgParams) (*api.TemplateResponse, error) {
	var org string
	if params != nil {
		org = params.Org
	}
	org, err := t.orgScope(org)
	if err != nil {
		return nil, err
	}
	if err := templateName(name); err != nil {
		return nil, err
	}
	return resourceJSON[api.TemplateResponse](func() (*http.Response, error) {
		return t.raw.ActivateOrgTemplate(ctx, org, name)
	}, http.StatusOK, "activate organization template")
}

func (t *templateClient) ArchiveForOrg(ctx context.Context, name string, params *ArchiveTemplateForOrgParams) (*api.TemplateResponse, error) {
	var org string
	if params != nil {
		org = params.Org
	}
	org, err := t.orgScope(org)
	if err != nil {
		return nil, err
	}
	if err := templateName(name); err != nil {
		return nil, err
	}
	return resourceJSON[api.TemplateResponse](func() (*http.Response, error) {
		return t.raw.ArchiveOrgTemplate(ctx, org, name)
	}, http.StatusOK, "archive organization template")
}

func (t *templateClient) ListForProject(ctx context.Context, org, project string, params *ListTemplatesForProjectParams) (*api.ListTemplatesBody, error) {
	org, project, err := t.projectScope(org, project)
	if err != nil {
		return nil, err
	}
	return resourceJSON[api.ListTemplatesBody](func() (*http.Response, error) {
		return t.raw.ListProjectTemplates(ctx, org, project, params)
	}, http.StatusOK, "list project templates")
}

func (t *templateClient) CreateForProject(ctx context.Context, params CreateTemplateForProjectParams) (*api.TemplateResponse, error) {
	org, project, err := t.projectScope(params.Org, params.Project)
	if err != nil {
		return nil, err
	}
	if err := templateName(params.Name); err != nil {
		return nil, err
	}
	return resourceJSON[api.TemplateResponse](func() (*http.Response, error) {
		return t.raw.CreateProjectTemplate(ctx, org, project, api.CreateProjectTemplateJSONRequestBody{Metadata: api.TemplateMetadata{
			Name: params.Name, DisplayName: params.DisplayName, Description: params.Description,
		}})
	}, http.StatusCreated, "create project template")
}

func (t *templateClient) GetForProject(ctx context.Context, name string, params *GetTemplateForProjectParams) (*api.TemplateResponse, error) {
	var org, project string
	if params != nil {
		org, project = params.Org, params.Project
	}
	org, project, err := t.projectScope(org, project)
	if err != nil {
		return nil, err
	}
	if err := templateName(name); err != nil {
		return nil, err
	}
	return resourceJSON[api.TemplateResponse](func() (*http.Response, error) {
		return t.raw.GetProjectTemplate(ctx, org, project, name)
	}, http.StatusOK, "get project template")
}

func (t *templateClient) UpdateForProject(ctx context.Context, name string, params UpdateTemplateForProjectParams) (*api.TemplateResponse, error) {
	org, project, err := t.projectScope(params.Org, params.Project)
	if err != nil {
		return nil, err
	}
	if err := templateName(name); err != nil {
		return nil, err
	}
	return resourceJSON[api.TemplateResponse](func() (*http.Response, error) {
		return t.raw.UpdateProjectTemplate(ctx, org, project, name, api.UpdateProjectTemplateJSONRequestBody{Metadata: api.TemplateMetadata{
			Name: name, DisplayName: params.DisplayName, Description: params.Description,
		}})
	}, http.StatusOK, "update project template")
}

func (t *templateClient) DeleteForProject(ctx context.Context, name string, params *DeleteTemplateForProjectParams) (*api.DeleteTemplateResponse, error) {
	var org, project string
	if params != nil {
		org, project = params.Org, params.Project
	}
	org, project, err := t.projectScope(org, project)
	if err != nil {
		return nil, err
	}
	if err := templateName(name); err != nil {
		return nil, err
	}
	return resourceJSON[api.DeleteTemplateResponse](func() (*http.Response, error) {
		return t.raw.DeleteProjectTemplate(ctx, org, project, name)
	}, http.StatusOK, "delete project template")
}

func (t *templateClient) ActivateForProject(ctx context.Context, name string, params *ActivateTemplateForProjectParams) (*api.TemplateResponse, error) {
	var org, project string
	if params != nil {
		org, project = params.Org, params.Project
	}
	org, project, err := t.projectScope(org, project)
	if err != nil {
		return nil, err
	}
	if err := templateName(name); err != nil {
		return nil, err
	}
	return resourceJSON[api.TemplateResponse](func() (*http.Response, error) {
		return t.raw.ActivateProjectTemplate(ctx, org, project, name)
	}, http.StatusOK, "activate project template")
}

func (t *templateClient) ArchiveForProject(ctx context.Context, name string, params *ArchiveTemplateForProjectParams) (*api.TemplateResponse, error) {
	var org, project string
	if params != nil {
		org, project = params.Org, params.Project
	}
	org, project, err := t.projectScope(org, project)
	if err != nil {
		return nil, err
	}
	if err := templateName(name); err != nil {
		return nil, err
	}
	return resourceJSON[api.TemplateResponse](func() (*http.Response, error) {
		return t.raw.ArchiveProjectTemplate(ctx, org, project, name)
	}, http.StatusOK, "archive project template")
}

func (t *templateClient) GetSystem(ctx context.Context, name string) (*api.TemplateResponse, error) {
	if err := templateName(name); err != nil {
		return nil, err
	}
	return resourceJSON[api.TemplateResponse](func() (*http.Response, error) {
		return t.raw.GetSystemTemplate(ctx, name)
	}, http.StatusOK, "get system template")
}
