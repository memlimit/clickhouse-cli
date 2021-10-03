package grpc

import (
	"context"

	"google.golang.org/grpc"
)

// CompressType - no/gzip
type CompressType string

const (
	No   CompressType = ""     //nolint:revive
	Gzip CompressType = "gzip" //nolint:revive
)

// Client gRPC
type Client struct {
	chClient clickHouseClient

	username string
	password string

	compressType CompressType
}

// New - returns new client via gRPC
func New(addr, username, password string, compress CompressType) (*Client, error) {
	var opts []grpc.DialOption
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := clickHouseClient{conn}

	return &Client{
		username:     username,
		password:     password,
		compressType: compress,
		chClient:     client,
	}, nil
}

// Query with gRPC
func (c *Client) Query(ctx context.Context, query string) (string, error) {
	var cp Compression
	cp.Level = CompressionLevel_COMPRESSION_HIGH

	switch c.compressType {
	case Gzip:
		cp.Algorithm = CompressionAlgorithm_GZIP
	case No:
		cp.Algorithm = CompressionAlgorithm_NO_COMPRESSION
	default:
		cp.Algorithm = CompressionAlgorithm_NO_COMPRESSION
	}

	q := QueryInfo{
		Query:             query,
		UserName:          c.username,
		Password:          c.password,
		ResultCompression: &cp,
	}

	result, err := c.chClient.ExecuteQuery(ctx, &q, nil)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}
