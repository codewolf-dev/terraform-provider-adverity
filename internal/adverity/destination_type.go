package adverity

import (
	"net/url"
)

type DestinationType struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	URL     string `json:"url"`
	Targets string `json:"targets"`
}

type DestinationTypeQueryResponse struct {
	Count    int               `json:"count"`
	Next     string            `json:"next"`
	Previous string            `json:"previous"`
	Results  []DestinationType `json:"results"`
}

func (c *Client) QueryDestinationTypes(searchTerm string) (*DestinationTypeQueryResponse, error) {
	r, _ := url.JoinPath("target-types", "/")
	p, _ := url.Parse(r)

	q := &url.Values{}
	q.Add("search", searchTerm)

	resp, err := Read[DestinationTypeQueryResponse](c, p, q)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
