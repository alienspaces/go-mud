package store

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // blank import intended

	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/storer"
)

// connectPostgresDB -
func connectPostgresDB(l logger.Logger, c *storer.ConnectionConfig) (*sqlx.DB, error) {

	if c == nil {
		return nil, fmt.Errorf("missing configurer, cannot connect to postgres database")
	}
	if l == nil {
		return nil, fmt.Errorf("missing logger, cannot connect to postgres database")
	}

	dbHost := c.Host
	if dbHost == "" {
		errMsg := "missing APP_SERVER_DB_HOST, cannot connect"
		l.Warn(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	dbPort := c.Port
	if dbPort == "" {
		errMsg := "missing APP_SERVER_DB_PORT, cannot connect"
		l.Warn(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	dbName := c.Database
	if dbName == "" {
		errMsg := "missing APP_SERVER_DB_NAME, cannot connect"
		l.Warn(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	dbUser := c.User
	if dbUser == "" {
		errMsg := "missing APP_SERVER_DB_USER, cannot connect"
		l.Warn(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	dbPass := c.Password
	if dbPass == "" {
		errMsg := "missing APP_SERVER_DB_PASSWORD, cannot connect"
		l.Warn(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	cs := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbUser, dbPass, dbName, dbHost, dbPort)
	l.Info("Connect string >%s<", fmt.Sprintf("user=%s password=******* dbname=%s host=%s port=%s sslmode=disable", dbUser, dbName, dbHost, dbPort))

	dbMaxOpenConnections := c.MaxOpenConnections
	if dbMaxOpenConnections == 0 {
		errMsg := "missing APP_SERVER_DB_MAX_OPEN_CONNECTIONS, cannot connect"
		l.Warn(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	dbMaxIdleConnections := c.MaxIdleConnections
	if dbMaxIdleConnections == 0 {
		errMsg := "missing APP_SERVER_DB_MAX_IDLE_CONNECTIONS, cannot connect"
		l.Warn(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	dbMaxIdleTimeMins := c.MaxIdleTimeMins
	if dbMaxIdleTimeMins == 0 {
		errMsg := "missing APP_SERVER_DB_MAX_IDLE_TIME_MINS, cannot connect"
		l.Warn(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	d, err := sqlx.Connect("postgres", cs)
	if err != nil {
		l.Warn("failed to db connect >%v<", err)
		return nil, err
	}

	d.SetMaxOpenConns(dbMaxOpenConnections)
	d.SetMaxIdleConns(dbMaxIdleConnections)
	d.SetConnMaxIdleTime(time.Duration(dbMaxIdleTimeMins) * time.Minute)

	return d, nil
}
