package adverity

import (
	"net/url"
)

type DatastreamType struct {
	ID              int      `json:"id"`
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

type DatastreamTypeQueryResponse struct {
	Count    int              `json:"count"`
	Next     string           `json:"next"`
	Previous string           `json:"previous"`
	Results  []DatastreamType `json:"results"`
}

func (c *Client) QueryDatastreamTypes(searchTerm string) (*DatastreamTypeQueryResponse, error) {
	r, _ := url.JoinPath("datastream-types", "/")
	p, _ := url.Parse(r)

	q := &url.Values{}
	q.Add("search", searchTerm)

	resp, err := Read[DatastreamTypeQueryResponse](c, p, q)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
