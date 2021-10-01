package cli

import (
	"github.com/memlimit/clickhouse-cli/cli/history"
	"github.com/memlimit/clickhouse-cli/pkg/clickhouse"
)

type CLI struct {
	client clickhouse.Client
	history *history.History
}

func New(client clickhouse.Client, history *history.History) *CLI {
	return &CLI{
		client: client,
		history: history,
	}
}