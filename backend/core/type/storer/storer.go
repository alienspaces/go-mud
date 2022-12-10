package storer

import "github.com/jmoiron/sqlx"

// Storer -
type Storer interface {
	Init() error
	GetDb() (*sqlx.DB, error)
	GetTx() (*sqlx.Tx, error)
}
