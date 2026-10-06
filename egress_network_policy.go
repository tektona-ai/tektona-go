package tektona

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tektona-ai/tektona-go/api"
)

// CreateEgressNetworkPolicyForOrgParams holds scope and policy fields.
type CreateEgressNetworkPolicyForOrgParams struct {
	Org            string
	Name           string
	AllowedCidrs   *[]string
	AllowedDomains *[]string
	DeniedCidrs    *[]string
	DeniedDomains  *[]string
	DisplayName    *string
	Extends        *string
}

// CreateEgressNetworkPolicyForProjectParams holds scope and policy fields.
type CreateEgressNetworkPolicyForProjectParams struct {
	Org            string
	Project        string
	Name           string
	AllowedCidrs   *[]string
	AllowedDomains *[]string
	DeniedCidrs    *[]string
	DeniedDomains  *[]string
	DisplayName    *string
	Extends        *string
}

// GetEgressNetworkPolicyForOrgParams holds an optional organization override.
type GetEgressNetworkPolicyForOrgParams struct{ Org string }

// GetEgressNetworkPolicyForProjectParams holds optional scope overrides.
type GetEgressNetworkPolicyForProjectParams struct{ Org, Project string }

// UpdateEgressNetworkPolicyForOrgParams holds scope and policy fields.
type UpdateEgressNetworkPolicyForOrgParams struct {
	Org            string
	AllowedCidrs   *[]string
	AllowedDomains *[]string
	DeniedCidrs    *[]string
	DeniedDomains  *[]string
	DisplayName    *string
	Extends        *string
}

// UpdateEgressNetworkPolicyForProjectParams holds scope and policy fields.
type UpdateEgressNetworkPolicyForProjectParams struct {
	Org            string
	Project        string
	AllowedCidrs   *[]string
	AllowedDomains *[]string
	DeniedCidrs    *[]string
	DeniedDomains  *[]string
	DisplayName    *string
	Extends        *string
}

// DeleteEgressNetworkPolicyForOrgParams holds an optional organization override.
type DeleteEgressNetworkPolicyForOrgParams struct{ Org string }

// DeleteEgressNetworkPolicyForProjectParams holds optional scope overrides.
type DeleteEgressNetworkPolicyForProjectParams struct{ Org, Project string }

// EgressNetworkPolicy supports system, organization, and project policy operations.
type EgressNetworkPolicy interface {
	ListSystem(ctx context.Context) (*api.ListSystemBody, error)
	GetSystem(ctx context.Context, name string) (*api.Policy, error)
	ListForOrg(ctx context.Context, org string) (*api.ListPoliciesBody, error)
	CreateForOrg(ctx context.Context, params CreateEgressNetworkPolicyForOrgParams) (*api.PolicyResponse, error)
	GetForOrg(ctx context.Context, name string, params *GetEgressNetworkPolicyForOrgParams) (*api.PolicyResponse, error)
	UpdateForOrg(ctx context.Context, name string, params UpdateEgressNetworkPolicyForOrgParams) (*api.PolicyResponse, error)
	DeleteForOrg(ctx context.Context, name string, params *DeleteEgressNetworkPolicyForOrgParams) error
	ListForProject(ctx context.Context, org, project string) (*api.ListPoliciesBody, error)
	CreateForProject(ctx context.Context, params CreateEgressNetworkPolicyForProjectParams) (*api.PolicyResponse, error)
	GetForProject(ctx context.Context, name string, params *GetEgressNetworkPolicyForProjectParams) (*api.PolicyResponse, error)
	UpdateForProject(ctx context.Context, name string, params UpdateEgressNetworkPolicyForProjectParams) (*api.PolicyResponse, error)
	DeleteForProject(ctx context.Context, name string, params *DeleteEgressNetworkPolicyForProjectParams) error
}

type egressNetworkPolicyClient struct {
	raw          *api.ClientWithResponses
	org, project string
}

// EgressNetworkPolicy returns policy methods that share the client's HTTP client and scope defaults.
func (c *Client) EgressNetworkPolicy() EgressNetworkPolicy {
	return &egressNetworkPolicyClient{raw: c.raw, org: c.org, project: c.project}
}

