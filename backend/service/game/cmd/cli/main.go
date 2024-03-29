package main

import (
	"fmt"
	"os"

	"gitlab.com/alienspaces/go-mud/backend/core/cli"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/dependencies"
	runner "gitlab.com/alienspaces/go-mud/backend/service/game/internal/runner/cli"
)

func main() {

	// Dependencies
	c, l, s, err := dependencies.Default()
	if err != nil {
		fmt.Printf("Failed default dependencies >%v<\n", err)
		os.Exit(0)
	}

	// Runner
	r, err := runner.NewRunner(c, l)
	if err != nil {
		fmt.Printf("Failed new runner >%v<\n", err)
		os.Exit(0)
	}

	// CLI
	cli, err := cli.NewCLI(c, l, s, r)
	if err != nil {
		fmt.Printf("Failed new cli >%v<\n", err)
		os.Exit(1)
	}

	args := make(map[string]interface{})

	err = cli.Run(args)
	if err != nil {
		fmt.Printf("Failed cli run >%v<\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
