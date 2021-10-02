package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/memlimit/clickhouse-cli/cli"
	"github.com/memlimit/clickhouse-cli/cli/completer"
	"github.com/memlimit/clickhouse-cli/cli/history"
	"github.com/memlimit/clickhouse-cli/pkg/clickhouse/http"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(0)
		return
	}
}

func run() error {
	client, err := http.New("https://gh-api.clickhouse.tech/", "play", "")
	if err != nil {
		return err
	}

	chVersion, err := client.Query(context.Background(), "SELECT version() FORMAT TabSeparated;")
	if err != nil {
		return errors.New("failed to connect to ClickHouse")
	}

	fmt.Printf("Connected to ClickHouse server version %s\n", chVersion)

	homeDirPath, err := os.UserHomeDir()
	h, err := history.New(homeDirPath + "/.clickhouse-client-history")
	if err != nil {
		return err
	}

	uh, err := h.Read()
	if err != nil {
		return err
	}

	c := cli.New(client, h, true)
	complete := completer.New()

	p := prompt.New(
		c.Executor,
		complete.Complete,
		prompt.OptionTitle("clickhouse-cli: cli for ClickHouse."),
		prompt.OptionHistory(h.RowsToStrArr(uh)),
		prompt.OptionPrefix(c.GetCurrentDB(context.Background()) + " :) "),
		prompt.OptionLivePrefix(c.GetLivePrefixState),
		prompt.OptionPrefixTextColor(prompt.White),
	)

	p.Run()

	return nil
}
