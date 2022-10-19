package bitbucket

import (
	"context"
	"fmt"
)

type ProjectService interface {
	GetProject(context.Context, string) (*Project, error)
	CreateProject(context.Context, *CreateProjectRequest) (*Project, error)
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

func (ps *projectService) GetProject(ctx context.Context, key string) (*Project, error) {
	req, err := ps.client.newRequest("GET", fmt.Sprintf("projects/%s", key), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request for getting projects: %w", err)
	}

	p := Project{}
	err = ps.client.do(ctx, req, &p)
	if err != nil {
		return nil, fmt.Errorf("error fetching projects: %w", err)
	}
	return &p, nil
}

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Key         string `json:"key"`
	Description string `json:"description,omitempty"`
	Public      bool   `json:"public,omitempty"`
}

func (ps *projectService) CreateProject(ctx context.Context, createReq *CreateProjectRequest) (*Project, error) {
	req, err := ps.client.newRequest("POST", "projects", createReq)
	if err != nil {
		return nil, fmt.Errorf("error creating request for creating project: %w", err)
	}

	p := Project{}
	err = ps.client.do(ctx, req, &p)
	if err != nil {
		return nil, fmt.Errorf("error creating project: %w", err)
	}

	return &p, nil
}
