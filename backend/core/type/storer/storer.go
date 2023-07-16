package storer

import "github.com/jmoiron/sqlx"

// Storer -
type Storer interface {
	GetDb() (*sqlx.DB, error)
	GetTx() (*sqlx.Tx, error)
}
