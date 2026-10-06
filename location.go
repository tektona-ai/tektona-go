package tektona

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tektona-ai/tektona-go/api"
)

// Location lists available sandbox locations.
type Location interface {
	List(ctx context.Context) (*api.ListLocationsOutputBody, error)
}

type locationClient struct{ raw *api.ClientWithResponses }

// Location returns the location methods.
func (c *Client) Location() Location { return c.location }

func (l *locationClient) List(ctx context.Context) (*api.ListLocationsOutputBody, error) {
	resp, err := l.raw.ListLocationsWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, responseError(resp.StatusCode(), resp.HTTPResponse.Header, resp.Body, resp.ApplicationproblemJSONDefault)
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("list locations: HTTP 200 lacks a JSON response")
	}
	return resp.JSON200, nil
}
