package bitbucket

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	apiPath = "/rest/api/1.0"
)

type Client struct {
	// client represents the HTTP client used for making HTTP requests.
	client *http.Client

	// base URL for the bitbucket server
	baseURL string

	Projects ProjectService
}

var (
	ErrPermission = errors.New("permission") // Permission denied.
	ErrNotFound   = errors.New("not_found")  // Resource not found.
)

func NewClient(baseURL string, base64creds string) (*Client, error) {
	fmt.Printf("creating bitbucket client, endpoint: %s\n", baseURL)

	c := &Client{
		baseURL: fmt.Sprintf("%s%s", baseURL, apiPath),
		client:  NewBasicAuthHttpClient(base64creds),
	}
	err := c.ping()
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

// func (c *Client) newRequest(method string, path string, body interface{}) (*http.Request, error) {
// 	url := fmt.Sprintf("%s/projects/%s", ps.client.baseURL, key)
// 	if err != nil {
// 		return nil, err
// 	}
// }

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
