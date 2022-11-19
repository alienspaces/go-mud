package store

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // blank import intended

	"gitlab.com/alienspaces/go-mud/server/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
)

// getPostgresDB -
func getPostgresDB(c configurer.Configurer, l logger.Logger) (*sqlx.DB, error) {

	dbHost := c.Get("APP_SERVER_DB_HOST")
	if dbHost == "" {
		errMsg := "missing APP_SERVER_DB_HOST, cannot connect"
		l.Warn(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	dbPort := c.Get("APP_SERVER_DB_PORT")
	if dbPort == "" {
		errMsg := "missing APP_SERVER_DB_PORT, cannot connect"
		l.Warn(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	dbName := c.Get("APP_SERVER_DB_NAME")
	if dbName == "" {
		errMsg := "missing APP_SERVER_DB_NAME, cannot connect"
		l.Warn(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	dbUser := c.Get("APP_SERVER_DB_USER")
	if dbUser == "" {
		errMsg := "missing APP_SERVER_DB_USER, cannot connect"
		l.Warn(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	dbPass := c.Get("APP_SERVER_DB_PASSWORD")
	if dbPass == "" {
		errMsg := "missing APP_SERVER_DB_PASSWORD, cannot connect"
		l.Warn(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	cs := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbUser, dbPass, dbName, dbHost, dbPort)
	l.Info("Connect string %s", cs)

	dbMaxOpenConnectionsStr := c.Get("APP_SERVER_DB_MAX_OPEN_CONNECTIONS")
	if dbPass == "" {
		errMsg := "missing APP_SERVER_DB_MAX_OPEN_CONNECTIONS, cannot connect"
		l.Warn(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	dbMaxOpenConnections, err := strconv.Atoi(dbMaxOpenConnectionsStr)
	if err != nil {
		errMsg := "APP_SERVER_DB_MAX_OPEN_CONNECTIONS is not int"
		l.Warn(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	dbMaxIdleConnectionsStr := c.Get("APP_SERVER_DB_MAX_IDLE_CONNECTIONS")
	if dbPass == "" {
		errMsg := "missing APP_SERVER_DB_MAX_IDLE_CONNECTIONS, cannot connect"
		l.Warn(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	dbMaxIdleConnections, err := strconv.Atoi(dbMaxIdleConnectionsStr)
	if err != nil {
		errMsg := "APP_SERVER_DB_MAX_IDLE_CONNECTIONS is not int"
		l.Warn(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	dbMaxIdleTimeMinsStr := c.Get("APP_SERVER_DB_MAX_IDLE_TIME_MINS")
	if dbPass == "" {
		errMsg := "missing APP_SERVER_DB_MAX_IDLE_TIME_MINS, cannot connect"
		l.Warn(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	dbMaxIdleTimeMins, err := strconv.Atoi(dbMaxIdleTimeMinsStr)
	if err != nil {
		errMsg := "APP_SERVER_DB_MAX_IDLE_TIME_MINS is not int"
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
