package bitbucket

import (
	"context"
	"fmt"
)

type ProjectService interface {
	GetProject(string) (*Project, error)
}

type projectService struct {
	client *Client
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

func (ps *projectService) GetProject(key string) (*Project, error) {
	req, err := ps.client.newRequest("GET", fmt.Sprintf("projects/%s", key), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request for getting projects: %w", err)
	}

	p := Project{}
	err = ps.client.do(context.Background(), req, &p)
	if err != nil {
		return nil, fmt.Errorf("error fetching projects: %w", err)
	}
	return &p, nil
}
