package http

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"github.com/memlimit/clickhouse-cli/pkg/clickhouse"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Client object
type Client struct {
	baseURL *url.URL

	username string
	password string

	client       *http.Client
	compressType clickhouse.CompressType
}

// New - returns client object
func New(addr, username, password string, compress clickhouse.CompressType) (*Client, error) {
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

	if c.compressType == clickhouse.Gzip {
		req.Header.Set("Accept-Encoding", string(clickhouse.Gzip))
	}

	q := req.URL.Query()

	q.Add("user", c.username)
	q.Add("password", c.password)
	q.Add("query", query)

	if c.compressType != clickhouse.No {
		q.Add("enable_http_compression", "1")
	}

	q.Add("default_format", "PrettyCompact")

	req.URL.RawQuery = q.Encode()

	return req, nil
}

// BareDo request with ctx
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

// Do http request
func (c *Client) Do(ctx context.Context, req *http.Request) (string, error) {
	resp, err := c.BareDo(ctx, req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var reader io.ReadCloser

	switch clickhouse.CompressType(resp.Header.Get("Content-Encoding")) {
	case clickhouse.Gzip:
		reader, err = gzip.NewReader(resp.Body) //nolint:ineffassign
	default:
		reader = resp.Body
	}
	if err != nil {
		return "", err
	}

	defer reader.Close()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	dataWithoutNewLine := strings.TrimSuffix(string(data), "\n")

	if err = c.handleError(resp); err != nil {
		return "", errors.New(dataWithoutNewLine)
	}

	return dataWithoutNewLine, err
}

// Query - build http request with query to clickhouse
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
