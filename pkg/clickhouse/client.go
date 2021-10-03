package clickhouse

import "context"

// Response of clickhouse
type Response struct {
	QueryID string
	Data    string
}

// Client - interface implementing Clickhouse
type Client interface {
	Query(ctx context.Context, query string) (string, error)
}
