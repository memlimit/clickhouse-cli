package cli

import (
	"context"
	"fmt"
	"os"
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
