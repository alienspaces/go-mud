package store

import (
	"fmt"
	"strconv"

	"gitlab.com/alienspaces/go-mud/backend/core/config"
)

const (
	defaultMaxOpenConnections int = 5
	defaultMaxIdleConnections int = 2
	defaultMaxIdleTimeMins    int = 2
)

func (s *Store) SetConnectionConfig(config *Config) error {
	l := s.Log

	if config.MaxOpenConnections == 0 {
		config.MaxOpenConnections = defaultMaxOpenConnections
	}
	if config.MaxIdleConnections == 0 {
		config.MaxIdleConnections = defaultMaxIdleConnections
	}
	if config.MaxIdleTimeMins == 0 {
		config.MaxIdleTimeMins = defaultMaxIdleTimeMins
	}

	err := validateConnectionConfig(config)
	if err != nil {
		l.Warn("failed validating connection configuration >%v<", err)
		return err
	}

	s.config = config

	return nil
}

func (s *Store) GetConnectionConfig() (*Config, error) {
	l := s.Log
	c := s.Config

	if s.config != nil {
		return s.config, nil
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

	config := Config{
		Host:               c.Get(config.AppServerDBHost),
		Port:               c.Get(config.AppServerDBPort),
		Database:           c.Get(config.AppServerDBName),
		User:               c.Get(config.AppServerDBUser),
		Password:           c.Get(config.AppServerDBPassword),
		MaxOpenConnections: dbMaxOpenConnections,
		MaxIdleConnections: dbMaxIdleConnections,
		MaxIdleTimeMins:    dbMaxIdleTimeMins,
	}

	if config.MaxOpenConnections == 0 {
		config.MaxOpenConnections = defaultMaxOpenConnections
	}
	if config.MaxIdleConnections == 0 {
		config.MaxIdleConnections = defaultMaxIdleConnections
	}
	if config.MaxIdleTimeMins == 0 {
		config.MaxIdleTimeMins = defaultMaxIdleTimeMins
	}

	err = validateConnectionConfig(&config)
	if err != nil {
		l.Warn("failed validating connection configuration >%v<", err)
		return nil, err
	}

	return &config, nil
}

func validateConnectionConfig(c *Config) error {

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
