package tektona_test

import (
	"context"

	tektona "github.com/tektona-ai/tektona-go"
	"github.com/tektona-ai/tektona-go/api"
)

type fakeTemplate struct{}

func (fakeTemplate) ListForOrg(context.Context, string, *tektona.ListTemplatesForOrgParams) (*api.ListTemplatesBody, error) {
	return &api.ListTemplatesBody{}, nil
}
func (fakeTemplate) CreateForOrg(context.Context, tektona.CreateTemplateForOrgParams) (*api.TemplateResponse, error) {
	return &api.TemplateResponse{}, nil
}
func (fakeTemplate) GetForOrg(context.Context, string, *tektona.GetTemplateForOrgParams) (*api.TemplateResponse, error) {
	return &api.TemplateResponse{}, nil
}
func (fakeTemplate) UpdateForOrg(context.Context, string, tektona.UpdateTemplateForOrgParams) (*api.TemplateResponse, error) {
	return &api.TemplateResponse{}, nil
}
func (fakeTemplate) DeleteForOrg(context.Context, string, *tektona.DeleteTemplateForOrgParams) (*api.DeleteTemplateResponse, error) {
	return &api.DeleteTemplateResponse{}, nil
}
func (fakeTemplate) ActivateForOrg(context.Context, string, *tektona.ActivateTemplateForOrgParams) (*api.TemplateResponse, error) {
	return &api.TemplateResponse{}, nil
}
func (fakeTemplate) ArchiveForOrg(context.Context, string, *tektona.ArchiveTemplateForOrgParams) (*api.TemplateResponse, error) {
	return &api.TemplateResponse{}, nil
}
func (fakeTemplate) ListForProject(context.Context, string, string, *tektona.ListTemplatesForProjectParams) (*api.ListTemplatesBody, error) {
	return &api.ListTemplatesBody{}, nil
}
func (fakeTemplate) CreateForProject(context.Context, tektona.CreateTemplateForProjectParams) (*api.TemplateResponse, error) {
	return &api.TemplateResponse{}, nil
}
func (fakeTemplate) GetForProject(context.Context, string, *tektona.GetTemplateForProjectParams) (*api.TemplateResponse, error) {
	return &api.TemplateResponse{}, nil
}
func (fakeTemplate) UpdateForProject(context.Context, string, tektona.UpdateTemplateForProjectParams) (*api.TemplateResponse, error) {
	return &api.TemplateResponse{}, nil
}
func (fakeTemplate) DeleteForProject(context.Context, string, *tektona.DeleteTemplateForProjectParams) (*api.DeleteTemplateResponse, error) {
	return &api.DeleteTemplateResponse{}, nil
}
func (fakeTemplate) ActivateForProject(context.Context, string, *tektona.ActivateTemplateForProjectParams) (*api.TemplateResponse, error) {
	return &api.TemplateResponse{}, nil
}
func (fakeTemplate) ArchiveForProject(context.Context, string, *tektona.ArchiveTemplateForProjectParams) (*api.TemplateResponse, error) {
	return &api.TemplateResponse{}, nil
}
func (fakeTemplate) GetSystem(context.Context, string) (*api.TemplateResponse, error) {
	return &api.TemplateResponse{}, nil
}

var _ tektona.Template = fakeTemplate{}
