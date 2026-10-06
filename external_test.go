package tektona_test

import (
	"context"

	tektona "github.com/tektona-ai/tektona-go"
	"github.com/tektona-ai/tektona-go/api"
)

type fakeLocation struct{}

func (fakeLocation) List(context.Context) (*api.ListLocationsOutputBody, error) {
	return &api.ListLocationsOutputBody{}, nil
}

type fakeMeta struct{}

func (fakeMeta) Get(context.Context) (*api.MetaResponse, error) { return &api.MetaResponse{}, nil }

type fakeOrganization struct{}

func (fakeOrganization) List(context.Context, *tektona.ListOrganizationsParams) (*api.ListOrgsBody, error) {
	return &api.ListOrgsBody{}, nil
}
func (fakeOrganization) Create(context.Context, tektona.CreateOrganizationParams) (*api.OrgResponse, error) {
	return &api.OrgResponse{}, nil
}
func (fakeOrganization) Get(context.Context, string) (*api.OrgResponse, error) {
	return &api.OrgResponse{}, nil
}
func (fakeOrganization) Update(context.Context, string, tektona.UpdateOrganizationParams) (*api.OrgResponse, error) {
	return &api.OrgResponse{}, nil
}
func (fakeOrganization) ListMembers(context.Context, string, *tektona.ListOrganizationMembersParams) (*api.ListMembersBody, error) {
	return &api.ListMembersBody{}, nil
}

type fakeProject struct{}

func (fakeProject) List(context.Context, *tektona.ListProjectsParams) (*api.ListAccessibleProjectsBody, error) {
	return &api.ListAccessibleProjectsBody{}, nil
}
func (fakeProject) ListForOrg(context.Context, string, *tektona.ListProjectsForOrgParams) (*api.ListProjectsBody, error) {
	return &api.ListProjectsBody{}, nil
}
func (fakeProject) Create(context.Context, tektona.CreateProjectParams) (*api.ProjectResponse, error) {
	return &api.ProjectResponse{}, nil
}
func (fakeProject) Get(context.Context, string, *tektona.GetProjectParams) (*api.ProjectResponse, error) {
	return &api.ProjectResponse{}, nil
}
func (fakeProject) Update(context.Context, string, tektona.UpdateProjectParams) (*api.ProjectResponse, error) {
	return &api.ProjectResponse{}, nil
}
func (fakeProject) Delete(context.Context, string, *tektona.DeleteProjectParams) error { return nil }
func (fakeProject) GetLifecycleDefaults(context.Context, string, *tektona.GetProjectLifecycleDefaultsParams) (*api.LifecycleDefaultsBody, error) {
	return &api.LifecycleDefaultsBody{}, nil
}
func (fakeProject) UpdateLifecycleDefaults(context.Context, string, tektona.UpdateProjectLifecycleDefaultsParams) (*api.LifecycleDefaultsBody, error) {
	return &api.LifecycleDefaultsBody{}, nil
}

var (
	_ tektona.Location     = fakeLocation{}
	_ tektona.Meta         = fakeMeta{}
	_ tektona.Organization = fakeOrganization{}
	_ tektona.Project      = fakeProject{}
)