func (p *egressNetworkPolicyClient) address(org, project string) (string, string, error) {
	if org == "" {
		org = p.org
	}
	if project == "" {
		project = p.project
	}
	if org == "" || project == "" {
		return "", "", fmt.Errorf("egress network policy requires org and project")
	}
	return org, project, nil
}

func (p *egressNetworkPolicyClient) organization(org string) (string, error) {
	if org == "" {
		org = p.org
	}
	if org == "" {
		return "", fmt.Errorf("egress network policy requires org")
	}
	return org, nil
}

func (p *egressNetworkPolicyClient) ListSystem(ctx context.Context) (*api.ListSystemBody, error) {
	return resourceJSON[api.ListSystemBody](func() (*http.Response, error) {
		return p.raw.ListSystemEgressNetworkPolicies(ctx)
	}, http.StatusOK, "list system egress network policies")
}

func (p *egressNetworkPolicyClient) GetSystem(ctx context.Context, name string) (*api.Policy, error) {
	if name == "" {
		return nil, fmt.Errorf("get system egress network policy requires name")
	}
	return resourceJSON[api.Policy](func() (*http.Response, error) {
		return p.raw.GetSystemEgressNetworkPolicy(ctx, name)
	}, http.StatusOK, "get system egress network policy")
}

func (p *egressNetworkPolicyClient) ListForOrg(ctx context.Context, org string) (*api.ListPoliciesBody, error) {
	org, err := p.organization(org)
	if err != nil {
		return nil, err
	}
	return resourceJSON[api.ListPoliciesBody](func() (*http.Response, error) {
		return p.raw.ListOrgEgressNetworkPolicies(ctx, org)
	}, http.StatusOK, "list organization egress network policies")
}

func (p *egressNetworkPolicyClient) CreateForOrg(ctx context.Context, params CreateEgressNetworkPolicyForOrgParams) (*api.PolicyResponse, error) {
	org, err := p.organization(params.Org)
	if err != nil {
		return nil, err
	}
	if params.Name == "" {
		return nil, fmt.Errorf("create organization egress network policy requires name")
	}
	return resourceJSON[api.PolicyResponse](func() (*http.Response, error) {
		return p.raw.CreateOrgEgressNetworkPolicy(ctx, org, api.CreateOrgEgressNetworkPolicyJSONRequestBody{
			Name: params.Name, AllowedCidrs: params.AllowedCidrs, AllowedDomains: params.AllowedDomains,
			DeniedCidrs: params.DeniedCidrs, DeniedDomains: params.DeniedDomains,
			DisplayName: params.DisplayName, Extends: params.Extends,
		})
	}, http.StatusCreated, "create organization egress network policy")
}

func (p *egressNetworkPolicyClient) GetForOrg(ctx context.Context, name string, params *GetEgressNetworkPolicyForOrgParams) (*api.PolicyResponse, error) {
	var org string
	if params != nil {
		org = params.Org
	}
	org, err := p.organization(org)
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, fmt.Errorf("get organization egress network policy requires name")
	}
	return resourceJSON[api.PolicyResponse](func() (*http.Response, error) {
		return p.raw.GetOrgEgressNetworkPolicy(ctx, org, name)
	}, http.StatusOK, "get organization egress network policy")
}

func (p *egressNetworkPolicyClient) UpdateForOrg(ctx context.Context, name string, params UpdateEgressNetworkPolicyForOrgParams) (*api.PolicyResponse, error) {
	org, err := p.organization(params.Org)
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, fmt.Errorf("update organization egress network policy requires name")
	}
	return resourceJSON[api.PolicyResponse](func() (*http.Response, error) {
		return p.raw.UpdateOrgEgressNetworkPolicy(ctx, org, name, api.UpdateOrgEgressNetworkPolicyJSONRequestBody{
			AllowedCidrs: params.AllowedCidrs, AllowedDomains: params.AllowedDomains,
			DeniedCidrs: params.DeniedCidrs, DeniedDomains: params.DeniedDomains,
			DisplayName: params.DisplayName, Extends: params.Extends,
		})
	}, http.StatusOK, "update organization egress network policy")
}

