package dependencies

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/core/log"
	"gitlab.com/alienspaces/go-mud/backend/core/store"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/config"
)

func Default() (*config.Config, *log.Log, *store.Store, error) {

	// Configurer
	c, err := config.NewConfig()
	if err != nil {
		fmt.Printf("failed new config >%v<", err)
		return nil, nil, nil, err
	}

	// Logger
	l, err := log.NewLogger(c)
	if err != nil {
		fmt.Printf("failed new logger >%v<", err)
		return nil, nil, nil, err
	}

	// Storer
	s, err := store.NewStore(c, l)
	if err != nil {
		fmt.Printf("failed new store >%v<", err)
		return nil, nil, nil, err
	}

	return c, l, s, nil
}
