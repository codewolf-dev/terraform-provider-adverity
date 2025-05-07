package adverity

import (
	"net/url"
)

type WorkspaceConfig struct {
	DatalakeID         int    `json:"datalake_id,omitempty"`
	Name               string `json:"name,omitempty"`
	ParentID           int    `json:"parent_id,omitempty"`
	Destination        int    `json:"destination,omitempty"`
	ManageExtractNames bool   `json:"default_manage_extract_names,omitempty"`
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
	ManageExtractNames bool   `json:"default_manage_extract_names"`
	Updated            string `json:"updated"`
	Created            string `json:"created"`
}

func (c *Client) CreateWorkspace(req *WorkspaceConfig) (*WorkspaceResponse, error) {
	r, _ := url.JoinPath("stacks", "/")
	p, _ := url.Parse(r)

	return Create[WorkspaceConfig, WorkspaceResponse](c, p, req, nil)
}

func (c *Client) ReadWorkspace(stackSlug string) (*WorkspaceResponse, error) {
	r, _ := url.JoinPath("stacks", stackSlug, "/")
	p, _ := url.Parse(r)

	return Read[WorkspaceResponse](c, p, nil)
}

func (c *Client) UpdateWorkspace(stackSlug string, req *WorkspaceConfig) (*WorkspaceResponse, error) {
	r, _ := url.JoinPath("stacks", stackSlug, "/")
	p, _ := url.Parse(r)

	return Update[WorkspaceConfig, WorkspaceResponse](c, p, req, nil)
}

func (c *Client) DeleteWorkspace(stackSlug string) (*WorkspaceResponse, error) {
	r, _ := url.JoinPath("stacks", stackSlug, "/")
	p, _ := url.Parse(r)

	return Delete[WorkspaceResponse](c, p, nil)
}
