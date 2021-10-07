package main

import (
	"context"
	"fmt"
	"os"

	"github.com/c-bata/go-prompt"

	"github.com/memlimit/clickhouse-cli/cli"
	"github.com/memlimit/clickhouse-cli/cli/completer"
	"github.com/memlimit/clickhouse-cli/cli/config"
	"github.com/memlimit/clickhouse-cli/cli/history"
	chHttp "github.com/memlimit/clickhouse-cli/pkg/clickhouse/http"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(0)
		return
	}
}

func run() error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	client, err := chHttp.New(cfg.HTTP.URL, cfg.Auth.UserName, cfg.Auth.Password, chHttp.CompressType(cfg.HTTP.Compress))
	if err != nil {
		return err
	}

	chVersion, err := client.Query(context.Background(), "SELECT version() FORMAT TabSeparated;")
	if err != nil {
		return fmt.Errorf("failed to connect to ClickHouse (%s)", err.Error())
	}

	fmt.Printf("Connected to ClickHouse server version %s\n", chVersion)

	h, uh, err := initHistory(cfg.CLI.HistoryPath)
	if err != nil {
		return err
	}

	c := cli.New(client, h, cfg.CLI.Multiline)
	complete := completer.New()

	p := prompt.New(
		c.Executor,
		complete.Complete,
		prompt.OptionTitle("clickhouse-cli: cli for ClickHouse."),
		prompt.OptionHistory(h.RowsToStrArr(uh)),
		prompt.OptionPrefix(c.GetCurrentDB(context.Background())+" :) "),
		prompt.OptionLivePrefix(c.GetLivePrefixState),
		prompt.OptionPrefixTextColor(prompt.White),
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.F3,
			Fn:  c.MultilineControl,
		}),
	)

	p.Run()

	return nil
}

func initHistory(path string) (*history.History, []*history.Row, error) {
	var historyPath string
	if path != "" {
		historyPath = path
	} else {
		home, _ := os.UserHomeDir()
		historyPath = home + "/.clickhouse-client-history"
	}

	h, err := history.New(historyPath)
	if err != nil {
		return nil, nil, err
	}

	uh, err := h.Read()
	if err != nil {
		return nil, nil, err
	}

	return h, uh, nil
}
