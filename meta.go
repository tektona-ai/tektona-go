package tektona

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tektona-ai/tektona-go/api"
)

// Meta reads public API metadata.
type Meta interface {
	Get(ctx context.Context) (*api.MetaResponse, error)
}

type metaClient struct{ raw *api.ClientWithResponses }

// Meta returns the metadata methods.
func (c *Client) Meta() Meta { return c.meta }

func (m *metaClient) Get(ctx context.Context) (*api.MetaResponse, error) {
	resp, err := m.raw.GetMetaWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, responseError(resp.StatusCode(), resp.HTTPResponse.Header, resp.Body, resp.ApplicationproblemJSONDefault)
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("get metadata: HTTP 200 lacks a JSON response")
	}
	return resp.JSON200, nil
}
