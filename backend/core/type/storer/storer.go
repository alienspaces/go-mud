package storer

import "github.com/jmoiron/sqlx"

type ConnectionConfig struct {
	Host               string
	Port               string
	User               string
	Database           string
	Password           string
	MaxOpenConnections int
	MaxIdleConnections int
	MaxIdleTimeMins    int
}

// Storer -
type Storer interface {
	GetDb() (*sqlx.DB, error)
	GetTx() (*sqlx.Tx, error)
	GetConnectionConfig() (*ConnectionConfig, error)
	SetConnectionConfig(*ConnectionConfig) error
}
