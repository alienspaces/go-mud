package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	"gitlab.com/alienspaces/go-mud/backend/core/nulltime"
)

// NOTE:
// Use sql.NullXxx types when the underlying database column:
// - has a NOT NULL check constraint;
// - does not otherwise have any other CHECK constraints;
// - is not an eval type;
// - is not an uuid type;
// - does not have a foreign key constraint; and
// - you don't want it to default to Go's default value for the property type.

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
	return nulltime.FromTime(timestamp())
}

// NewProcessedAt -
func NewProcessedAt() sql.NullTime {
	return nulltime.FromTime(timestamp())
}

// NewFailedAt -
func NewFailedAt() sql.NullTime {
	return nulltime.FromTime(timestamp())
}

// NewDeletedAt -
func NewDeletedAt() sql.NullTime {
	return nulltime.FromTime(timestamp())
}

func timestamp() time.Time {
	return time.Now().UTC()
}
