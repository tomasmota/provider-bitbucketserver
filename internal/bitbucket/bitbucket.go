package bitbucket

import (
	"encoding/json"
	"errors"
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

var (
	ErrPermission = errors.New("permission") // Permission denied.
	ErrNotFound   = errors.New("not_found")  // Resource not found.
)

func NewClient(baseURL string, base64creds string) (*Client, error) {
	fmt.Printf("creating bitbucket client, endpoint: %s\n", baseURL)

	fmt.Printf("\n\n\n\n creds: %s \n\n\n\n", base64creds)
	c := &Client{
		baseURL: fmt.Sprintf("%s%s", baseURL, apiPath),
		client:  NewBasicAuthHttpClient(base64creds),
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

type Project struct {
	Name        string `json:"name"`
	Key         string `json:"key"`
	ID          int    `json:"id"`
	Description string `json:"description"`
	Scope       string `json:"scope,omitempty"`
	Type        string `json:"type"`
	Public      bool   `json:"public"`
}

func (c *Client) GetProject(key string) (*Project, error) {
	r, err := c.client.Get(fmt.Sprintf("%s/projects/%s", c.baseURL, key))
	if err != nil {
		return nil, fmt.Errorf("error fetching project with key %s: %w", key, err)
	}
	switch r.StatusCode {
	case 404:
		return nil, ErrNotFound
	case 401:
		return nil, ErrPermission
	}
	p := &Project{}

	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(p)
	if err != nil {
		return nil, errors.New("error decoding project from api response")
	}
	fmt.Println("GOT PROJECT BACK, TYPE: " + string(p.Type)) // TODO: remove after testing
	return p, nil
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
