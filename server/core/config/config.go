// Package config provides methods for managing configuration
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"gitlab.com/alienspaces/go-boilerplate/server/core/type/configurer"
)

// Item defines a valid environment variable and whether it is required
type Item struct {
	Key      string
	Required bool
}

// Config defines a container of items and corresponding values
type Config struct {
	Items  []*Item
	Values map[string]string
}

var _ configurer.Configurer = &Config{}

// NewConfig creates a new environment object
func NewConfig(items []Item, dotEnv bool) (*Config, error) {

	e := Config{
		Values: make(map[string]string),
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
	dir := os.Getenv("APP_SERVER_HOME")
	if dir == "" {
		dir, err = os.Getwd()
		os.Setenv("APP_SERVER_HOME", dir)
	}
	err = e.Add("APP_SERVER_HOME", true)
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

	for _, item := range e.Items {
		if item.Key == key {
			value = e.Values[key]
			return value
		}
	}

	return ""
}

// Set a config item value
func (e *Config) Set(key string, value string) {

	for _, item := range e.Items {
		if item.Key == key {
			e.Values[key] = value
			return
		}
	}
}

// Add will add a new config item
func (e *Config) Add(key string, required bool) (err error) {

	item := Item{
		Key:      key,
		Required: required,
	}

	e.Items = append(e.Items, &item)

	err = e.sourceItem(&item)
	if err != nil {
		return err
	}

	err = e.checkItem(&item)
	if err != nil {
		return err
	}

	return nil
}

// Verify checks whether the provided items have values set
func (e *Config) Verify(items []Item) (err error) {

	for _, item := range items {
		err = e.checkItem(&item)
		if err != nil {
			return err
		}
	}

	return nil
}

// sourceItem - sources and sets a config item value
func (e *Config) sourceItem(item *Item) error {

	value := os.Getenv(item.Key)
	e.Set(item.Key, value)

	return nil
}

// checkItem - checks a config item
func (e *Config) checkItem(item *Item) error {

	if item.Required && e.Values[item.Key] == "" {
		return fmt.Errorf("Failed checking env value >%s<", item.Key)
	}

	return nil
}
