package adverity

import (
	"net/url"
)

type ConnectionType struct {
	ID           int      `json:"id"`
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

type ConnectionTypeQueryResponse struct {
	Count    int              `json:"count"`
	Next     string           `json:"next"`
	Previous string           `json:"previous"`
	Results  []ConnectionType `json:"results"`
}

func (c *Client) QueryConnectionTypes(searchTerm string) (*ConnectionTypeQueryResponse, error) {
	r, _ := url.JoinPath("connection-types", "/")
	p, _ := url.Parse(r)

	q := &url.Values{}
	q.Add("search", searchTerm)

	resp, err := Read[ConnectionTypeQueryResponse](c, p, q)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
