// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package adverity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"slices"
	"time"
)

// Client holds http.Client, endpoint and token.
type Client struct {
	httpClient *http.Client
	endpoint   *url.URL
	token      string
}

// NewClient constructs a new Client with given endpoint and token.
func NewClient(instanceUrl, authToken *string) (*Client, error) {
	baseUrl, err := url.Parse(*instanceUrl)
	if err != nil {
		return nil, err
	}

	restPath := "api"
	apiEndpoint := baseUrl.JoinPath(restPath, "/") //  Needed to make buildURL work properly (see https://datatracker.ietf.org/doc/html/rfc3986#section-5.2.3)

	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return nil, err
	}

	log.Printf("Building API client for %s", apiEndpoint.String())
	c := Client{
		httpClient: &http.Client{Timeout: 30 * time.Second, Jar: jar},
		endpoint:   apiEndpoint,
		token:      *authToken,
	}

	return &c, nil
}

// buildURL constructs a full URL for a given path.
func (c *Client) buildURL(path *url.URL) *url.URL {
	return c.endpoint.ResolveReference(path)
}

func (c *Client) Create(path *url.URL, payload io.Reader, query *url.Values) ([]byte, error) {
	return c.doRequest(http.MethodPost, path, payload, query, nil)
}

func (c *Client) Read(path *url.URL, query *url.Values) ([]byte, error) {
	return c.doRequest(http.MethodGet, path, nil, query, nil)
}

func (c *Client) Update(path *url.URL, payload io.Reader, query *url.Values) ([]byte, error) {
	return c.doRequest(http.MethodPatch, path, payload, query, nil)
}

func (c *Client) Delete(path *url.URL, query *url.Values) ([]byte, error) {
	return c.doRequest(http.MethodDelete, path, nil, query, nil)
}

//nolint:unparam // authToken may be used to overwrite a token set in Client
func (c *Client) doRequest(method string, path *url.URL, payload io.Reader, query *url.Values, authToken *string) ([]byte, error) {
	token := c.token

	if authToken != nil {
		token = *authToken
	}

	// Build the resource URL
	u := c.buildURL(path)

	// Add query parameters to the URL
	if query != nil {
		u.RawQuery = query.Encode()
	}

	// Create the request
	req, err := http.NewRequest(method, u.String(), payload)
	if err != nil {
		return nil, err
	}

	// Check allowed methods
	allowedMethods := []string{http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}
	if !slices.Contains(allowedMethods, method) {
		return nil, fmt.Errorf("unsupported method: %s, allowed: %v", method, allowedMethods)
	}

	// Add headers (e.g., auth token)
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Handle HTTP errors
	expectedStatusCodes := []int{http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent}
	if !slices.Contains(expectedStatusCodes, resp.StatusCode) {
		return nil, fmt.Errorf("status: %d, body: %s, expected: %v", resp.StatusCode, body, expectedStatusCodes)
	}

	return body, nil
}

func Create[ReqT any, RespT any](c *Client, path *url.URL, resource *ReqT, query *url.Values) (*RespT, error) {
	return execute[ReqT, RespT](c, http.MethodPost, path, resource, query)
}

func Read[RespT any](c *Client, path *url.URL, query *url.Values) (*RespT, error) {
	return execute[any, RespT](c, http.MethodGet, path, nil, query)
}

func Update[ReqT any, RespT any](c *Client, path *url.URL, resource *ReqT, query *url.Values) (*RespT, error) {
	return execute[ReqT, RespT](c, http.MethodPatch, path, resource, query)
}

func Delete[RespT any](c *Client, path *url.URL, query *url.Values) (*RespT, error) {
	return execute[any, RespT](c, http.MethodDelete, path, nil, query)
}

func execute[ReqT any, RespT any](c *Client, method string, path *url.URL, resource *ReqT, query *url.Values) (*RespT, error) {
	log.Printf("%s %s", method, path.String())

	var r io.Reader

	if resource != nil {
		payload, err := json.Marshal(resource)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal %T: %w", *resource, err)
		}
		log.Printf("request payload: %s", payload)
		r = bytes.NewReader(payload)
	}

	body, err := c.doRequest(method, path, r, query, nil)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, nil
	}

	resp := new(RespT)
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal into %T: %w", *resp, err)
	}
	log.Printf("response body: %s", body)

	return resp, nil
}
