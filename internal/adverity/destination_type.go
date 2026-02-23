// Copyright codewolf.dev 2025, 0
// SPDX-License-Identifier: MPL-2.0

package adverity

import (
	"net/url"
)

type DestinationType struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	URL          string `json:"url"`
	Destinations string `json:"targets"`
}

type destinationTypeQueryResponse struct {
	Count    int64             `json:"count"`
	Next     string            `json:"next"`
	Previous string            `json:"previous"`
	Results  []DestinationType `json:"results"`
}

func (c *Client) QueryDestinationTypes(searchTerm string) ([]DestinationType, error) {
	r, _ := url.JoinPath("target-types", "/")
	p, _ := url.Parse(r)

	q := &url.Values{}
	q.Add("search", searchTerm)

	resp, err := Read[destinationTypeQueryResponse](c, p, q)
	if err != nil {
		return nil, err
	}

	return resp.Results, nil
}
