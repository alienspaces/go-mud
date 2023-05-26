// Package config provides methods for managing configuration
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
)

// Config defines a container of Items and corresponding Values. Items specifies whether the Item.Key is required.
type Config struct {
	Required map[string]bool
	Values   map[string]string
}

// NewConfig creates a new environment object
func NewConfig(items []Item, dotEnv bool) (*Config, error) {

	e := Config{
		Required: make(map[string]bool),
		Values:   make(map[string]string),
	}

	err := e.Init(items, dotEnv)
	if err != nil {
		return nil, fmt.Errorf("NewConfig failed init >%v<", err)
	}

	return &e, nil
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

// MustGet returns a non-empty config item value, otherwise an error
func (e *Config) MustGet(key string) (string, error) {
	v := e.Values[key]
	if v == "" {
		return "", fmt.Errorf("config key >%s< must not be empty", key)
	}
	return v, nil
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
		return fmt.Errorf("failed checking required env value >%s<", item.Key)
	}

	return nil
}

func (e *Config) Clone() configurer.Configurer {
	cfg := Config{
		Required: make(map[string]bool),
		Values:   make(map[string]string),
	}

	for k, v := range e.Required {
		cfg.Required[k] = v
	}

	for k, v := range e.Values {
		cfg.Values[k] = v
	}

	return &cfg
}
