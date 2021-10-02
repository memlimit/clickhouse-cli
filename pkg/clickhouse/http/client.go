package http

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type CompressType string

const (
	Gzip CompressType = "gzip"
)

type Client struct {
	baseURL *url.URL

	username string
	password string

	client       *http.Client
	compressType CompressType
}

func New(addr, username, password string, compress CompressType) (*Client, error) {
	baseURL, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	return &Client{
		baseURL:      baseURL,
		username:     username,
		password:     password,
		client:       http.DefaultClient,
		compressType: compress,
	}, nil
}

// NewRequest creates an new request.
func (c *Client) NewRequest(url, query, method string) (*http.Request, error) {
	if !strings.HasSuffix(c.baseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.baseURL)
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	switch c.compressType {
	case Gzip:
		req.Header.Set("Accept-Encoding", string(Gzip))
	}

	q := req.URL.Query()

	q.Add("user", c.username)
	q.Add("password", c.password)
	q.Add("query", query)
	q.Add("enable_http_compression", "1")
	q.Add("default_format", "PrettyCompact")

	req.URL.RawQuery = q.Encode()

	return req, nil
}

func (c *Client) BareDo(ctx context.Context, req *http.Request) (*http.Response, error) {
	if ctx == nil {
		return nil, errors.New("ctx must be not nil")
	}

	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}

	return resp, err
}

func (c *Client) Do(ctx context.Context, req *http.Request) (string, error) {
	resp, err := c.BareDo(ctx, req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var reader io.ReadCloser

	switch CompressType(resp.Header.Get("Content-Encoding")) {
	case Gzip:
		reader, err = gzip.NewReader(resp.Body)
	default:
		reader = resp.Body
	}

	defer reader.Close()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	if err = c.handleError(resp); err != nil {
		return "", errors.New(string(data))
	}

	return string(data), err
}

func (c *Client) Query(ctx context.Context, query string) (string, error) {
	req, err := c.NewRequest(c.baseURL.String(), query, "POST")
	if err != nil {
		return "", err
	}

	data, err := c.Do(ctx, req)
	if err != nil {
		return "", err
	}

	return data, nil
}

func (c *Client) handleError(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	return errors.New("clickhouse returned error")
}
