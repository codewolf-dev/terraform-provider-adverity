package adverity

import (
	"net/url"
	"strconv"
)

type ConnectionConfig struct {
	Name       *string      `json:"name,omitempty"`
	StackID    *int64       `json:"stack,omitempty"`
	Parameters *[]Parameter `json:"-"`
}

func (c *ConnectionConfig) MarshalJSON() ([]byte, error) {
	return FlattenedMarshal(c, c.Parameters)
}

type ConnectionResponse struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	MetadataSlack int64  `json:"metadata_slack"`
	StackID       int64  `json:"stack"`
	App           int64  `json:"app"`
	User          int64  `json:"user"`
	IsAuthorized  bool   `json:"is_authorized"`
}

func (c *Client) CreateConnection(connectionTypeId int, req *ConnectionConfig) (*ConnectionResponse, error) {
	r, _ := url.JoinPath("connection-types", strconv.Itoa(connectionTypeId), "connections", "/")
	p, _ := url.Parse(r)

	return Create[ConnectionConfig, ConnectionResponse](c, p, req, nil)
}

func (c *Client) ReadConnection(connectionTypeId, connectionId int) (*ConnectionResponse, error) {
	r, _ := url.JoinPath("connection-types", strconv.Itoa(connectionTypeId), "connections", strconv.Itoa(connectionId), "/")
	p, _ := url.Parse(r)

	return Read[ConnectionResponse](c, p, nil)
}

func (c *Client) UpdateConnection(connectionTypeId, connectionId int, req *ConnectionConfig) (*ConnectionResponse, error) {
	r, _ := url.JoinPath("connection-types", strconv.Itoa(connectionTypeId), "connections", strconv.Itoa(connectionId), "/")
	p, _ := url.Parse(r)

	return Update[ConnectionConfig, ConnectionResponse](c, p, req, nil)
}

func (c *Client) DeleteConnection(connectionTypeId, connectionId int) (*ConnectionResponse, error) {
	r, _ := url.JoinPath("connection-types", strconv.Itoa(connectionTypeId), "connections", strconv.Itoa(connectionId), "/")
	p, _ := url.Parse(r)

	return Delete[ConnectionResponse](c, p, nil)
}
