package store

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // blank import intended

	"gitlab.com/alienspaces/go-mud/server/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
)

// getPostgresDB -
func getPostgresDB(c configurer.Configurer, l logger.Logger) (*sqlx.DB, error) {

	dbHost := c.Get("APP_SERVER_DB_HOST")
	dbPort := c.Get("APP_SERVER_DB_PORT")
	dbName := c.Get("APP_SERVER_DB_NAME")
	dbUser := c.Get("APP_SERVER_DB_USER")
	dbPass := c.Get("APP_SERVER_DB_PASSWORD")

	if dbHost == "" {
		l.Warn("missing APP_SERVER_DB_HOST, cannot connect")
		return nil, fmt.Errorf("missing APP_SERVER_DB_HOST, cannot connect")
	}

	if dbPort == "" {
		l.Warn("missing APP_SERVER_DB_PORT, cannot connect")
		return nil, fmt.Errorf("missing APP_SERVER_DB_PORT, cannot connect")
	}

	if dbName == "" {
		l.Warn("missing APP_SERVER_DB_NAME, cannot connect")
		return nil, fmt.Errorf("missing APP_SERVER_DB_NAME, cannot connect")
	}

	if dbUser == "" {
		l.Warn("missing APP_SERVER_DB_USER, cannot connect")
		return nil, fmt.Errorf("missing APP_SERVER_DB_USER, cannot connect")
	}

	if dbPass == "" {
		l.Warn("missing APP_SERVER_DB_PASS, cannot connect")
		return nil, fmt.Errorf("missing APP_SERVER_DB_PASS, cannot connect")
	}

	cs := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbUser, dbPass, dbName, dbHost, dbPort)

	l.Info("Connect string %s", cs)

	d, err := sqlx.Connect("postgres", cs)
	if err != nil {
		l.Warn("Failed to db connect >%v<", err)
		return nil, err
	}

	return d, nil
}
