package tektona_test

import (
	"context"
	"io"
	"strings"

	tektona "github.com/tektona-ai/tektona-go"
	"github.com/tektona-ai/tektona-go/api"
)

type fakeTemplateBuild struct{}

func (fakeTemplateBuild) ListForOrg(context.Context, string, *tektona.ListTemplateBuildsForOrgParams) (*api.ListTemplateBuildsBody, error) {
	return &api.ListTemplateBuildsBody{}, nil
}
func (fakeTemplateBuild) CreateForOrg(context.Context, tektona.CreateTemplateBuildParams) (*api.TemplateBuildResponse, error) {
	return &api.TemplateBuildResponse{}, nil
}
func (fakeTemplateBuild) ListForProject(context.Context, string, string, *tektona.ListTemplateBuildsForProjectParams) (*api.ListTemplateBuildsBody, error) {
	return &api.ListTemplateBuildsBody{}, nil
}
func (fakeTemplateBuild) CreateForProject(context.Context, tektona.CreateTemplateBuildParams) (*api.TemplateBuildResponse, error) {
	return &api.TemplateBuildResponse{}, nil
}
func (fakeTemplateBuild) Get(context.Context, string) (*api.TemplateBuildResponse, error) {
	return &api.TemplateBuildResponse{}, nil
}
func (fakeTemplateBuild) Cancel(context.Context, string) (*api.TemplateBuildResponse, error) {
	return &api.TemplateBuildResponse{}, nil
}
func (fakeTemplateBuild) GetLogs(context.Context, string, *tektona.GetTemplateBuildLogsParams) (*api.TemplateBuildLogsBody, error) {
	return &api.TemplateBuildLogsBody{}, nil
}
func (fakeTemplateBuild) StreamLogs(context.Context, string, *tektona.StreamTemplateBuildLogsParams) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader("")), nil
}

var _ tektona.TemplateBuild = fakeTemplateBuild{}
