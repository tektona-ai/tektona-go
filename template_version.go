package tektona

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tektona-ai/tektona-go/api"
)

// ListTemplateVersionsForOrgParams holds scope, filters, and pagination.
type ListTemplateVersionsForOrgParams struct {
	Org    string
	State  *api.ListOrgTemplateVersionsParamsState
	Tagged *bool
	Q      *string
	Cursor *string
	Limit  *int32
}

// ListTemplateVersionsForProjectParams holds scope, filters, and pagination.
type ListTemplateVersionsForProjectParams struct {
	Org     string
	Project string
	State   *api.ListProjectTemplateVersionsParamsState
	Tagged  *bool
	Q       *string
	Cursor  *string
	Limit   *int32
}

// ListTemplateVersionsSystemParams holds filters and pagination.
type ListTemplateVersionsSystemParams = api.ListSystemTemplateVersionsParams

// TemplateVersionForOrgParams holds an organization override.
type TemplateVersionForOrgParams struct{ Org string }

// TemplateVersionForProjectParams holds organization and project overrides.
type TemplateVersionForProjectParams struct{ Org, Project string }

// PruneTemplateVersionsForOrgParams holds scope and the prune body.
type PruneTemplateVersionsForOrgParams struct {
	Org    string
	DryRun *bool
}

// PruneTemplateVersionsForProjectParams holds scope and the prune body.
type PruneTemplateVersionsForProjectParams struct {
	Org     string
	Project string
	DryRun  *bool
}

// TemplateVersion supports organization, project, and system template versions.
type TemplateVersion interface {
	ListForOrg(context.Context, string, *ListTemplateVersionsForOrgParams) (*api.ListTemplateVersionsBody, error)
	GetForOrg(context.Context, string, string, *TemplateVersionForOrgParams) (*api.TemplateVersion, error)
	DeleteForOrg(context.Context, string, string, *TemplateVersionForOrgParams) error
	PruneForOrg(context.Context, string, PruneTemplateVersionsForOrgParams) (*api.PruneVersionsReport, error)
	ActivateForOrg(context.Context, string, string, *TemplateVersionForOrgParams) (*api.TemplateVersion, error)
	ArchiveForOrg(context.Context, string, string, *TemplateVersionForOrgParams) (*api.TemplateVersion, error)
	ListForProject(context.Context, string, *ListTemplateVersionsForProjectParams) (*api.ListTemplateVersionsBody, error)
	GetForProject(context.Context, string, string, *TemplateVersionForProjectParams) (*api.TemplateVersion, error)
	DeleteForProject(context.Context, string, string, *TemplateVersionForProjectParams) error
	PruneForProject(context.Context, string, PruneTemplateVersionsForProjectParams) (*api.PruneVersionsReport, error)
	ActivateForProject(context.Context, string, string, *TemplateVersionForProjectParams) (*api.TemplateVersion, error)
	ArchiveForProject(context.Context, string, string, *TemplateVersionForProjectParams) (*api.TemplateVersion, error)
	ListSystem(context.Context, string, *ListTemplateVersionsSystemParams) (*api.ListTemplateVersionsBody, error)
	GetSystem(context.Context, string, string) (*api.TemplateVersion, error)
}

type templateVersionClient struct {
	raw          *api.ClientWithResponses
	org, project string
}

// TemplateVersion returns the template version methods.
func (c *Client) TemplateVersion() TemplateVersion {
	return &templateVersionClient{raw: c.raw, org: c.org, project: c.project}
}

func (v *templateVersionClient) orgAddress(org, name string) (string, error) {
	if org == "" {
		org = v.org
	}
	if org == "" || name == "" {
		return "", fmt.Errorf("template version requires org and template name")
	}
	return org, nil
}

func (v *templateVersionClient) projectAddress(org, project, name string) (string, string, error) {
	org, err := v.orgAddress(org, name)
	if err != nil {
		return "", "", err
	}
	if project == "" {
		project = v.project
	}
	if project == "" {
		return "", "", fmt.Errorf("template version requires project")
	}
	return org, project, nil
}

