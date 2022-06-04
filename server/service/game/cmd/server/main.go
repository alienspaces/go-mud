package main

import (
	"fmt"
	"os"

	"gitlab.com/alienspaces/go-mud/server/core/server"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/dependencies"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/server/runner"
)

func main() {

	c, l, s, err := dependencies.Default()
	if err != nil {
		fmt.Printf("Failed default dependencies >%v<", err)
		os.Exit(0)
	}

	r, err := runner.NewRunner(c, l)
	if err != nil {
		fmt.Printf("Failed new runner >%v<", err)
		os.Exit(0)
	}

	svc, err := server.NewServer(c, l, s, r)
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
