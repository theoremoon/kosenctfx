package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/pkg/errors"
)

type Client struct {
	URL    *url.URL
	APIKey string

	HTTPClient *http.Client
}

func NewClient(url *url.URL, apiKey string) *Client {
	c := &Client{
		URL:        url,
		APIKey:     apiKey,
		HTTPClient: &http.Client{},
	}

	return c
}

func (c *Client) NewRequest(ctx context.Context, method, p string, body io.Reader) (*http.Request, error) {
	u := *c.URL
	u.Path = path.Join(c.URL.Path, p)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, errors.Wrap(err, "NewRequest")
	}
	req = req.WithContext(ctx)

	req.SetBasicAuth("api", c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "kosenctfx-agent")

	return req, nil
}

func (c *Client) Get(ctx context.Context, path string, payload *url.Values) (*http.Response, error) {
	req, err := c.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	if payload != nil {
		req.URL.RawQuery = payload.Encode()
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		res.Body.Close()
		return nil, fmt.Errorf("status: %s", res.Status)
	}
	return res, nil
}

func (c *Client) Post(ctx context.Context, path string, payload interface{}) (*http.Response, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := c.NewRequest(ctx, "POST", path, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		res.Body.Close()
		return nil, fmt.Errorf("status: %s", res.Status)
	}
	return res, nil
}

func DecodeBody(res *http.Response, out interface{}) error {
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	return decoder.Decode(out)
}
