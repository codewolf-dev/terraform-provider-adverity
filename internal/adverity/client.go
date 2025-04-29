// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package adverity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"slices"
	"strings"
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

// buildURL constructs a full URL for a given path.
func (c *Client) buildURL(path string) string {
	return strings.TrimRight(c.Endpoint.String(), "/") + "/" + strings.TrimLeft(path, "/")
}

func (c *Client) create(path string, payload io.Reader) ([]byte, error) {
	res, err := c.doRequest(http.MethodPost, path, payload, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) read(path string) ([]byte, error) {
	res, err := c.doRequest(http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) update(path string, payload io.Reader) ([]byte, error) {
	res, err := c.doRequest(http.MethodPatch, path, payload, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) delete(path string) ([]byte, error) {
	res, err := c.doRequest(http.MethodDelete, path, nil, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) doRequest(method, path string, payload io.Reader, authToken *string) ([]byte, error) {
	token := c.Token

	if authToken != nil {
		token = *authToken
	}

	req, err := http.NewRequest(method, c.buildURL(path), payload)
	if err != nil {
		return nil, err
	}

	allowedMethods := []string{http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}
	if !slices.Contains(allowedMethods, method) {
		return nil, fmt.Errorf("unsupported method: %s, allowed: %v", method, allowedMethods)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	expectedStatusCodes := []int{http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent}
	if !slices.Contains(expectedStatusCodes, resp.StatusCode) {
		return nil, fmt.Errorf("status: %d, body: %s, expected: %v", resp.StatusCode, body, expectedStatusCodes)
	}

	return body, nil
}

func Create[ReqT any, RespT any](c *Client, path string, resource *ReqT) (*RespT, error) {
	payload, err := json.Marshal(resource)
	if err != nil {
		return nil, err
	}

	body, err := c.create(path, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	resp := new(RespT)
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func Read[RespT any](c *Client, path string) (*RespT, error) {
	body, err := c.read(path)
	if err != nil {
		return nil, err
	}

	resp := new(RespT)
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func Update[ReqT any, RespT any](c *Client, path string, resource *ReqT) (*RespT, error) {
	payload, err := json.Marshal(resource)
	if err != nil {
		return nil, err
	}

	body, err := c.update(path, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	resp := new(RespT)
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func Delete[RespT any](c *Client, path string) (*RespT, error) {
	body, err := c.delete(path)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, nil
	}

	resp := new(RespT)
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
