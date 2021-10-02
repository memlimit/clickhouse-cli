package cli

import (
	"context"
)

func (c *CLI) GetCurrentDB(ctx context.Context) string {
	db, _ := c.client.Query(ctx, "SELECT DATABASE() FORMAT TabSeparated")

	return db
}
