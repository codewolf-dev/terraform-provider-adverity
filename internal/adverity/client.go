// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package adverity

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"slices"
	"time"
)

// Client -
type Client struct {
	HTTPClient *http.Client
	Endpoint   *url.URL
	Token      string
}

// NewClient -
func NewClient(instanceUrl, authToken *string) (*Client, error) {
	baseUrl := *instanceUrl
	restPath := "/api/"

	apiEndpoint, err := url.ParseRequestURI(baseUrl + restPath)
	if err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return nil, err
	}

	c := Client{
		HTTPClient: &http.Client{Timeout: 30 * time.Second, Jar: jar},
		Endpoint:   apiEndpoint,
		Token:      *authToken,
	}

	return &c, nil
}

func (client *Client) create(path string, payload *bytes.Reader) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, client.Endpoint.String()+path, io.NopCloser(payload))
	if err != nil {
		return nil, err
	}

	res, err := client.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (client *Client) read(path string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, client.Endpoint.String()+path, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (client *Client) update(path string, payload *bytes.Reader) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPatch, client.Endpoint.String()+path, io.NopCloser(payload))
	if err != nil {
		return nil, err
	}

	res, err := client.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (client *Client) delete(path string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodDelete, client.Endpoint.String()+path, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (client *Client) doRequest(req *http.Request, authToken *string) ([]byte, error) {
	token := client.Token

	if authToken != nil {
		token = *authToken
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if !slices.Contains([]int{http.StatusOK, http.StatusCreated, http.StatusAccepted}, res.StatusCode) {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, nil
}
