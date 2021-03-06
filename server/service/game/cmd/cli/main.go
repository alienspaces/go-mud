package main

import (
	"fmt"
	"os"

	"gitlab.com/alienspaces/go-mud/server/core/cli"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/cli/runner"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/dependencies"
)

func main() {

	// Dependencies
	c, l, s, m, err := dependencies.Default()
	if err != nil {
		fmt.Printf("Failed default dependencies >%v<", err)
		os.Exit(0)
	}

	// Runner
	r := runner.NewRunner()

	// CLI
	cli, err := cli.NewCLI(c, l, s, m, r)
	if err != nil {
		fmt.Printf("Failed new cli >%v<", err)
		os.Exit(1)
	}

	args := make(map[string]interface{})

	err = cli.Run(args)
	if err != nil {
		fmt.Printf("Failed cli run >%v<", err)
		os.Exit(1)
	}

	os.Exit(0)
}
