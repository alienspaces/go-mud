package main

import (
	"fmt"
	"os"

	"gitlab.com/alienspaces/go-mud/backend/core/server"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/dependencies"
	runner "gitlab.com/alienspaces/go-mud/backend/service/game/internal/runner/server"
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

	svc, err := server.NewServer(c, l, s, r)
	if err != nil {
		fmt.Printf("Failed new server >%v<\n", err)
		os.Exit(0)
	}

	args := make(map[string]interface{})

	err = svc.Run(args)
	if err != nil {
		fmt.Printf("Failed server run >%v<\n", err)
		os.Exit(0)
	}

	os.Exit(1)
}
