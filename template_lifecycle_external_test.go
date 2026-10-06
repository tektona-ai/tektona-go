package tektona_test

import (
	"context"

	tektona "github.com/tektona-ai/tektona-go"
	"github.com/tektona-ai/tektona-go/api"
)

type fakeTemplateLifecycle struct{}

func (fakeTemplateLifecycle) GetForOrg(context.Context, string) (*api.TemplateLifecycleSettingsResponse, error) {
	return &api.TemplateLifecycleSettingsResponse{}, nil
}
func (fakeTemplateLifecycle) UpdateForOrg(context.Context, string, tektona.UpdateTemplateLifecycleForOrgParams) (*api.TemplateLifecycleSettingsResponse, error) {
	return &api.TemplateLifecycleSettingsResponse{}, nil
}
func (fakeTemplateLifecycle) PreviewForOrg(context.Context, string, tektona.PreviewTemplateLifecycleForOrgParams) (*api.TemplateLifecyclePreviewResponse, error) {
	return &api.TemplateLifecyclePreviewResponse{}, nil
}
func (fakeTemplateLifecycle) GetForProject(context.Context, string, *tektona.GetTemplateLifecycleForProjectParams) (*api.TemplateLifecycleSettingsResponse, error) {
	return &api.TemplateLifecycleSettingsResponse{}, nil
}
func (fakeTemplateLifecycle) UpdateForProject(context.Context, string, tektona.UpdateTemplateLifecycleForProjectParams) (*api.TemplateLifecycleSettingsResponse, error) {
	return &api.TemplateLifecycleSettingsResponse{}, nil
}
func (fakeTemplateLifecycle) PreviewForProject(context.Context, string, tektona.PreviewTemplateLifecycleForProjectParams) (*api.TemplateLifecyclePreviewResponse, error) {
	return &api.TemplateLifecyclePreviewResponse{}, nil
}

var _ tektona.TemplateLifecycle = fakeTemplateLifecycle{}
