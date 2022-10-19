package bitbucket

import (
	"context"
	"fmt"
)

type ProjectService interface {
	GetProject(context.Context, *GetProjectRequest) (*Project, error)
	CreateProject(context.Context, *CreateProjectRequest) (*Project, error)
	DeleteProject(context.Context, *DeleteProjectRequest) error
	UpdateProject(context.Context, *UpdateProjectRequest) (*Project, error)
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

type GetProjectRequest struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

func (ps *projectService) GetProject(ctx context.Context, getReq *GetProjectRequest) (*Project, error) {
	req, err := ps.client.newRequest("GET", fmt.Sprintf("projects/%s", getReq.Key), nil)
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

type DeleteProjectRequest struct {
	Key string `json:"key"`
}

func (ps *projectService) DeleteProject(ctx context.Context, deleteReq *DeleteProjectRequest) error {
	req, err := ps.client.newRequest("DELETE", fmt.Sprintf("projects/%s", deleteReq.Key), nil)
	if err != nil {
		return fmt.Errorf("error creating request for deleting project: %w", err)
	}

	err = ps.client.do(ctx, req, nil)
	if err != nil {
		return fmt.Errorf("error deleting project: %w", err)
	}

	return nil
}

type UpdateProjectRequest struct {
	Key         string `json:"key"`
	Description string `json:"description,omitempty"`
	Public      bool   `json:"public,omitempty"`
}

func (ps *projectService) UpdateProject(ctx context.Context, updateReq *UpdateProjectRequest) (*Project, error) {
	req, err := ps.client.newRequest("PUT", fmt.Sprintf("projects/%s", updateReq.Key), updateReq)
	if err != nil {
		return nil, fmt.Errorf("error creating request for updating project: %w", err)
	}

	p := Project{}
	err = ps.client.do(ctx, req, &p)
	if err != nil {
		return nil, fmt.Errorf("error updating project: %w", err)
	}

	return &p, nil
}
