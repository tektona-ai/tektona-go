package tektona_test

import (
	"context"

	tektona "github.com/tektona-ai/tektona-go"
	"github.com/tektona-ai/tektona-go/api"
)

type fakeEgressNetworkPolicy struct{}

func (fakeEgressNetworkPolicy) ListSystem(context.Context) (*api.ListSystemBody, error) {
	return &api.ListSystemBody{}, nil
}
func (fakeEgressNetworkPolicy) GetSystem(context.Context, string) (*api.Policy, error) {
	return &api.Policy{}, nil
}
func (fakeEgressNetworkPolicy) ListForOrg(context.Context, string) (*api.ListPoliciesBody, error) {
	return &api.ListPoliciesBody{}, nil
}
func (fakeEgressNetworkPolicy) CreateForOrg(context.Context, tektona.CreateEgressNetworkPolicyForOrgParams) (*api.PolicyResponse, error) {
	return &api.PolicyResponse{}, nil
}
func (fakeEgressNetworkPolicy) GetForOrg(context.Context, string, *tektona.GetEgressNetworkPolicyForOrgParams) (*api.PolicyResponse, error) {
	return &api.PolicyResponse{}, nil
}
func (fakeEgressNetworkPolicy) UpdateForOrg(context.Context, string, tektona.UpdateEgressNetworkPolicyForOrgParams) (*api.PolicyResponse, error) {
	return &api.PolicyResponse{}, nil
}
func (fakeEgressNetworkPolicy) DeleteForOrg(context.Context, string, *tektona.DeleteEgressNetworkPolicyForOrgParams) error {
	return nil
}
func (fakeEgressNetworkPolicy) ListForProject(context.Context, string, string) (*api.ListPoliciesBody, error) {
	return &api.ListPoliciesBody{}, nil
}
func (fakeEgressNetworkPolicy) CreateForProject(context.Context, tektona.CreateEgressNetworkPolicyForProjectParams) (*api.PolicyResponse, error) {
	return &api.PolicyResponse{}, nil
}
func (fakeEgressNetworkPolicy) GetForProject(context.Context, string, *tektona.GetEgressNetworkPolicyForProjectParams) (*api.PolicyResponse, error) {
	return &api.PolicyResponse{}, nil
}
func (fakeEgressNetworkPolicy) UpdateForProject(context.Context, string, tektona.UpdateEgressNetworkPolicyForProjectParams) (*api.PolicyResponse, error) {
	return &api.PolicyResponse{}, nil
}
func (fakeEgressNetworkPolicy) DeleteForProject(context.Context, string, *tektona.DeleteEgressNetworkPolicyForProjectParams) error {
	return nil
}

var _ tektona.EgressNetworkPolicy = fakeEgressNetworkPolicy{}
