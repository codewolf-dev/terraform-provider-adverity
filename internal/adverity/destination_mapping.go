package adverity

import (
	"net/url"
	"strconv"
)

type DestinationMappingConfig struct {
	DatastreamId *int64       `json:"datastream,omitempty"`
	Enabled      *bool        `json:"enabled,omitempty"`
	TableName    *string      `json:"table_name,omitempty"`
	Parameters   *[]Parameter `json:"-"`
}

func (c *DestinationMappingConfig) MarshalJSON() ([]byte, error) {
	return FlattenedMarshal(c, c.Parameters)
}

type DestinationMappingResponse struct {
	ID            int64  `json:"id"`
	DestinationID int64  `json:"target"`
	DatastreamID  int64  `json:"datastream"`
	Enabled       bool   `json:"enabled"`
	TableName     string `json:"table_name"`
}

func (c *Client) CreateDestinationMapping(destinationTypeId, destinationId int, req *DestinationMappingConfig) (*DestinationMappingResponse, error) {
	r, _ := url.JoinPath("target-types", strconv.Itoa(destinationTypeId), "targets", strconv.Itoa(destinationId), "mappings", "/")
	p, _ := url.Parse(r)

	resp, err := Create[DestinationMappingConfig, DestinationMappingResponse](c, p, req, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) ReadDestinationMapping(destinationTypeId, destinationId, destinationMappingId int) (*DestinationMappingResponse, error) {
	r, _ := url.JoinPath("target-types", strconv.Itoa(destinationTypeId), "targets", strconv.Itoa(destinationId), "mappings", strconv.Itoa(destinationMappingId), "/")
	p, _ := url.Parse(r)

	resp, err := Read[DestinationMappingResponse](c, p, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) UpdateDestinationMapping(destinationTypeId, destinationId, destinationMappingId int, req *DestinationMappingConfig) (*DestinationMappingResponse, error) {
	r, _ := url.JoinPath("target-types", strconv.Itoa(destinationTypeId), "targets", strconv.Itoa(destinationId), "mappings", strconv.Itoa(destinationMappingId), "/")
	p, _ := url.Parse(r)

	resp, err := Update[DestinationMappingConfig, DestinationMappingResponse](c, p, req, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) DeleteDestinationMapping(destinationTypeId, destinationId, destinationMappingId int) (*DestinationMappingResponse, error) {
	r, _ := url.JoinPath("target-types", strconv.Itoa(destinationTypeId), "targets", strconv.Itoa(destinationId), "mappings", strconv.Itoa(destinationMappingId), "/")
	p, _ := url.Parse(r)

	resp, err := Delete[DestinationMappingResponse](c, p, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
