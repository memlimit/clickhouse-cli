package cli

import (
	"context"
	"fmt"
	"github.com/memlimit/clickhouse-cli/cli/history"
	"os"
	"time"
)

var exitCodes = [...]string{
	"exit",
	"quit",
	"logout",
	"учше",
	"йгше",
	"дщпщге",
	"exit;",
	"quit;",
	"logout;",
	"учшеж",
	"йгшеж",
	"дщпщгеж",
	"q",
	"й",
	"Q",
	":q",
	"Й",
	"Жй",
}

func (c *CLI) Executor(s string) {
	if err := c.history.Write(&history.Row{
		CreatedAt: time.Now(),
		Query:     s,
	}); err != nil {
		fmt.Println(err)
	}

	for _, code := range exitCodes {
		if s == code {
			fmt.Println("Bye.")
			os.Exit(0)
			return
		}
	}

	data, _ := c.client.Query(context.Background(), s)
	fmt.Println(data)
}
