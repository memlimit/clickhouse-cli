package cli

import "github.com/memlimit/clickhouse-cli/pkg/clickhouse"

type CLI struct {
	client clickhouse.Client
}

func New(client clickhouse.Client) *CLI {
	return &CLI{
		client: client,
	}
}