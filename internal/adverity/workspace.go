package adverity

import (
	"net/url"
)

type CreateWorkspaceRequest struct {
	DatalakeID int    `json:"datalake_id,omitempty" url:"datalake_id,omitempty"`
	ParentID   int    `json:"parent_id,omitempty"   url:"parent_id,omitempty"`
	Name       string `json:"name,omitempty"        url:"name,omitempty"`
}

type ReadWorkspaceRequest struct {
	StackSlug string `json:"stack_slug,omitempty" url:"stack_slug,omitempty"`
}

type UpdateWorkspaceRequest struct {
	StackSlug  string `json:"stack_slug,omitempty"  url:"stack_slug,omitempty"`
	DatalakeID int    `json:"datalake_id,omitempty" url:"datalake_id,omitempty"`
	ParentID   int    `json:"parent_id,omitempty"   url:"parent_id,omitempty"`
	Name       string `json:"name,omitempty"        url:"name,omitempty"`
}

type DeleteWorkspaceRequest struct {
	StackSlug string `json:"stack_slug,omitempty" url:"stack_slug,omitempty"`
}

type WorkspaceResponse struct {
	AddConnectionURL string      `json:"add_connection_url"`
	AddDatastreamURL string      `json:"add_datastream_url"`
	ChangeURL        string      `json:"change_url"`
	Datalake         string      `json:"datalake"`
	Destination      interface{} `json:"destination"`
	ExtractsURL      string      `json:"extracts_url"`
	IssuesURL        string      `json:"issues_url"`
	ID               int         `json:"id"`
	Name             string      `json:"name"`
	OverviewURL      string      `json:"overview_url"`
	Parent           string      `json:"parent"`
	ParentID         int         `json:"parent_id"`
	Slug             string      `json:"slug"`
	URL              string      `json:"url"`
	Counts           struct {
		Connections int `json:"connections"`
		Datastreams int `json:"datastreams"`
	} `json:"counts"`
	Permissions struct {
		IsCreator           bool `json:"isCreator"`
		IsDatastreamManager bool `json:"isDatastreamManager"`
		IsViewer            bool `json:"isViewer"`
	} `json:"permissions"`
	DefaultManageExtractNames bool   `json:"default_manage_extract_names"`
	Updated                   string `json:"updated"`
	Created                   string `json:"created"`
}

func (c *Client) CreateWorkspace(req *CreateWorkspaceRequest) (*WorkspaceResponse, error) {
	path, _ := url.JoinPath("stacks", "/")

	resp, err := Create[CreateWorkspaceRequest, WorkspaceResponse](c, path, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) ReadWorkspace(req *ReadWorkspaceRequest) (*WorkspaceResponse, error) {
	path, _ := url.JoinPath("stacks", req.StackSlug, "/")

	resp, err := Read[WorkspaceResponse](c, path)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) UpdateWorkspace(req *UpdateWorkspaceRequest) (*WorkspaceResponse, error) {
	path, _ := url.JoinPath("stacks", req.StackSlug, "/")

	resp, err := Update[UpdateWorkspaceRequest, WorkspaceResponse](c, path, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) DeleteWorkspace(req *DeleteWorkspaceRequest) (*WorkspaceResponse, error) {
	path, _ := url.JoinPath("stacks", req.StackSlug, "/")

	resp, err := Delete[WorkspaceResponse](c, path)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
