package main

import (
	"fmt"
	"os"

	"gitlab.com/alienspaces/go-mud/backend/core/cli"

	"gitlab.com/alienspaces/go-mud/backend/service/template/internal/cli/runner"
	"gitlab.com/alienspaces/go-mud/backend/service/template/internal/dependencies"
)

func main() {

	c, l, s, err := dependencies.Default()
	if err != nil {
		fmt.Printf("Failed default dependencies >%v<\n", err)
		os.Exit(0)
	}

	r, err := runner.NewRunner(c, l)
	if err != nil {
		fmt.Printf("Failed new runner >%v<\n", err)
		os.Exit(0)
	}

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
