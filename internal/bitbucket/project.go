package bitbucket

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ProjectService interface {
	GetProject(string) (*Project, error)
}

type projectService struct {
	client *Client
}

func (ps *projectService) GetProject(key string) (*Project, error) {
	r, err := ps.client.client.Get(fmt.Sprintf("%s/projects/%s", ps.client.baseURL, key))
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
	return p, nil
}