func versionID(versionID string) error {
	if versionID == "" {
		return fmt.Errorf("template version requires version ID")
	}
	return nil
}

func (v *templateVersionClient) ListForOrg(ctx context.Context, name string, params *ListTemplateVersionsForOrgParams) (*api.ListTemplateVersionsBody, error) {
	var p ListTemplateVersionsForOrgParams
	if params != nil {
		p = *params
	}
	org, err := v.orgAddress(p.Org, name)
	if err != nil {
		return nil, err
	}
	query := &api.ListOrgTemplateVersionsParams{State: p.State, Tagged: p.Tagged, Q: p.Q, Cursor: p.Cursor, Limit: p.Limit}
	return resourceJSON[api.ListTemplateVersionsBody](func() (*http.Response, error) {
		return v.raw.ListOrgTemplateVersions(ctx, org, name, query)
	}, http.StatusOK, "list organization template versions")
}

func (v *templateVersionClient) GetForOrg(ctx context.Context, name, id string, params *TemplateVersionForOrgParams) (*api.TemplateVersion, error) {
	org, err := v.orgVersionAddress(name, id, params)
	if err != nil {
		return nil, err
	}
	return resourceJSON[api.TemplateVersion](func() (*http.Response, error) {
		return v.raw.GetOrgTemplateVersion(ctx, org, name, id)
	}, http.StatusOK, "get organization template version")
}

func (v *templateVersionClient) DeleteForOrg(ctx context.Context, name, id string, params *TemplateVersionForOrgParams) error {
	org, err := v.orgVersionAddress(name, id, params)
	if err != nil {
		return err
	}
	return resourceNoContent(func() (*http.Response, error) {
		return v.raw.DeleteOrgTemplateVersion(ctx, org, name, id)
	})
}

func (v *templateVersionClient) PruneForOrg(ctx context.Context, name string, params PruneTemplateVersionsForOrgParams) (*api.PruneVersionsReport, error) {
	org, err := v.orgAddress(params.Org, name)
	if err != nil {
		return nil, err
	}
	return resourceJSON[api.PruneVersionsReport](func() (*http.Response, error) {
		return v.raw.PruneOrgTemplateVersions(ctx, org, name, api.PruneOrgTemplateVersionsJSONRequestBody{DryRun: params.DryRun})
	}, http.StatusOK, "prune organization template versions")
}

func (v *templateVersionClient) ActivateForOrg(ctx context.Context, name, id string, params *TemplateVersionForOrgParams) (*api.TemplateVersion, error) {
	org, err := v.orgVersionAddress(name, id, params)
	if err != nil {
		return nil, err
	}
	return resourceJSON[api.TemplateVersion](func() (*http.Response, error) {
		return v.raw.ActivateOrgTemplateVersion(ctx, org, name, id)
	}, http.StatusOK, "activate organization template version")
}

func (v *templateVersionClient) ArchiveForOrg(ctx context.Context, name, id string, params *TemplateVersionForOrgParams) (*api.TemplateVersion, error) {
	org, err := v.orgVersionAddress(name, id, params)
	if err != nil {
		return nil, err
	}
	return resourceJSON[api.TemplateVersion](func() (*http.Response, error) {
		return v.raw.ArchiveOrgTemplateVersion(ctx, org, name, id)
	}, http.StatusOK, "archive organization template version")
}

func (v *templateVersionClient) orgVersionAddress(name, id string, params *TemplateVersionForOrgParams) (string, error) {
	var org string
	if params != nil {
		org = params.Org
	}
	org, err := v.orgAddress(org, name)
	if err != nil {
		return "", err
	}
	return org, versionID(id)
}

func (v *templateVersionClient) ListForProject(ctx context.Context, name string, params *ListTemplateVersionsForProjectParams) (*api.ListTemplateVersionsBody, error) {
	var p ListTemplateVersionsForProjectParams
	if params != nil {
		p = *params
	}
	org, project, err := v.projectAddress(p.Org, p.Project, name)
	if err != nil {
		return nil, err
	}
	query := &api.ListProjectTemplateVersionsParams{State: p.State, Tagged: p.Tagged, Q: p.Q, Cursor: p.Cursor, Limit: p.Limit}
	return resourceJSON[api.ListTemplateVersionsBody](func() (*http.Response, error) {
		return v.raw.ListProjectTemplateVersions(ctx, org, project, name, query)
	}, http.StatusOK, "list project template versions")
}

