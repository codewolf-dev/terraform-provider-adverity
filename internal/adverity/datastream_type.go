// Copyright (c) codewolf.dev
// SPDX-License-Identifier: MPL-2.0

package adverity

import (
	"net/url"
)

type DatastreamType struct {
	ID              int64    `json:"id"`
	Name            string   `json:"name"`
	Slug            string   `json:"slug"`
	URL             string   `json:"url"`
	Categories      []string `json:"categories"`
	Keywords        []string `json:"keywords"`
	IsDeprecated    bool     `json:"is_deprecated"`
	LogoURL         string   `json:"logo_url"`
	CreateURL       string   `json:"create_url"`
	Datastream      string   `json:"datastreams"`
	ConnectionTypes []string `json:"connection_types"`
}

type datastreamTypeQueryResponse struct {
	Count    int64            `json:"count"`
	Next     string           `json:"next"`
	Previous string           `json:"previous"`
	Results  []DatastreamType `json:"results"`
}

func (c *Client) QueryDatastreamTypes(searchTerm string) ([]DatastreamType, error) {
	r, _ := url.JoinPath("datastream-types", "/")
	p, _ := url.Parse(r)

	q := &url.Values{}
	q.Add("search", searchTerm)

	resp, err := Read[datastreamTypeQueryResponse](c, p, q)
	if err != nil {
		return nil, err
	}

	return resp.Results, nil
}
