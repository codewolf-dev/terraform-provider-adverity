// Copyright codewolf.dev 2025, 0
// SPDX-License-Identifier: MPL-2.0

package adverity

import (
	"net/url"
	"strconv"
)

type Schedule struct {
	CronPreset         *string `json:"cron_preset,omitempty"`
	CronType           *string `json:"cron_type,omitempty"`
	CronInterval       *int64  `json:"cron_interval,omitempty"`
	CronIntervalStart  *int64  `json:"cron_interval_start,omitempty"`
	CronStartOfDay     *string `json:"cron_start_of_day,omitempty"`
	TimeRangePreset    *int64  `json:"time_range_preset,omitempty"`
	DeltaType          *int64  `json:"delta_type,omitempty"`
	DeltaInterval      *int64  `json:"delta_interval,omitempty"`
	DeltaIntervalStart *int64  `json:"delta_interval_start,omitempty"`
	DeltaStartOfDay    *string `json:"delta_start_of_day,omitempty"`
	FixedStart         *string `json:"fixed_start,omitempty"`
	FixedEnd           *string `json:"fixed_end,omitempty"`
	OffsetDays         *int64  `json:"offset_days,omitempty"`
	NotBeforeDate      *string `json:"not_before_date,omitempty"`
	NotBeforeTime      *string `json:"not_before_time,omitempty"`
}

type DatastreamScheduleConfig struct {
	Schedules *[]Schedule `json:"schedules,omitempty"`
	Enabled   *bool       `json:"enabled,omitempty"`
	DataType  *string     `json:"datatype,omitempty"`
}

type DatastreamUpdateConfig struct {
	Name                *string      `json:"name,omitempty"`
	Description         *string      `json:"description,omitempty"`
	StackID             *int64       `json:"stack,omitempty"`
	AuthID              *int64       `json:"auth,omitempty"`
	RetentionType       *int64       `json:"retention_type,omitempty"`
	RetentionNumber     *int64       `json:"retention_number,omitempty"`
	OverwriteKeyColumns *bool        `json:"overwrite_key_columns,omitempty"`
	OverwriteDatastream *bool        `json:"overwrite_datastream,omitempty"`
	OverwriteFileName   *bool        `json:"overwrite_filename,omitempty"`
	IsInsightsMediaplan *bool        `json:"is_insights_mediaplan,omitempty"`
	ManageExtractNames  *bool        `json:"manage_extract_names,omitempty"`
	ExtractNameKeys     *string      `json:"extract_name_keys,omitempty"`
	Parameters          *[]Parameter `json:"-"`
}

func (c *DatastreamUpdateConfig) MarshalJSON() ([]byte, error) {
	return FlattenedMarshal(c, c.Parameters)
}

type DatastreamCreateConfig struct {
	Name                *string      `json:"name,omitempty"`
	Description         *string      `json:"description,omitempty"`
	StackID             *int64       `json:"stack,omitempty"`
	AuthID              *int64       `json:"auth,omitempty"`
	DataType            *string      `json:"datatype,omitempty"`
	RetentionType       *int64       `json:"retention_type,omitempty"`
	RetentionNumber     *int64       `json:"retention_number,omitempty"`
	OverwriteKeyColumns *bool        `json:"overwrite_key_columns,omitempty"`
	OverwriteDatastream *bool        `json:"overwrite_datastream,omitempty"`
	OverwriteFileName   *bool        `json:"overwrite_filename,omitempty"`
	IsInsightsMediaplan *bool        `json:"is_insights_mediaplan,omitempty"`
	ManageExtractNames  *bool        `json:"manage_extract_names,omitempty"`
	ExtractNameKeys     *string      `json:"extract_name_keys,omitempty"`
	Parameters          *[]Parameter `json:"-"`
	Schedules           *[]Schedule  `json:"schedules,omitempty"`
	Enabled             *bool        `json:"enabled,omitempty"`
}

func (c *DatastreamCreateConfig) MarshalJSON() ([]byte, error) {
	return FlattenedMarshal(c, c.Parameters)
}

type DatastreamResponse struct {
	ID                  int64      `json:"id"`
	DataType            string     `json:"datatype"`
	Creator             string     `json:"creator"`
	DatastreamTypeID    int64      `json:"datastream_type_id"`
	AbsoluteURL         string     `json:"absolute_url"`
	Created             string     `json:"created"`
	Updated             string     `json:"updated"`
	Slug                string     `json:"slug"`
	Name                string     `json:"name"`
	Description         string     `json:"description"`
	Enabled             bool       `json:"enabled"`
	AuthID              int64      `json:"auth"`
	Frequency           string     `json:"frequency"`
	LastFetch           string     `json:"last_fetch"`
	NextRun             string     `json:"next_run"`
	OverviewURL         string     `json:"overview_url"`
	StackID             int64      `json:"stack_id"`
	Schedules           []Schedule `json:"schedules"`
	RetentionType       int64      `json:"retention_type"`
	RetentionNumber     int64      `json:"retention_number"`
	OverwriteKeyColumns bool       `json:"overwrite_key_columns"`
	OverwriteDatastream bool       `json:"overwrite_datastream"`
	OverwriteFileName   bool       `json:"overwrite_filename"`
	IsInsightsMediaplan bool       `json:"is_insights_mediaplan"`
	ManageExtractNames  bool       `json:"manage_extract_names"`
	ExtractNameKeys     string     `json:"extract_name_keys"`
}

func (c *Client) CreateDatastream(datastreamTypeId int, req *DatastreamCreateConfig) (*DatastreamResponse, error) {
	r, _ := url.JoinPath("datastream-types", strconv.Itoa(datastreamTypeId), "datastreams", "/")
	p, _ := url.Parse(r)

	return Create[DatastreamCreateConfig, DatastreamResponse](c, p, req, nil)
}

func (c *Client) ReadDatastream(datastreamTypeId, datastreamId int) (*DatastreamResponse, error) {
	r, _ := url.JoinPath("datastream-types", strconv.Itoa(datastreamTypeId), "datastreams", strconv.Itoa(datastreamId), "/")
	p, _ := url.Parse(r)

	return Read[DatastreamResponse](c, p, nil)
}

func (c *Client) UpdateDatastream(datastreamTypeId, datastreamId int, req *DatastreamUpdateConfig) (*DatastreamResponse, error) {
	r, _ := url.JoinPath("datastream-types", strconv.Itoa(datastreamTypeId), "datastreams", strconv.Itoa(datastreamId), "/")
	p, _ := url.Parse(r)

	return Update[DatastreamUpdateConfig, DatastreamResponse](c, p, req, nil)
}

func (c *Client) DeleteDatastream(datastreamTypeId, datastreamId int) (*DatastreamResponse, error) {
	r, _ := url.JoinPath("datastream-types", strconv.Itoa(datastreamTypeId), "datastreams", strconv.Itoa(datastreamId), "/")
	p, _ := url.Parse(r)

	return Delete[DatastreamResponse](c, p, nil)
}

func (c *Client) UpdateDatastreamSchedule(datastreamId int, req *DatastreamScheduleConfig) (*DatastreamResponse, error) {
	r, _ := url.JoinPath("datastreams", strconv.Itoa(datastreamId), "/")
	p, _ := url.Parse(r)

	return Update[DatastreamScheduleConfig, DatastreamResponse](c, p, req, nil)
}
