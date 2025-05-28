// Copyright (c) codewolf.dev
// SPDX-License-Identifier: MPL-2.0

package adverity

import (
	"net/url"
)

type ConnectionType struct {
	ID           int64    `json:"id"`
	Name         string   `json:"name"`
	Slug         string   `json:"slug"`
	URL          string   `json:"url"`
	Categories   []string `json:"categories"`
	Keywords     []string `json:"keywords"`
	IsDeprecated bool     `json:"is_deprecated"`
	LogoURL      string   `json:"logo_url"`
	CreateURL    string   `json:"create_url"`
	Connections  string   `json:"connections"`
}

type connectionTypeQueryResponse struct {
	Count    int64            `json:"count"`
	Next     string           `json:"next"`
	Previous string           `json:"previous"`
	Results  []ConnectionType `json:"results"`
}

func (c *Client) QueryConnectionTypes(searchTerm string) ([]ConnectionType, error) {
	r, _ := url.JoinPath("connection-types", "/")
	p, _ := url.Parse(r)

	q := &url.Values{}
	q.Add("search", searchTerm)

	resp, err := Read[connectionTypeQueryResponse](c, p, q)
	if err != nil {
		return nil, err
	}

	return resp.Results, nil
}
