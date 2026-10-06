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
	resp, err := o.raw.ListOrgsWithResponse(ctx, params)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, responseError(resp.StatusCode(), resp.HTTPResponse.Header, resp.Body, resp.ApplicationproblemJSONDefault)
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("list organizations: HTTP 200 lacks a JSON response")
	}
	return resp.JSON200, nil
}

func (o *organizationClient) Create(ctx context.Context, params CreateOrganizationParams) (*api.OrgResponse, error) {
	if params.Name == "" || params.DisplayName == "" {
		return nil, fmt.Errorf("create organization requires name and display name")
	}
	resp, err := o.raw.CreateOrgWithResponse(ctx, params)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusCreated {
		return nil, responseError(resp.StatusCode(), resp.HTTPResponse.Header, resp.Body, resp.ApplicationproblemJSONDefault)
	}
	if resp.JSON201 == nil {
		return nil, fmt.Errorf("create organization: HTTP 201 lacks a JSON response")
	}
	return resp.JSON201, nil
}

func (o *organizationClient) Get(ctx context.Context, org string) (*api.OrgResponse, error) {
	if org == "" {
		org = o.org
	}
	if org == "" {
		return nil, fmt.Errorf("get organization requires org")
	}
	resp, err := o.raw.GetOrgWithResponse(ctx, org)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, responseError(resp.StatusCode(), resp.HTTPResponse.Header, resp.Body, resp.ApplicationproblemJSONDefault)
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("get organization: HTTP 200 lacks a JSON response")
	}
	return resp.JSON200, nil
}

func (o *organizationClient) Update(ctx context.Context, org string, params UpdateOrganizationParams) (*api.OrgResponse, error) {
	if org == "" {
		org = o.org
	}
	if org == "" || params.DisplayName == "" {
		return nil, fmt.Errorf("update organization requires org and display name")
	}
	resp, err := o.raw.UpdateOrgWithResponse(ctx, org, params)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, responseError(resp.StatusCode(), resp.HTTPResponse.Header, resp.Body, resp.ApplicationproblemJSONDefault)
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("update organization: HTTP 200 lacks a JSON response")
	}
	return resp.JSON200, nil
}

func (o *organizationClient) ListMembers(ctx context.Context, org string, params *ListOrganizationMembersParams) (*api.ListMembersBody, error) {
	if org == "" {
		org = o.org
	}
	if org == "" {
		return nil, fmt.Errorf("list organization members requires org")
	}
	resp, err := o.raw.ListOrgMembersWithResponse(ctx, org, params)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, responseError(resp.StatusCode(), resp.HTTPResponse.Header, resp.Body, resp.ApplicationproblemJSONDefault)
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("list organization members: HTTP 200 lacks a JSON response")
	}
	return resp.JSON200, nil
}
