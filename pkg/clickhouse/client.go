package clickhouse

import "context"

// Response of clickhouse
type Response struct {
	QueryID string
	Data    string
}

// CompressType for types of compression :D
// e.g. gzip/deflate/xz and others.
type CompressType string

const (
	No   CompressType = ""     //nolint:revive
	Gzip CompressType = "gzip" //nolint:revive
)

// Client - interface implementing Clickhouse
type Client interface {
	Query(ctx context.Context, query string) (string, error)
}
