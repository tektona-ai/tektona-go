package tektona

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tektona-ai/tektona-go/api"
)

// ListOrganizationsParams holds pagination options.
type ListOrganizationsParams = api.ListOrgsParams

// CreateOrganizationParams holds organization fields.
type CreateOrganizationParams = api.CreateOrgJSONRequestBody

// UpdateOrganizationParams holds organization fields.
type UpdateOrganizationParams = api.UpdateOrgJSONRequestBody

// ListOrganizationMembersParams holds member filters and pagination options.
type ListOrganizationMembersParams = api.ListOrgMembersParams

// Organization supports organization and member operations.
type Organization interface {
	List(ctx context.Context, params *ListOrganizationsParams) (*api.ListOrgsBody, error)
	Create(ctx context.Context, params CreateOrganizationParams) (*api.OrgResponse, error)
	Get(ctx context.Context, org string) (*api.OrgResponse, error)
	Update(ctx context.Context, org string, params UpdateOrganizationParams) (*api.OrgResponse, error)
	ListMembers(ctx context.Context, org string, params *ListOrganizationMembersParams) (*api.ListMembersBody, error)
}

type organizationClient struct {
	raw *api.ClientWithResponses
	org string
}

// Organization returns the organization methods.
func (c *Client) Organization() Organization { return c.organization }

func (o *organizationClient) List(ctx context.Context, params *ListOrganizationsParams) (*api.ListOrgsBody, error) {
	return resourceJSON[api.ListOrgsBody](func() (*http.Response, error) {
		return o.raw.ListOrgs(ctx, params)
	}, http.StatusOK, "list organizations")
}

func (o *organizationClient) Create(ctx context.Context, params CreateOrganizationParams) (*api.OrgResponse, error) {
	if params.Name == "" || params.DisplayName == "" {
		return nil, fmt.Errorf("create organization requires name and display name")
	}
	return resourceJSON[api.OrgResponse](func() (*http.Response, error) {
		return o.raw.CreateOrg(ctx, params)
	}, http.StatusCreated, "create organization")
}

func (o *organizationClient) Get(ctx context.Context, org string) (*api.OrgResponse, error) {
	if org == "" {
		return nil, fmt.Errorf("get organization requires org")
	}
	return resourceJSON[api.OrgResponse](func() (*http.Response, error) {
		return o.raw.GetOrg(ctx, org)
	}, http.StatusOK, "get organization")
}

func (o *organizationClient) Update(ctx context.Context, org string, params UpdateOrganizationParams) (*api.OrgResponse, error) {
	if org == "" || params.DisplayName == "" {
		return nil, fmt.Errorf("update organization requires org and display name")
	}
	return resourceJSON[api.OrgResponse](func() (*http.Response, error) {
		return o.raw.UpdateOrg(ctx, org, params)
	}, http.StatusOK, "update organization")
}

func (o *organizationClient) ListMembers(ctx context.Context, org string, params *ListOrganizationMembersParams) (*api.ListMembersBody, error) {
	if org == "" {
		org = o.org
	}
	if org == "" {
		return nil, fmt.Errorf("list organization members requires org")
	}
	return resourceJSON[api.ListMembersBody](func() (*http.Response, error) {
		return o.raw.ListOrgMembers(ctx, org, params)
	}, http.StatusOK, "list organization members")
}
