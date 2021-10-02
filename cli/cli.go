package cli

import (
	"github.com/memlimit/clickhouse-cli/cli/history"
	"github.com/memlimit/clickhouse-cli/pkg/clickhouse"
)

// CLI object of cli :)
type CLI struct {
	client  clickhouse.Client
	history *history.History

	Multiline               bool
	isMultilineInputStarted bool
	query                   string
}

// New - returns CLI object
func New(client clickhouse.Client, history *history.History, multiline bool) *CLI {
	return &CLI{
		client:    client,
		history:   history,
		Multiline: multiline,

		isMultilineInputStarted: false,
	}
}
