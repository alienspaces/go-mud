package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Record -
type Record struct {
	ID        string       `db:"id"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

// NewRecordID -
func NewRecordID() string {
	uuidByte, _ := uuid.NewRandom()
	uuidString := uuidByte.String()
	return uuidString
}

// NewCreatedAt -
func NewCreatedAt() time.Time {
	return timestamp()
}

// NewUpdatedAt -
func NewUpdatedAt() sql.NullTime {
	return NewNullTime(timestamp())
}

// NewDeletedAt -
func NewDeletedAt() sql.NullTime {
	return NewNullTime(timestamp())
}

func timestamp() time.Time {
	return time.Now().UTC()
}

// NewNullTime - converts time type to sql.NewNullTime type
func NewNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{}
	}
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}

// NewNullString - converts string type to sql.NullString type
func NewNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
