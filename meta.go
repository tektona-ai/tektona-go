package tektona

import (
	"context"
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
	return resourceJSON[api.MetaResponse](func() (*http.Response, error) {
		return m.raw.GetMeta(ctx)
	}, http.StatusOK, "get metadata")
}
