package main

import (
	"fmt"
	"os"

	"gitlab.com/alienspaces/go-mud/server/core/server"
	"gitlab.com/alienspaces/go-mud/server/service/template/internal/dependencies"
	"gitlab.com/alienspaces/go-mud/server/service/template/internal/server/runner"
)

func main() {

	// Dependencies
	c, l, s, m, err := dependencies.Default()
	if err != nil {
		fmt.Printf("Failed default dependencies >%v<", err)
		os.Exit(0)
	}

	r := runner.NewRunner()

	svc, err := server.NewServer(c, l, s, m, r)
	if err != nil {
		fmt.Printf("Failed new server >%v<", err)
		os.Exit(0)
	}

	args := make(map[string]interface{})

	err = svc.Run(args)
	if err != nil {
		fmt.Printf("Failed server run >%v<", err)
		os.Exit(0)
	}

	os.Exit(1)
}
