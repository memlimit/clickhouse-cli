package cli

import (
	"context"
)

// GetCurrentDB from clickhouse
func (c *CLI) GetCurrentDB(ctx context.Context) string {
	db, _ := c.client.Query(ctx, "SELECT DATABASE() FORMAT TabSeparated")

	return db
}
