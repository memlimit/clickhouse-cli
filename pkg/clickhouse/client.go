package clickhouse

import "context"

type Response struct {
	QueryID string
	Data string
}

type Client interface {
	Query(ctx context.Context, query string) (string, error)
}
