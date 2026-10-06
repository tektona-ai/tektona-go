package tektona_test

import (
	"context"

	tektona "github.com/tektona-ai/tektona-go"
	"github.com/tektona-ai/tektona-go/api"
)

type fakeTemplateVersion struct{}

func (fakeTemplateVersion) ListForOrg(context.Context, string, *tektona.ListTemplateVersionsForOrgParams) (*api.ListTemplateVersionsBody, error) {
	return &api.ListTemplateVersionsBody{}, nil
}
func (fakeTemplateVersion) GetForOrg(context.Context, string, string, *tektona.TemplateVersionForOrgParams) (*api.TemplateVersion, error) {
	return &api.TemplateVersion{}, nil
}
func (fakeTemplateVersion) DeleteForOrg(context.Context, string, string, *tektona.TemplateVersionForOrgParams) error {
	return nil
}
func (fakeTemplateVersion) PruneForOrg(context.Context, string, tektona.PruneTemplateVersionsForOrgParams) (*api.PruneVersionsReport, error) {
	return &api.PruneVersionsReport{}, nil
}
func (fakeTemplateVersion) ActivateForOrg(context.Context, string, string, *tektona.TemplateVersionForOrgParams) (*api.TemplateVersion, error) {
	return &api.TemplateVersion{}, nil
}
func (fakeTemplateVersion) ArchiveForOrg(context.Context, string, string, *tektona.TemplateVersionForOrgParams) (*api.TemplateVersion, error) {
	return &api.TemplateVersion{}, nil
}
func (fakeTemplateVersion) ListForProject(context.Context, string, *tektona.ListTemplateVersionsForProjectParams) (*api.ListTemplateVersionsBody, error) {
	return &api.ListTemplateVersionsBody{}, nil
}
func (fakeTemplateVersion) GetForProject(context.Context, string, string, *tektona.TemplateVersionForProjectParams) (*api.TemplateVersion, error) {
	return &api.TemplateVersion{}, nil
}
func (fakeTemplateVersion) DeleteForProject(context.Context, string, string, *tektona.TemplateVersionForProjectParams) error {
	return nil
}
func (fakeTemplateVersion) PruneForProject(context.Context, string, tektona.PruneTemplateVersionsForProjectParams) (*api.PruneVersionsReport, error) {
	return &api.PruneVersionsReport{}, nil
}
func (fakeTemplateVersion) ActivateForProject(context.Context, string, string, *tektona.TemplateVersionForProjectParams) (*api.TemplateVersion, error) {
	return &api.TemplateVersion{}, nil
}
func (fakeTemplateVersion) ArchiveForProject(context.Context, string, string, *tektona.TemplateVersionForProjectParams) (*api.TemplateVersion, error) {
	return &api.TemplateVersion{}, nil
}
func (fakeTemplateVersion) ListSystem(context.Context, string, *tektona.ListTemplateVersionsSystemParams) (*api.ListTemplateVersionsBody, error) {
	return &api.ListTemplateVersionsBody{}, nil
}
func (fakeTemplateVersion) GetSystem(context.Context, string, string) (*api.TemplateVersion, error) {
	return &api.TemplateVersion{}, nil
}

var _ tektona.TemplateVersion = fakeTemplateVersion{}
