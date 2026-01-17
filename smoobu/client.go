package smoobu

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	BaseURL  string
	APIToken string
	http     *http.Client
}

func NewClient(baseURL, apiToken string) *Client {
	return &Client{
		BaseURL:  baseURL,
		APIToken: apiToken,
		http:     &http.Client{},
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Add("Api-Key", c.APIToken)
	return c.http.Do(req)
}

func (c *Client) Get(path string) (*http.Response, error) {
	if !strings.HasPrefix(path, "/") {
		return nil, fmt.Errorf("path must start with '/'")
	}

	req, err := http.NewRequest("GET", c.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) Post(path string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.BaseURL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) Delete(path string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", c.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
