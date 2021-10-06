package store

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-boilerplate/server/core/type/configurer"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
)

const (
	// DBPostgres -
	DBPostgres string = "postgres"
)

// Store -
type Store struct {
	Log      logger.Logger
	Config   configurer.Configurer
	Database string
	DB       *sqlx.DB
}

// NewStore -
func NewStore(c configurer.Configurer, l logger.Logger) (*Store, error) {

	dt := c.Get("APP_SERVER_DATABASE")
	if dt == "" {
		l.Debug("Defaulting to postgres")
		dt = DBPostgres
	}

	s := Store{
		Log:      l,
		Config:   c,
		Database: dt,
	}

	return &s, nil
}

// Init sets up a database connection if one doesn't already exist
func (s *Store) Init() error {

	// Get a database connection
	c, err := s.GetDb()
	if err != nil {
		s.Log.Warn("Failed database init >%v<", err)
		return err
	}
	s.DB = c

	return nil
}

// GetDb -
func (s *Store) GetDb() (*sqlx.DB, error) {

	// Return existing database handle if we have one
	if s.DB != nil {
		return s.DB, nil
	}

	if s.Database == DBPostgres {
		s.Log.Debug("Connecting to postgres")
		c, err := getPostgresDB(s.Config, s.Log)
		if err != nil {
			s.Log.Warn("Failed getting postgres DB handle >%v<", err)
			return nil, err
		}
		s.DB = c
		return c, nil
	}

	return nil, fmt.Errorf("Unsupported database")
}

// GetTx -
func (s *Store) GetTx() (*sqlx.Tx, error) {

	if s.DB == nil {
		s.Log.Warn("Not connected")
		return nil, fmt.Errorf("Not Not connected")
	}

	return s.DB.Beginx()
}
