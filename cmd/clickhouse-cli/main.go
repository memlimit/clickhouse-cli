package main

import (
	"context"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/memlimit/clickhouse-cli/cli"
	"github.com/memlimit/clickhouse-cli/cli/completer"
	"github.com/memlimit/clickhouse-cli/pkg/clickhouse/http"
	"os"
)

func main() {
	client, err := http.New("https://gh-api.clickhouse.tech/", "play", "")
	if err != nil {
		return 
	}

	chVersion, err := client.Query(context.Background(), "SELECT version() FORMAT TabSeparated;")
	if err != nil {
		fmt.Println("Failed to connect to ClickHouse")
		os.Exit(0)
	}

	fmt.Printf("Connected to ClickHouse server version %s", chVersion)

	c := cli.New(client)
	complete := completer.New()

	p := prompt.New(
		c.Executor,
		complete.Complete,
		prompt.OptionTitle("clickhouse-cli: cli for ClickHouse."),
	)

	p.Run()
}
