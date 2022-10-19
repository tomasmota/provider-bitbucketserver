package bitbucket

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	apiPath       = "/rest/api/1.0/"
	jsonMediaType = "application/json"
)

type Client struct {
	// client represents the HTTP client used for making HTTP requests.
	client *http.Client

	// headers are used to override request headers for every single HTTP request
	headers map[string]string

	// base URL for the bitbucket server + apiPath
	baseURL *url.URL

	Projects ProjectService
}

var (
	ErrPermission        = errors.New("permission")         // Permission denied.
	ErrNotFound          = errors.New("not_found")          // Resource not found.
	ErrResponseMalformed = errors.New("response_malformed") // Resource not found.
)

func NewClient(baseURL string, base64creds string) (*Client, error) {
	fmt.Printf("creating bitbucket client, endpoint: %s\n", baseURL)
	pBaseURL, err := url.Parse(fmt.Sprintf("%s%s", baseURL, apiPath))
	if err != nil {
		return nil, err
	}

	c := &Client{
		baseURL: pBaseURL,
		client:  &http.Client{Timeout: time.Second * 10},
		headers: map[string]string{"Authorization": "Basic " + base64creds},
	}

	err = c.ping()
	if err != nil {
		return nil, fmt.Errorf("error creating bitbucket client: %w", err)
	}

	c.Projects = &projectService{client: c}

	return c, nil
}

func (c *Client) ping() error {
	req, err := c.newRequest("GET", "projects", nil)
	if err != nil {
		return fmt.Errorf("error creating request for getting projects: %w", err)
	}

	err = c.do(context.Background(), req, nil)
	if err != nil {
		return fmt.Errorf("error fetching projects: %w", err)
	}
	return nil
}

func (c *Client) newRequest(method string, path string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	switch method {
	case http.MethodGet:
		req, err = http.NewRequest(method, u.String(), nil)
		if err != nil {
			return nil, err
		}
	default:
		buf := new(bytes.Buffer)
		if body != nil {
			err = json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}
		}

		req, err = http.NewRequest(method, u.String(), buf)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", jsonMediaType)
	}

	req.Header.Set("Accept", jsonMediaType)

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

// do makes an HTTP request and populates the given struct v from the response.
func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) error {
	req = req.WithContext(ctx)
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return c.handleResponse(ctx, res, v)
}

// handleResponse makes an HTTP request and populates the given struct v from
// the response.  This is meant for internal testing and shouldn't be used
// directly. Instead please use `Client.do`.
func (c *Client) handleResponse(ctx context.Context, res *http.Response, v interface{}) error {
	out, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	switch res.StatusCode {
	case 404:
		return ErrNotFound
	case 401:
		return ErrPermission
	}

	// this means we don't care about unmarshaling the response body into v
	if v == nil || res.StatusCode == http.StatusNoContent {
		return nil
	}

	err = json.Unmarshal(out, &v)
	if err != nil {
		var jsonErr *json.SyntaxError
		if errors.As(err, &jsonErr) {
			return ErrResponseMalformed
		}
		return err
	}

	fmt.Println(v)
	return nil
}

// type bearerRoundTripper struct {
// 	base64creds string // Basic auth creds
// }

// func (t *bearerRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
// 	r.Header.Add("Authorization", "Basic "+t.base64creds)
// 	return http.DefaultTransport.RoundTrip(r)
// }

// func NewBasicAuthHttpClient(creds string) *http.Client {
// 	return &http.Client{
// 		Timeout: time.Second * 10,
// 		Transport: &bearerRoundTripper{
// 			base64creds: creds,
// 		},
// 	}
// }
