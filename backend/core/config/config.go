// Package config provides methods for managing configuration
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
)

const (
	AppServerHome                 string = "APP_SERVER_HOME"
	AppServerEnv                  string = "APP_SERVER_ENV"
	AppServerPort                 string = "APP_SERVER_PORT"
	AppServerHost                 string = "APP_SERVER_HOST"
	AppServerLogLevel             string = "APP_SERVER_LOG_LEVEL"
	AppServerLogPretty            string = "APP_SERVER_LOG_PRETTY"
	AppServerDbHost               string = "APP_SERVER_DB_HOST"
	AppServerDbPort               string = "APP_SERVER_DB_PORT"
	AppServerDbName               string = "APP_SERVER_DB_NAME"
	AppServerDbUser               string = "APP_SERVER_DB_USER"
	AppServerDbPassword           string = "APP_SERVER_DB_PASSWORD"
	AppServerDbMaxOpenConnections string = "APP_SERVER_DB_MAX_OPEN_CONNECTIONS"
	AppServerDbMaxIdleConnections string = "APP_SERVER_DB_MAX_IDLE_CONNECTIONS"
	AppServerDbMaxIdleTimeMins    string = "APP_SERVER_DB_MAX_IDLE_TIME_MINS"
	AppServerSchemaPath           string = "APP_SERVER_SCHEMA_PATH"
	AppServerJwtSigningKey        string = "APP_SERVER_JWT_SIGNING_KEY"
)

// Config defines a container of Items and corresponding Values. Items specifies whether the Item.Key is required.
type Config struct {
	Required map[string]bool
	Values   map[string]string
}

var _ configurer.Configurer = &Config{}

// NewConfig creates a new environment object
func NewConfig(items []Item, dotEnv bool) (*Config, error) {

	c := Config{
		Required: make(map[string]bool),
		Values:   make(map[string]string),
	}

	err := c.Init(items, dotEnv)
	if err != nil {
		return nil, fmt.Errorf("NewConfig failed init >%v<", err)
	}

	return &c, nil
}

// Init initialises and checks environment values
func (e *Config) Init(items []Item, dotEnv bool) (err error) {

	// app home
	dir := os.Getenv(AppServerHome)
	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return err
		}
		err := os.Setenv(AppServerHome, dir)
		if err != nil {
			return err
		}
	}
	err = e.Add(AppServerHome, true)
	if err != nil {
		return err
	}

	if dotEnv {
		envFile := fmt.Sprintf("%s/%s", dir, ".env")
		err = godotenv.Load(envFile)
		if err != nil {
			return err
		}
	}

	for _, item := range items {
		err = e.Add(item.Key, item.Required)
		if err != nil {
			return err
		}
	}

	return nil
}

// Get returns a config item value
func (e *Config) Get(key string) (value string) {
	return e.Values[key]
}

// Set a config item value
func (e *Config) Set(key string, value string) {
	e.Values[key] = value
}

// Add will add a new config item
func (e *Config) Add(key string, required bool) (err error) {

	item := Item{
		Key:      key,
		Required: required,
	}

	e.Required[key] = required

	err = e.sourceItem(item)
	if err != nil {
		return err
	}

	err = e.checkItem(item)
	if err != nil {
		return err
	}

	return nil
}

func (e *Config) GetAll() map[string]string {
	return e.Values
}

// sourceItem - sources and sets a config item value
func (e *Config) sourceItem(item Item) error {

	value := os.Getenv(item.Key)
	e.Set(item.Key, value)

	return nil
}

// checkItem - checks a config item
func (e *Config) checkItem(item Item) error {

	if item.Required && e.Values[item.Key] == "" {
		return fmt.Errorf("failed checking env value >%s<", item.Key)
	}

	return nil
}
