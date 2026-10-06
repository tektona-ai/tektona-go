package tektona

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tektona-ai/tektona-go/api"
)

// ListEgressProxyProfilesForOrgParams holds pagination options for one organization.
type ListEgressProxyProfilesForOrgParams = api.ListOrgEgressProxyProfilesParams

// ListEgressProxyProfilesForProjectParams holds the organization override and pagination options.
type ListEgressProxyProfilesForProjectParams struct {
	Org    string
	Cursor *string
	Limit  *int32
}

// CreateEgressProxyProfileForOrgParams holds the organization and profile fields.
type CreateEgressProxyProfileForOrgParams struct {
	Org         string
	Name        string
	DisplayName *string
	Default     bool
}

// CreateEgressProxyProfileForProjectParams holds the project scope and profile fields.
type CreateEgressProxyProfileForProjectParams struct {
	Org         string
	Project     string
	Name        string
	DisplayName *string
	Default     bool
}

// AddEgressProxyProfileRuleForProjectParams holds the project scope and rule fields.
type AddEgressProxyProfileRuleForProjectParams struct {
	Org           string
	Project       string
	DomainPattern string
	PathMatch     *string
	Recipe        *[]api.InjectOp
}

// DeleteEgressProxyProfileForProjectParams holds optional scope overrides.
type DeleteEgressProxyProfileForProjectParams struct{ Org, Project string }

// SetEgressProxyProfileDefaultForProjectParams holds the project scope and default flag.
type SetEgressProxyProfileDefaultForProjectParams struct {
	Org, Project string
	IsDefault    bool
}

// DeleteEgressProxyProfileRuleForProjectParams holds optional scope overrides.
type DeleteEgressProxyProfileRuleForProjectParams struct{ Org, Project string }

// UpdateEgressProxyProfileRuleForProjectParams holds the project scope and rule fields.
type UpdateEgressProxyProfileRuleForProjectParams struct {
	Org           string
	Project       string
	DomainPattern string
	PathMatch     *string
	Recipe        *[]api.InjectOp
}

// EgressProxyProfile supports organization and project proxy profile operations.
type EgressProxyProfile interface {
	ListForOrg(ctx context.Context, org string, params *ListEgressProxyProfilesForOrgParams) (*api.ProfileListBody, error)
	CreateForOrg(ctx context.Context, params CreateEgressProxyProfileForOrgParams) (*api.ProfileCreateBody, error)
	ListForProject(ctx context.Context, project string, params *ListEgressProxyProfilesForProjectParams) (*api.ProfileListBody, error)
	CreateForProject(ctx context.Context, params CreateEgressProxyProfileForProjectParams) (*api.ProfileCreateBody, error)
	AddRuleForProject(ctx context.Context, name string, params AddEgressProxyProfileRuleForProjectParams) (*api.RuleCreatedBody, error)
	DeleteForProject(ctx context.Context, profileID string, params *DeleteEgressProxyProfileForProjectParams) error
	SetDefaultForProject(ctx context.Context, profileID string, params SetEgressProxyProfileDefaultForProjectParams) error
	DeleteRuleForProject(ctx context.Context, profileID, ruleID string, params *DeleteEgressProxyProfileRuleForProjectParams) error
	UpdateRuleForProject(ctx context.Context, profileID, ruleID string, params UpdateEgressProxyProfileRuleForProjectParams) error
}

type egressProxyProfileClient struct {
	raw          *api.ClientWithResponses
	org, project string
}

// EgressProxyProfile returns the proxy profile methods.
func (c *Client) EgressProxyProfile() EgressProxyProfile {
	return &egressProxyProfileClient{raw: c.raw, org: c.org, project: c.project}
}

func (p *egressProxyProfileClient) organization(org string) (string, error) {
	if org == "" {
		org = p.org
	}
	if org == "" {
		return "", fmt.Errorf("egress proxy profile requires org")
	}
	return org, nil
}

func (p *egressProxyProfileClient) scope(org, project string) (string, string, error) {
	org, err := p.organization(org)
	if project == "" {
		project = p.project
	}
	if err != nil || project == "" {
		return "", "", fmt.Errorf("egress proxy profile requires org and project")
	}
	return org, project, nil
}

func (p *egressProxyProfileClient) ListForOrg(ctx context.Context, org string, params *ListEgressProxyProfilesForOrgParams) (*api.ProfileListBody, error) {
	org, err := p.organization(org)
	if err != nil {
		return nil, err
	}
	return resourceJSON[api.ProfileListBody](func() (*http.Response, error) {
		return p.raw.ListOrgEgressProxyProfiles(ctx, org, params)
	}, http.StatusOK, "list organization egress proxy profiles")
}

func (p *egressProxyProfileClient) CreateForOrg(ctx context.Context, params CreateEgressProxyProfileForOrgParams) (*api.ProfileCreateBody, error) {
	org, err := p.organization(params.Org)
	if err != nil || params.Name == "" {
		return nil, fmt.Errorf("create organization egress proxy profile requires org and name")
	}
	return resourceJSON[api.ProfileCreateBody](func() (*http.Response, error) {
		return p.raw.CreateOrgEgressProxyProfile(ctx, org, api.CreateProfileBody{
			Name: params.Name, DisplayName: params.DisplayName, Default: params.Default, Scope: api.CreateProfileBodyScope("org"),
		})
	}, http.StatusCreated, "create organization egress proxy profile")
}

