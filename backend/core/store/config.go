package store

import (
	"fmt"
	"strconv"

	"gitlab.com/alienspaces/go-mud/backend/core/config"
	"gitlab.com/alienspaces/go-mud/backend/core/type/storer"
)

const (
	defaultMaxOpenConnections int = 5
	defaultMaxIdleConnections int = 2
	defaultMaxIdleTimeMins    int = 2
)

func (s *Store) SetConnectionConfig(connectionConfig *storer.ConnectionConfig) error {
	l := s.Log

	if connectionConfig.MaxOpenConnections == 0 {
		connectionConfig.MaxOpenConnections = defaultMaxOpenConnections
	}
	if connectionConfig.MaxIdleConnections == 0 {
		connectionConfig.MaxIdleConnections = defaultMaxIdleConnections
	}
	if connectionConfig.MaxIdleTimeMins == 0 {
		connectionConfig.MaxIdleTimeMins = defaultMaxIdleTimeMins
	}

	err := validateConnectionConfig(connectionConfig)
	if err != nil {
		l.Warn("failed validating connection configuration >%v<", err)
		return err
	}

	s.connectionConfig = connectionConfig

	return nil
}

func (s *Store) GetConnectionConfig() (*storer.ConnectionConfig, error) {
	l := s.Log
	c := s.Config

	if s.connectionConfig != nil {
		return s.connectionConfig, nil
	}

	dbMaxOpenConnectionsStr := c.Get(config.AppServerDBMaxOpenConnections)
	dbMaxOpenConnections, err := strconv.Atoi(dbMaxOpenConnectionsStr)
	if err != nil {
		errMsg := "APP_SERVER_DB_MAX_OPEN_CONNECTIONS is not int"
		l.Warn(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	dbMaxIdleConnectionsStr := c.Get(config.AppServerDBMaxIdleConnections)
	dbMaxIdleConnections, err := strconv.Atoi(dbMaxIdleConnectionsStr)
	if err != nil {
		errMsg := "APP_SERVER_DB_MAX_IDLE_CONNECTIONS is not int"
		l.Warn(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	dbMaxIdleTimeMinsStr := c.Get(config.AppServerDBMaxIdleTimeMins)
	dbMaxIdleTimeMins, err := strconv.Atoi(dbMaxIdleTimeMinsStr)
	if err != nil {
		errMsg := "APP_SERVER_DB_MAX_IDLE_TIME_MINS is not int"
		l.Warn(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	connectionConfig := storer.ConnectionConfig{
		Host:               c.Get(config.AppServerDBHost),
		Port:               c.Get(config.AppServerDBPort),
		Database:           c.Get(config.AppServerDBName),
		User:               c.Get(config.AppServerDBUser),
		Password:           c.Get(config.AppServerDBPassword),
		MaxOpenConnections: dbMaxOpenConnections,
		MaxIdleConnections: dbMaxIdleConnections,
		MaxIdleTimeMins:    dbMaxIdleTimeMins,
	}

	if connectionConfig.MaxOpenConnections == 0 {
		connectionConfig.MaxOpenConnections = defaultMaxOpenConnections
	}
	if connectionConfig.MaxIdleConnections == 0 {
		connectionConfig.MaxIdleConnections = defaultMaxIdleConnections
	}
	if connectionConfig.MaxIdleTimeMins == 0 {
		connectionConfig.MaxIdleTimeMins = defaultMaxIdleTimeMins
	}

	err = validateConnectionConfig(&connectionConfig)
	if err != nil {
		l.Warn("failed validating connection configuration >%v<", err)
		return nil, err
	}

	return &connectionConfig, nil
}

func validateConnectionConfig(c *storer.ConnectionConfig) error {

	if c.Host == "" {
		errMsg := "configuration missing Host, cannot connect"
		return fmt.Errorf(errMsg)
	}

	if c.Port == "" {
		errMsg := "configuration missing Port, cannot connect"
		return fmt.Errorf(errMsg)
	}

	if c.Database == "" {
		errMsg := "configuration missing Database, cannot connect"
		return fmt.Errorf(errMsg)
	}

	if c.User == "" {
		errMsg := "configuration missing User, cannot connect"
		return fmt.Errorf(errMsg)
	}

	if c.Password == "" {
		errMsg := "configuration missing Password, cannot connect"
		return fmt.Errorf(errMsg)
	}

	if c.MaxOpenConnections == 0 {
		errMsg := "configuration missing MaxOpenConnections, cannot connect"
		return fmt.Errorf(errMsg)
	}

	if c.MaxIdleConnections == 0 {
		errMsg := "configuration missing MaxIdleConnections, cannot connect"
		return fmt.Errorf(errMsg)
	}

	if c.MaxIdleTimeMins == 0 {
		errMsg := "configuration missing MaxIdleTimeMins, cannot connect"
		return fmt.Errorf(errMsg)
	}

	return nil
}
