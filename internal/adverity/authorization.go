// Copyright codewolf.dev 2025, 0
// SPDX-License-Identifier: MPL-2.0

package adverity

import (
	"net/url"
	"strconv"
)

type AuthorizationConfig struct {
	Name       *string      `json:"name,omitempty"`
	StackID    *int64       `json:"stack,omitempty"`
	Parameters *[]Parameter `json:"-"`
}

func (c *AuthorizationConfig) MarshalJSON() ([]byte, error) {
	return FlattenedMarshal(c, c.Parameters)
}

type AuthorizationResponse struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	MetadataSlack int64  `json:"metadata_slack"`
	StackID       int64  `json:"stack"`
	App           int64  `json:"app"`
	User          int64  `json:"user"`
	IsAuthorized  bool   `json:"is_authorized"`
}

func (c *Client) CreateAuthorization(connectionTypeId int, req *AuthorizationConfig) (*AuthorizationResponse, error) {
	r, _ := url.JoinPath("connection-types", strconv.Itoa(connectionTypeId), "connections", "/")
	p, _ := url.Parse(r)

	return Create[AuthorizationConfig, AuthorizationResponse](c, p, req, nil)
}

func (c *Client) ReadAuthorization(connectionTypeId, connectionId int) (*AuthorizationResponse, error) {
	r, _ := url.JoinPath("connection-types", strconv.Itoa(connectionTypeId), "connections", strconv.Itoa(connectionId), "/")
	p, _ := url.Parse(r)

	return Read[AuthorizationResponse](c, p, nil)
}

func (c *Client) UpdateAuthorization(connectionTypeId, connectionId int, req *AuthorizationConfig) (*AuthorizationResponse, error) {
	r, _ := url.JoinPath("connection-types", strconv.Itoa(connectionTypeId), "connections", strconv.Itoa(connectionId), "/")
	p, _ := url.Parse(r)

	return Update[AuthorizationConfig, AuthorizationResponse](c, p, req, nil)
}

func (c *Client) DeleteAuthorization(connectionTypeId, connectionId int) (*AuthorizationResponse, error) {
	r, _ := url.JoinPath("connection-types", strconv.Itoa(connectionTypeId), "connections", strconv.Itoa(connectionId), "/")
	p, _ := url.Parse(r)

	return Delete[AuthorizationResponse](c, p, nil)
}
