package adverity

import (
	"net/url"
	"strconv"
)

type DestinationConfig struct {
	Name                   *string      `json:"name,omitempty"`
	StackID                *int64       `json:"stack,omitempty"`
	AuthID                 *int64       `json:"auth,omitempty"`
	SchemaMapping          *bool        `json:"schema_mapping,omitempty"`
	ForceString            *bool        `json:"force_string,omitempty"`
	FormatHeaders          *bool        `json:"format_headers,omitempty"`
	ColumnNamesToLowerCase *bool        `json:"column_names_to_lowercase,omitempty"`
	HeadersFormatting      *int64       `json:"headers_formatting,omitempty"`
	Parameters             *[]Parameter `json:"-"`
}

func (c *DestinationConfig) MarshalJSON() ([]byte, error) {
	return FlattenedMarshal(c, c.Parameters)
}

type DestinationResponse struct {
	ID                      int64  `json:"id"`
	LogoURL                 string `json:"logo_url"`
	IsSchemaMappingRequired bool   `json:"is_schema_mapping_required"`
	Name                    string `json:"name"`
	SchemaMapping           bool   `json:"schema_mapping"`
	ForceString             bool   `json:"force_string"`
	FormatHeaders           bool   `json:"format_headers"`
	ColumnNamesToLowerCase  bool   `json:"column_names_to_lowercase"`
	Project                 string `json:"project"`
	Dataset                 string `json:"dataset"`
	HeadersFormatting       int64  `json:"headers_formatting"`
	StackID                 int64  `json:"stack"`
	AuthID                  int64  `json:"auth"`
}

func (c *Client) CreateDestination(destinationTypeId int, req *DestinationConfig) (*DestinationResponse, error) {
	r, _ := url.JoinPath("target-types", strconv.Itoa(destinationTypeId), "targets", "/")
	p, _ := url.Parse(r)

	resp, err := Create[DestinationConfig, DestinationResponse](c, p, req, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) ReadDestination(destinationTypeId, destinationId int) (*DestinationResponse, error) {
	r, _ := url.JoinPath("target-types", strconv.Itoa(destinationTypeId), "targets", strconv.Itoa(destinationId), "/")
	p, _ := url.Parse(r)

	resp, err := Read[DestinationResponse](c, p, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) UpdateDestination(destinationTypeId, destinationId int, req *DestinationConfig) (*DestinationResponse, error) {
	r, _ := url.JoinPath("target-types", strconv.Itoa(destinationTypeId), "targets", strconv.Itoa(destinationId), "/")
	p, _ := url.Parse(r)

	resp, err := Update[DestinationConfig, DestinationResponse](c, p, req, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) DeleteDestination(destinationTypeId, destinationId int) (*DestinationResponse, error) {
	r, _ := url.JoinPath("target-types", strconv.Itoa(destinationTypeId), "targets", strconv.Itoa(destinationId), "/")
	p, _ := url.Parse(r)

	resp, err := Delete[DestinationResponse](c, p, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