func (p *egressProxyProfileClient) ListForProject(ctx context.Context, project string, params *ListEgressProxyProfilesForProjectParams) (*api.ProfileListBody, error) {
	var org string
	var query *api.ListOrgProjectEgressProxyProfilesParams
	if params != nil {
		org = params.Org
		query = &api.ListOrgProjectEgressProxyProfilesParams{Cursor: params.Cursor, Limit: params.Limit}
	}
	org, project, err := p.scope(org, project)
	if err != nil {
		return nil, err
	}
	return resourceJSON[api.ProfileListBody](func() (*http.Response, error) {
		return p.raw.ListOrgProjectEgressProxyProfiles(ctx, org, project, query)
	}, http.StatusOK, "list project egress proxy profiles")
}

func (p *egressProxyProfileClient) CreateForProject(ctx context.Context, params CreateEgressProxyProfileForProjectParams) (*api.ProfileCreateBody, error) {
	org, project, err := p.scope(params.Org, params.Project)
	if err != nil || params.Name == "" {
		return nil, fmt.Errorf("create project egress proxy profile requires org, project, and name")
	}
	return resourceJSON[api.ProfileCreateBody](func() (*http.Response, error) {
		return p.raw.CreateOrgProjectEgressProxyProfile(ctx, org, project, api.CreateProfileBody{
			Name: params.Name, DisplayName: params.DisplayName, Default: params.Default, Scope: api.CreateProfileBodyScope("project"),
		})
	}, http.StatusCreated, "create project egress proxy profile")
}

func (p *egressProxyProfileClient) AddRuleForProject(ctx context.Context, name string, params AddEgressProxyProfileRuleForProjectParams) (*api.RuleCreatedBody, error) {
	org, project, err := p.scope(params.Org, params.Project)
	if err != nil || name == "" {
		return nil, fmt.Errorf("add egress proxy profile rule requires org, project, and name")
	}
	return resourceJSON[api.RuleCreatedBody](func() (*http.Response, error) {
		return p.raw.AddOrgProjectEgressProxyProfileRule(ctx, org, project, name, api.AddRuleBody{
			DomainPattern: params.DomainPattern, PathMatch: params.PathMatch, Recipe: params.Recipe,
		})
	}, http.StatusCreated, "add egress proxy profile rule")
}

func (p *egressProxyProfileClient) DeleteForProject(ctx context.Context, profileID string, params *DeleteEgressProxyProfileForProjectParams) error {
	var org, project string
	if params != nil {
		org, project = params.Org, params.Project
	}
	org, project, err := p.scope(org, project)
	if err != nil || profileID == "" {
		return fmt.Errorf("delete egress proxy profile requires org, project, and profile ID")
	}
	return resourceNoContent(func() (*http.Response, error) {
		return p.raw.DeleteOrgProjectEgressProxyProfile(ctx, org, project, profileID)
	})
}

func (p *egressProxyProfileClient) SetDefaultForProject(ctx context.Context, profileID string, params SetEgressProxyProfileDefaultForProjectParams) error {
	org, project, err := p.scope(params.Org, params.Project)
	if err != nil || profileID == "" {
		return fmt.Errorf("set egress proxy profile default requires org, project, and profile ID")
	}
	return resourceNoContent(func() (*http.Response, error) {
		return p.raw.SetOrgProjectEgressProxyProfileDefault(ctx, org, project, profileID, api.SetDefaultBody{IsDefault: params.IsDefault})
	})
}

func (p *egressProxyProfileClient) DeleteRuleForProject(ctx context.Context, profileID, ruleID string, params *DeleteEgressProxyProfileRuleForProjectParams) error {
	var org, project string
	if params != nil {
		org, project = params.Org, params.Project
	}
	org, project, err := p.scope(org, project)
	if err != nil || profileID == "" || ruleID == "" {
		return fmt.Errorf("delete egress proxy profile rule requires org, project, profile ID, and rule ID")
	}
	return resourceNoContent(func() (*http.Response, error) {
		return p.raw.DeleteOrgProjectEgressProxyProfileRule(ctx, org, project, profileID, ruleID)
	})
}

func (p *egressProxyProfileClient) UpdateRuleForProject(ctx context.Context, profileID, ruleID string, params UpdateEgressProxyProfileRuleForProjectParams) error {
	org, project, err := p.scope(params.Org, params.Project)
	if err != nil || profileID == "" || ruleID == "" {
		return fmt.Errorf("update egress proxy profile rule requires org, project, profile ID, and rule ID")
	}
	return resourceNoContent(func() (*http.Response, error) {
		return p.raw.UpdateOrgProjectEgressProxyProfileRule(ctx, org, project, profileID, ruleID, api.AddRuleBody{
			DomainPattern: params.DomainPattern, PathMatch: params.PathMatch, Recipe: params.Recipe,
		})
	})
}