func (p *egressNetworkPolicyClient) DeleteForOrg(ctx context.Context, name string, params *DeleteEgressNetworkPolicyForOrgParams) error {
	var org string
	if params != nil {
		org = params.Org
	}
	org, err := p.organization(org)
	if err != nil {
		return err
	}
	if name == "" {
		return fmt.Errorf("delete organization egress network policy requires name")
	}
	return resourceNoContent(func() (*http.Response, error) {
		return p.raw.DeleteOrgEgressNetworkPolicy(ctx, org, name)
	})
}

func (p *egressNetworkPolicyClient) ListForProject(ctx context.Context, org, project string) (*api.ListPoliciesBody, error) {
	org, project, err := p.address(org, project)
	if err != nil {
		return nil, err
	}
	return resourceJSON[api.ListPoliciesBody](func() (*http.Response, error) {
		return p.raw.ListProjectEgressNetworkPolicies(ctx, org, project)
	}, http.StatusOK, "list project egress network policies")
}

func (p *egressNetworkPolicyClient) CreateForProject(ctx context.Context, params CreateEgressNetworkPolicyForProjectParams) (*api.PolicyResponse, error) {
	org, project, err := p.address(params.Org, params.Project)
	if err != nil {
		return nil, err
	}
	if params.Name == "" {
		return nil, fmt.Errorf("create project egress network policy requires name")
	}
	return resourceJSON[api.PolicyResponse](func() (*http.Response, error) {
		return p.raw.CreateProjectEgressNetworkPolicy(ctx, org, project, api.CreateProjectEgressNetworkPolicyJSONRequestBody{
			Name: params.Name, AllowedCidrs: params.AllowedCidrs, AllowedDomains: params.AllowedDomains,
			DeniedCidrs: params.DeniedCidrs, DeniedDomains: params.DeniedDomains,
			DisplayName: params.DisplayName, Extends: params.Extends,
		})
	}, http.StatusCreated, "create project egress network policy")
}

func (p *egressNetworkPolicyClient) GetForProject(ctx context.Context, name string, params *GetEgressNetworkPolicyForProjectParams) (*api.PolicyResponse, error) {
	var org, project string
	if params != nil {
		org, project = params.Org, params.Project
	}
	org, project, err := p.address(org, project)
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, fmt.Errorf("get project egress network policy requires name")
	}
	return resourceJSON[api.PolicyResponse](func() (*http.Response, error) {
		return p.raw.GetProjectEgressNetworkPolicy(ctx, org, project, name)
	}, http.StatusOK, "get project egress network policy")
}

func (p *egressNetworkPolicyClient) UpdateForProject(ctx context.Context, name string, params UpdateEgressNetworkPolicyForProjectParams) (*api.PolicyResponse, error) {
	org, project, err := p.address(params.Org, params.Project)
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, fmt.Errorf("update project egress network policy requires name")
	}
	return resourceJSON[api.PolicyResponse](func() (*http.Response, error) {
		return p.raw.UpdateProjectEgressNetworkPolicy(ctx, org, project, name, api.UpdateProjectEgressNetworkPolicyJSONRequestBody{
			AllowedCidrs: params.AllowedCidrs, AllowedDomains: params.AllowedDomains,
			DeniedCidrs: params.DeniedCidrs, DeniedDomains: params.DeniedDomains,
			DisplayName: params.DisplayName, Extends: params.Extends,
		})
	}, http.StatusOK, "update project egress network policy")
}

func (p *egressNetworkPolicyClient) DeleteForProject(ctx context.Context, name string, params *DeleteEgressNetworkPolicyForProjectParams) error {
	var org, project string
	if params != nil {
		org, project = params.Org, params.Project
	}
	org, project, err := p.address(org, project)
	if err != nil {
		return err
	}
	if name == "" {
		return fmt.Errorf("delete project egress network policy requires name")
	}
	return resourceNoContent(func() (*http.Response, error) {
		return p.raw.DeleteProjectEgressNetworkPolicy(ctx, org, project, name)
	})
}
