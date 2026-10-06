// Package tektona calls the Tektona public HTTP API. The api package holds generated types and raw methods.
package tektona

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/tektona-ai/tektona-go/api"
)

const defaultBaseURL = "https://api.tektona.ai"

type config struct {
	apiKey, baseURL, org, project                 string
	httpClient                                    *http.Client
	keySet, urlSet, orgSet, projectSet, clientSet bool
}

// Option configures a Client.
type Option func(*config)

// WithAPIKey selects an API key instead of TEKTONA_API_KEY.
func WithAPIKey(key string) Option { return func(c *config) { c.apiKey, c.keySet = key, true } }

// WithBaseURL selects the API origin or an origin with a path prefix.
func WithBaseURL(baseURL string) Option {
	return func(c *config) { c.baseURL, c.urlSet = baseURL, true }
}

// WithHTTPClient selects the HTTP client used for requests.
func WithHTTPClient(client *http.Client) Option {
	return func(c *config) { c.httpClient, c.clientSet = client, true }
}

// WithOrg sets the default organization for scoped calls.
func WithOrg(org string) Option { return func(c *config) { c.org, c.orgSet = org, true } }

// WithProject sets the default project for scoped calls.
func WithProject(project string) Option {
	return func(c *config) { c.project, c.projectSet = project, true }
}

// Client shares authentication and scope defaults across resources.
type Client struct {
	raw          *api.ClientWithResponses
	org, project string
	sandbox      Sandbox
	location     Location
	meta         Meta
	organization Organization
	projectAPI   Project
}

// NewClient creates a client. It reads TEKTONA_API_KEY, TEKTONA_API_URL, TEKTONA_ORG, and TEKTONA_PROJECT.
func NewClient(opts ...Option) (*Client, error) {
	cfg := config{}
	for _, opt := range opts {
		opt(&cfg)
	}
	if !cfg.keySet {
		cfg.apiKey = os.Getenv("TEKTONA_API_KEY")
	}
	if cfg.apiKey == "" {
		return nil, fmt.Errorf("API key is required: set TEKTONA_API_KEY or use WithAPIKey")
	}
	if !cfg.urlSet {
		cfg.baseURL = os.Getenv("TEKTONA_API_URL")
		if cfg.baseURL == "" {
			cfg.baseURL = defaultBaseURL
		}
	}
	base, err := url.Parse(cfg.baseURL)
	if err != nil || (base.Scheme != "http" && base.Scheme != "https") || base.Host == "" || base.User != nil || base.RawQuery != "" || base.Fragment != "" {
		return nil, fmt.Errorf("invalid API URL %q: expected an HTTP origin or path prefix", cfg.baseURL)
	}
	if !cfg.orgSet {
		cfg.org = os.Getenv("TEKTONA_ORG")
	}
	if !cfg.projectSet {
		cfg.project = os.Getenv("TEKTONA_PROJECT")
	}
	if cfg.clientSet && cfg.httpClient == nil {
		return nil, fmt.Errorf("HTTP client cannot be nil")
	}
	if cfg.httpClient == nil {
		cfg.httpClient = http.DefaultClient
	}
	client := *cfg.httpClient
	transport := client.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	client.Transport = bearerTransport{next: transport, key: cfg.apiKey, host: base.Host, scheme: base.Scheme}
	base.Path = strings.TrimSuffix(base.Path, "/") + "/"
	raw, err := api.NewClientWithResponses(base.String(), api.WithHTTPClient(&client))
	if err != nil {
		return nil, err
	}
	c := &Client{raw: raw, org: cfg.org, project: cfg.project}
	c.sandbox = &sandboxClient{raw: raw, org: cfg.org, project: cfg.project}
	c.location = &locationClient{raw: raw}
	c.meta = &metaClient{raw: raw}
	c.organization = &organizationClient{raw: raw, org: c.org}
	c.projectAPI = &projectClient{raw: raw, org: c.org}
	return c, nil
}

type bearerTransport struct {
	next   http.RoundTripper
	key    string
	host   string
	scheme string
}

func (t bearerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.Host != t.host || req.URL.Scheme != t.scheme {
		return t.next.RoundTrip(req)
	}
	copy := req.Clone(req.Context())
	copy.Header.Set("Authorization", "Bearer "+t.key)
	return t.next.RoundTrip(copy)
}

// Sandbox returns the sandbox methods. The same interface is returned on each call.
func (c *Client) Sandbox() Sandbox { return c.sandbox }

// Raw returns generated methods without scope defaults or resource error conversion.
func (c *Client) Raw() *api.ClientWithResponses { return c.raw }
