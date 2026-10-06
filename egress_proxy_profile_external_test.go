package tektona_test

import (
	"context"

	tektona "github.com/tektona-ai/tektona-go"
	"github.com/tektona-ai/tektona-go/api"
)

// The fake compiles outside the SDK package without a network connection.
type fakeEgressProxyProfile struct{}

var _ tektona.EgressProxyProfile = fakeEgressProxyProfile{}

func (fakeEgressProxyProfile) ListForOrg(context.Context, string, *tektona.ListEgressProxyProfilesForOrgParams) (*api.ProfileListBody, error) {
	return nil, nil
}
func (fakeEgressProxyProfile) CreateForOrg(context.Context, tektona.CreateEgressProxyProfileForOrgParams) (*api.ProfileCreateBody, error) {
	return nil, nil
}
func (fakeEgressProxyProfile) ListForProject(context.Context, string, *tektona.ListEgressProxyProfilesForProjectParams) (*api.ProfileListBody, error) {
	return nil, nil
}
func (fakeEgressProxyProfile) CreateForProject(context.Context, tektona.CreateEgressProxyProfileForProjectParams) (*api.ProfileCreateBody, error) {
	return nil, nil
}
func (fakeEgressProxyProfile) AddRuleForProject(context.Context, string, tektona.AddEgressProxyProfileRuleForProjectParams) (*api.RuleCreatedBody, error) {
	return nil, nil
}
func (fakeEgressProxyProfile) DeleteForProject(context.Context, string, *tektona.DeleteEgressProxyProfileForProjectParams) error {
	return nil
}
func (fakeEgressProxyProfile) SetDefaultForProject(context.Context, string, tektona.SetEgressProxyProfileDefaultForProjectParams) error {
	return nil
}
func (fakeEgressProxyProfile) DeleteRuleForProject(context.Context, string, string, *tektona.DeleteEgressProxyProfileRuleForProjectParams) error {
	return nil
}
func (fakeEgressProxyProfile) UpdateRuleForProject(context.Context, string, string, tektona.UpdateEgressProxyProfileRuleForProjectParams) error {
	return nil
}
