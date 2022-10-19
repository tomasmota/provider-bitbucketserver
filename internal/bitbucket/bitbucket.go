package bitbucket

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	apiPath = "/rest/api/1.0"
)

type Client struct {
	// client represents the HTTP client used for making HTTP requests.
	client *http.Client

	// base URL for the bitbucket server
	baseURL *url.URL

	Projects ProjectService
}

var (
	ErrPermission = errors.New("permission") // Permission denied.
	ErrNotFound   = errors.New("not_found")  // Resource not found.
)

func NewClient(baseURL string, base64creds string) (*Client, error) {
	fmt.Printf("creating bitbucket client, endpoint: %s\n", baseURL)
	pBaseURL, err := url.Parse(fmt.Sprintf("%s%s", baseURL, apiPath))
	if err != nil {
		return nil, err
	}

	c := &Client{
		baseURL: pBaseURL,
		client:  NewBasicAuthHttpClient(base64creds),
	}

	err = c.ping()
	if err != nil {
		return nil, fmt.Errorf("error creating bitbucket client: %w", err)
	}

	c.Projects = &projectService{client: c}

	return c, nil
}

func (c *Client) ping() error {
	_, err := c.client.Get(fmt.Sprintf("%s/projects", c.baseURL))
	if err != nil {
		return fmt.Errorf("error fetching projects: %w", err)
	}
	return nil
}

type Project struct {
	Name        string `json:"name"`
	Key         string `json:"key"`
	ID          int    `json:"id"`
	Description string `json:"description"`
	Scope       string `json:"scope,omitempty"`
	Type        string `json:"type"`
	Public      bool   `json:"public"`
}

type bearerRoundTripper struct {
	base64creds string // Basic auth creds
}

func (t *bearerRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Add("Authorization", "Basic "+t.base64creds)
	return http.DefaultTransport.RoundTrip(r)
}

func NewBasicAuthHttpClient(creds string) *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
		Transport: &bearerRoundTripper{
			base64creds: creds,
		},
	}
}
