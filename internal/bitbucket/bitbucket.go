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
		return nil, errors.New("error fetching project with key " + key)
	}
	if r.StatusCode == 404 {
		return nil, &Error{
			msg:  fmt.Sprintf("Project with key %s does not exist", key),
			Code: ErrNotFound,
		}
	}
	p := &Project{}

	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(p)
	if err != nil {
		return nil, errors.New("error decoding project from api response")
	}
	fmt.Println("GOT PROJECT BACK, TYPE: " + string(p.Type))
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

// ErrorCode defines the code of an error.
type ErrorCode string

const (
	ErrPermission ErrorCode = "permission" // Permission denied.
	ErrNotFound   ErrorCode = "not_found"  // Resource not found.
)

// Error represents common errors originating from the Client.
type Error struct {
	// msg contains the human readable string
	msg string

	// Code specifies the error code. i.e; NotFound, Permission, etc...
	Code ErrorCode
}

// Error returns the string representation of the error.
func (e *Error) Error() string { return e.msg }
