package tektona

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/tektona-ai/tektona-go/api"
)

// AutoPauseMode names an auto-pause mode from the public API.
type AutoPauseMode = api.CreateSandboxRequestBodyAutoPauseMode

// SandboxResources sets CPU, memory, and disk limits.
type SandboxResources = api.PublicResources

// CreateSandboxParams holds sandbox creation options without HTTP body or query details.
type CreateSandboxParams struct {
	Org                 string
	Project             string
	Template            string
	Wait                bool
	Name                *string
	Location            *string
	Resources           *SandboxResources
	Env                 *map[string]string
	User                *string
	Workdir             *string
	Tags                *[]string
	Public              *bool
	EgressNetworkPolicy *string
	EgressProxyProfile  *string
	AutoPauseAfter      *string
	AutoPauseMode       *AutoPauseMode
	AutoResume          *bool
	AutoDeleteAfter     *string
}

// GetSandboxParams holds optional scope for a sandbox name.
type GetSandboxParams = api.GetSandboxParams

// ListSandboxesParams holds sandbox filters and pagination options.
type ListSandboxesParams = api.ListSandboxesParams

// DeleteSandboxParams holds optional scope for a sandbox name.
type DeleteSandboxParams = api.DeleteSandboxParams

// Sandbox supports the first four sandbox operations. Callers can replace it with a fake.
type Sandbox interface {
	Create(ctx context.Context, params CreateSandboxParams) (*api.CreateSandboxResponseBody, error)
	Get(ctx context.Context, idOrName string, params *GetSandboxParams) (*api.SandboxObject, error)
	List(ctx context.Context, params *ListSandboxesParams) (*api.ListSandboxesBody, error)
	Delete(ctx context.Context, idOrName string, params *DeleteSandboxParams) error
}

type sandboxClient struct {
	raw          *api.ClientWithResponses
	org, project string
}

func (s *sandboxClient) Create(ctx context.Context, params CreateSandboxParams) (*api.CreateSandboxResponseBody, error) {
	org, project := params.Org, params.Project
	if org == "" {
		org = s.org
	}
	if project == "" {
		project = s.project
	}
	if org == "" || project == "" || params.Template == "" {
		return nil, fmt.Errorf("create sandbox requires org, project, and template")
	}
	var query *api.CreateSandboxParams
	if params.Wait {
		query = &api.CreateSandboxParams{Wait: &params.Wait}
	}
	return resourceJSON[api.CreateSandboxResponseBody](func() (*http.Response, error) {
		return s.raw.CreateSandbox(ctx, query, api.CreateSandboxJSONRequestBody{
			Org: org, Project: project, Template: params.Template, Name: params.Name,
			Location: params.Location, Resources: params.Resources, Env: params.Env,
			User: params.User, Workdir: params.Workdir, Tags: params.Tags, Public: params.Public,
			EgressNetworkPolicy: params.EgressNetworkPolicy, EgressProxyProfile: params.EgressProxyProfile,
			AutoPauseAfter: params.AutoPauseAfter, AutoPauseMode: params.AutoPauseMode,
			AutoResume: params.AutoResume, AutoDeleteAfter: params.AutoDeleteAfter,
		})
	}, http.StatusCreated, "create sandbox")
}

func (s *sandboxClient) List(ctx context.Context, params *ListSandboxesParams) (*api.ListSandboxesBody, error) {
	var p ListSandboxesParams
	if params != nil {
		p = *params
	}
	if p.Org == "" {
		p.Org = s.org
	}
	if p.Org == "" {
		return nil, fmt.Errorf("list sandboxes requires org")
	}
	if p.Project == nil && s.project != "" {
		p.Project = &s.project
	}
	return resourceJSON[api.ListSandboxesBody](func() (*http.Response, error) {
		return s.raw.ListSandboxes(ctx, &p)
	}, http.StatusOK, "list sandboxes")
}

func (s *sandboxClient) Get(ctx context.Context, idOrName string, params *GetSandboxParams) (*api.SandboxObject, error) {
	if idOrName == "" {
		return nil, fmt.Errorf("get sandbox requires an ID or name")
	}
	var p GetSandboxParams
	if params != nil {
		p = *params
	}
	if err := s.address(idOrName, &p.Org, &p.Project); err != nil {
		return nil, err
	}
	return resourceJSON[api.SandboxObject](func() (*http.Response, error) {
		return s.raw.GetSandbox(ctx, idOrName, &p)
	}, http.StatusOK, "get sandbox")
}

func (s *sandboxClient) Delete(ctx context.Context, idOrName string, params *DeleteSandboxParams) error {
	if idOrName == "" {
		return fmt.Errorf("delete sandbox requires an ID or name")
	}
	var p DeleteSandboxParams
	if params != nil {
		p = *params
	}
	if err := s.address(idOrName, &p.Org, &p.Project); err != nil {
		return err
	}
	return resourceNoContent(func() (*http.Response, error) {
		return s.raw.DeleteSandbox(ctx, idOrName, &p)
	})
}

func (s *sandboxClient) address(value string, org, project **string) error {
	if isSandboxID(value) {
		return nil
	}
	if *org == nil {
		*org = &s.org
	}
	if *project == nil {
		*project = &s.project
	}
	if **org == "" || **project == "" {
		return fmt.Errorf("sandbox name requires org and project")
	}
	return nil
}

func isSandboxID(value string) bool {
	if len(value) != 26 {
		return false
	}
	const alphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	for _, ch := range strings.ToUpper(value) {
		if !strings.ContainsRune(alphabet, ch) {
			return false
		}
	}
	return true
}
