package bitbucket

import (
	"fmt"
	"net/http"
	"time"
)

const (
	apiPath = "/rest/api/1.0"
)

type Client struct {
	baseURL string
	client  *http.Client
}

func NewClient(baseURL string, token string) (*Client, error) {
	fmt.Printf("creating bitbucket client, endpoint: %s\n", baseURL)

	c := &Client{
		baseURL: fmt.Sprintf("%s%s", baseURL, apiPath),
		client:  NewBearerHttpClient(token),
	}
	err := c.ping()
	if err != nil {
		return nil, fmt.Errorf("error creating bitbucket client: %w", err)
	}

	return c, nil

}

func (c *Client) ping() error {
	_, err := c.client.Get(fmt.Sprintf("%s/projects", c.baseURL))
	if err != nil {
		return fmt.Errorf("error fetching projects: %w", err)
	}
	return nil
}

// TODO: replace this with oauth package roundtripper
type bearerRoundTripper struct {
	token string // Bearer token
}

func (t *bearerRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Add("Authorization", "Bearer "+t.token)
	return http.DefaultTransport.RoundTrip(r)
}

func NewBearerHttpClient(token string) *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
		Transport: &bearerRoundTripper{
			token: token,
		},
	}
}
