package bitbucket

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
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
		client: &http.Client{
			Timeout: time.Second * 10,
			Transport: &oauth2.Transport{
				Source: oauth2.StaticTokenSource(
					&oauth2.Token{
						AccessToken: strings.TrimSpace(token),
					},
				),
			},
		},
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
