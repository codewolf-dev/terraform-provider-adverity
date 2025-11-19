// Copyright (c) codewolf.dev
// SPDX-License-Identifier: MPL-2.0

package adverity

import (
	"net/url"
)

type AuthorizationType struct {
	ID             int64    `json:"id"`
	Name           string   `json:"name"`
	Slug           string   `json:"slug"`
	URL            string   `json:"url"`
	Categories     []string `json:"categories"`
	Keywords       []string `json:"keywords"`
	IsDeprecated   bool     `json:"is_deprecated"`
	LogoURL        string   `json:"logo_url"`
	CreateURL      string   `json:"create_url"`
	Authorizations string   `json:"connections"`
}

type authorizationTypeQueryResponse struct {
	Count    int64               `json:"count"`
	Next     string              `json:"next"`
	Previous string              `json:"previous"`
	Results  []AuthorizationType `json:"results"`
}

func (c *Client) QueryAuthorizationTypes(searchTerm string) ([]AuthorizationType, error) {
	r, _ := url.JoinPath("connection-types", "/")
	p, _ := url.Parse(r)

	q := &url.Values{}
	q.Add("search", searchTerm)

	resp, err := Read[authorizationTypeQueryResponse](c, p, q)
	if err != nil {
		return nil, err
	}

	return resp.Results, nil
}
