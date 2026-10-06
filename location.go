package tektona

import (
	"context"
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
	return resourceJSON[api.ListLocationsOutputBody](func() (*http.Response, error) {
		return l.raw.ListLocations(ctx)
	}, http.StatusOK, "list locations")
}
