package tektona

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/tektona-ai/tektona-go/api"
)

// ListTemplateBuildsForOrgParams holds filters and pagination options.
type ListTemplateBuildsForOrgParams = api.ListOrgTemplateBuildsParams

// ListTemplateBuildsForProjectParams holds filters and pagination options.
type ListTemplateBuildsForProjectParams = api.ListProjectTemplateBuildsParams

// CreateTemplateBuildParams holds a build request and its optional retry key.
type CreateTemplateBuildParams struct {
	Org            string
	Project        string
	Template       string
	Build          api.TemplateBuildSpec
	Sandbox        *api.SandboxSpec
	Tag            *string
	IdempotencyKey *string
}

// GetTemplateBuildLogsParams holds log filters and pagination options.
type GetTemplateBuildLogsParams = api.GetTemplateBuildLogsParams

// StreamTemplateBuildLogsParams holds the stream cursor and last event ID.
type StreamTemplateBuildLogsParams = api.StreamTemplateBuildLogsParams

// TemplateBuild supports template build operations.
type TemplateBuild interface {
	ListForOrg(ctx context.Context, org string, params *ListTemplateBuildsForOrgParams) (*api.ListTemplateBuildsBody, error)
	CreateForOrg(ctx context.Context, params CreateTemplateBuildParams) (*api.TemplateBuildResponse, error)
	ListForProject(ctx context.Context, org, project string, params *ListTemplateBuildsForProjectParams) (*api.ListTemplateBuildsBody, error)
	CreateForProject(ctx context.Context, params CreateTemplateBuildParams) (*api.TemplateBuildResponse, error)
	Get(ctx context.Context, buildID string) (*api.TemplateBuildResponse, error)
	Cancel(ctx context.Context, buildID string) (*api.TemplateBuildResponse, error)
	GetLogs(ctx context.Context, buildID string, params *GetTemplateBuildLogsParams) (*api.TemplateBuildLogsBody, error)
	StreamLogs(ctx context.Context, buildID string, params *StreamTemplateBuildLogsParams) (io.ReadCloser, error)
}

type templateBuildClient struct {
	raw          *api.ClientWithResponses
	org, project string
}

// TemplateBuild returns the template build methods.
func (c *Client) TemplateBuild() TemplateBuild {
	return &templateBuildClient{raw: c.raw, org: c.org, project: c.project}
}

func (b *templateBuildClient) orgName(org string) (string, error) {
	if org == "" {
		org = b.org
	}
	if org == "" {
		return "", fmt.Errorf("template build requires org")
	}
	return org, nil
}

func (b *templateBuildClient) projectName(org, project string) (string, string, error) {
	var err error
	org, err = b.orgName(org)
	if project == "" {
		project = b.project
	}
	if err != nil || project == "" {
		return "", "", fmt.Errorf("project template build requires org and project")
	}
	return org, project, nil
}

func (b *templateBuildClient) ListForOrg(ctx context.Context, org string, params *ListTemplateBuildsForOrgParams) (*api.ListTemplateBuildsBody, error) {
	org, err := b.orgName(org)
	if err != nil {
		return nil, err
	}
	return resourceJSON[api.ListTemplateBuildsBody](func() (*http.Response, error) {
		return b.raw.ListOrgTemplateBuilds(ctx, org, params)
	}, http.StatusOK, "list organization template builds")
}

func (b *templateBuildClient) CreateForOrg(ctx context.Context, params CreateTemplateBuildParams) (*api.TemplateBuildResponse, error) {
	org, err := b.orgName(params.Org)
	if err != nil {
		return nil, err
	}
	if params.Template == "" {
		return nil, fmt.Errorf("create template build requires template")
	}
	return resourceJSON[api.TemplateBuildResponse](func() (*http.Response, error) {
		return b.raw.CreateOrgTemplateBuild(ctx, org, &api.CreateOrgTemplateBuildParams{IdempotencyKey: params.IdempotencyKey}, api.CreateOrgTemplateBuildJSONRequestBody{
			Template: params.Template, Build: params.Build, Sandbox: params.Sandbox, Tag: params.Tag,
		})
	}, http.StatusAccepted, "create organization template build")
}

func (b *templateBuildClient) ListForProject(ctx context.Context, org, project string, params *ListTemplateBuildsForProjectParams) (*api.ListTemplateBuildsBody, error) {
	org, project, err := b.projectName(org, project)
	if err != nil {
		return nil, err
	}
	return resourceJSON[api.ListTemplateBuildsBody](func() (*http.Response, error) {
		return b.raw.ListProjectTemplateBuilds(ctx, org, project, params)
	}, http.StatusOK, "list project template builds")
}

func (b *templateBuildClient) CreateForProject(ctx context.Context, params CreateTemplateBuildParams) (*api.TemplateBuildResponse, error) {
	org, project, err := b.projectName(params.Org, params.Project)
	if err != nil {
		return nil, err
	}
	if params.Template == "" {
		return nil, fmt.Errorf("create template build requires template")
	}
	return resourceJSON[api.TemplateBuildResponse](func() (*http.Response, error) {
		return b.raw.CreateProjectTemplateBuild(ctx, org, project, &api.CreateProjectTemplateBuildParams{IdempotencyKey: params.IdempotencyKey}, api.CreateProjectTemplateBuildJSONRequestBody{
			Template: params.Template, Build: params.Build, Sandbox: params.Sandbox, Tag: params.Tag,
		})
	}, http.StatusAccepted, "create project template build")
}

func (b *templateBuildClient) Get(ctx context.Context, buildID string) (*api.TemplateBuildResponse, error) {
	if buildID == "" {
		return nil, fmt.Errorf("get template build requires build ID")
	}
	return resourceJSON[api.TemplateBuildResponse](func() (*http.Response, error) {
		return b.raw.GetTemplateBuild(ctx, buildID)
	}, http.StatusOK, "get template build")
}

func (b *templateBuildClient) Cancel(ctx context.Context, buildID string) (*api.TemplateBuildResponse, error) {
	if buildID == "" {
		return nil, fmt.Errorf("cancel template build requires build ID")
	}
	return resourceJSON[api.TemplateBuildResponse](func() (*http.Response, error) {
		return b.raw.CancelTemplateBuild(ctx, buildID)
	}, http.StatusAccepted, "cancel template build")
}

func (b *templateBuildClient) GetLogs(ctx context.Context, buildID string, params *GetTemplateBuildLogsParams) (*api.TemplateBuildLogsBody, error) {
	if buildID == "" {
		return nil, fmt.Errorf("get template build logs requires build ID")
	}
	return resourceJSON[api.TemplateBuildLogsBody](func() (*http.Response, error) {
		return b.raw.GetTemplateBuildLogs(ctx, buildID, params)
	}, http.StatusOK, "get template build logs")
}

// StreamLogs returns SSE frames. Close the stream when it is no longer needed.
// Cancel ctx to stop a blocked read.
func (b *templateBuildClient) StreamLogs(ctx context.Context, buildID string, params *StreamTemplateBuildLogsParams) (io.ReadCloser, error) {
	if buildID == "" {
		return nil, fmt.Errorf("stream template build logs requires build ID")
	}
	resp, err := b.raw.StreamTemplateBuildLogs(ctx, buildID, params)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		_, _, err = resourceBody(func() (*http.Response, error) { return resp, nil }, http.StatusOK)
		return nil, err
	}
	return resp.Body, nil
}