func (v *templateVersionClient) projectVersionAddress(name, id string, params *TemplateVersionForProjectParams) (string, string, error) {
	var org, project string
	if params != nil {
		org, project = params.Org, params.Project
	}
	org, project, err := v.projectAddress(org, project, name)
	if err != nil {
		return "", "", err
	}
	if err := versionID(id); err != nil {
		return "", "", err
	}
	return org, project, nil
}

func (v *templateVersionClient) GetForProject(ctx context.Context, name, id string, params *TemplateVersionForProjectParams) (*api.TemplateVersion, error) {
	org, project, err := v.projectVersionAddress(name, id, params)
	if err != nil {
		return nil, err
	}
	return resourceJSON[api.TemplateVersion](func() (*http.Response, error) {
		return v.raw.GetProjectTemplateVersion(ctx, org, project, name, id)
	}, http.StatusOK, "get project template version")
}

func (v *templateVersionClient) DeleteForProject(ctx context.Context, name, id string, params *TemplateVersionForProjectParams) error {
	org, project, err := v.projectVersionAddress(name, id, params)
	if err != nil {
		return err
	}
	return resourceNoContent(func() (*http.Response, error) {
		return v.raw.DeleteProjectTemplateVersion(ctx, org, project, name, id)
	})
}

func (v *templateVersionClient) PruneForProject(ctx context.Context, name string, params PruneTemplateVersionsForProjectParams) (*api.PruneVersionsReport, error) {
	org, project, err := v.projectAddress(params.Org, params.Project, name)
	if err != nil {
		return nil, err
	}
	return resourceJSON[api.PruneVersionsReport](func() (*http.Response, error) {
		return v.raw.PruneProjectTemplateVersions(ctx, org, project, name, api.PruneProjectTemplateVersionsJSONRequestBody{DryRun: params.DryRun})
	}, http.StatusOK, "prune project template versions")
}

func (v *templateVersionClient) ActivateForProject(ctx context.Context, name, id string, params *TemplateVersionForProjectParams) (*api.TemplateVersion, error) {
	org, project, err := v.projectVersionAddress(name, id, params)
	if err != nil {
		return nil, err
	}
	return resourceJSON[api.TemplateVersion](func() (*http.Response, error) {
		return v.raw.ActivateProjectTemplateVersion(ctx, org, project, name, id)
	}, http.StatusOK, "activate project template version")
}

func (v *templateVersionClient) ArchiveForProject(ctx context.Context, name, id string, params *TemplateVersionForProjectParams) (*api.TemplateVersion, error) {
	org, project, err := v.projectVersionAddress(name, id, params)
	if err != nil {
		return nil, err
	}
	return resourceJSON[api.TemplateVersion](func() (*http.Response, error) {
		return v.raw.ArchiveProjectTemplateVersion(ctx, org, project, name, id)
	}, http.StatusOK, "archive project template version")
}

func (v *templateVersionClient) ListSystem(ctx context.Context, name string, params *ListTemplateVersionsSystemParams) (*api.ListTemplateVersionsBody, error) {
	if name == "" {
		return nil, fmt.Errorf("list system template versions requires template name")
	}
	return resourceJSON[api.ListTemplateVersionsBody](func() (*http.Response, error) {
		return v.raw.ListSystemTemplateVersions(ctx, name, params)
	}, http.StatusOK, "list system template versions")
}

func (v *templateVersionClient) GetSystem(ctx context.Context, name, id string) (*api.TemplateVersion, error) {
	if name == "" {
		return nil, fmt.Errorf("get system template version requires template name")
	}
	if err := versionID(id); err != nil {
		return nil, err
	}
	return resourceJSON[api.TemplateVersion](func() (*http.Response, error) {
		return v.raw.GetSystemTemplateVersion(ctx, name, id)
	}, http.StatusOK, "get system template version")
}
