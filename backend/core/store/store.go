package store

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/core/type/storer"

	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
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
	config   *Config
}

type Config struct {
	Host               string
	Port               string
	User               string
	Database           string
	Password           string
	MaxOpenConnections int
	MaxIdleConnections int
	MaxIdleTimeMins    int
}

var _ storer.Storer = &Store{}

// NewStore -
func NewStore(c configurer.Configurer, l logger.Logger) (*Store, error) {

	dt := ""
	if c != nil {
		dt = c.Get("APP_SERVER_DATABASE")
	}

	if dt == "" {
		l.Debug("defaulting to postgres")
		dt = DBPostgres
	}

	s := Store{
		Log:      l,
		Config:   c,
		Database: dt,
	}

	return &s, nil
}

// GetDb -
func (s *Store) GetDb() (*sqlx.DB, error) {
	l := s.Log

	// Return existing database handle if we have one
	if s.DB != nil {
		return s.DB, nil
	}

	var err error
	config := s.config
	if config == nil {
		config, err = s.GetConnectionConfig()
		if err != nil {
			s.Log.Warn("failed to get connection config >%v<", err)
			return nil, err
		}
	}

	if s.Database == DBPostgres {
		l.Debug("connecting to postgres")
		conn, err := connectPostgresDB(s.Log, config)
		if err != nil {
			l.Warn("failed getting postgres DB handle >%v<", err)
			return nil, err
		}
		s.DB = conn
		return conn, nil
	}

	return nil, fmt.Errorf("unsupported database")
}

// GetTx -
func (s *Store) GetTx() (*sqlx.Tx, error) {
	l := s.Log

	if s.DB == nil {
		_, err := s.GetDb()
		if err != nil {
			l.Warn("failed getting postgres DB handle >%v<", err)
			return nil, err
		}
	}

	return s.DB.Beginx()
}
