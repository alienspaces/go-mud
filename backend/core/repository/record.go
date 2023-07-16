package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/r3labs/diff/v3"
	"gitlab.com/alienspaces/go-mud/backend/core/null"
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

func (r *Record) clearID() {
	r.ID = ""
}
func (r *Record) clearTimestamps() {
	r.CreatedAt = time.Time{}
	r.UpdatedAt = null.NullTimeFromTime(time.Time{})
	r.DeletedAt = null.NullTimeFromTime(time.Time{})
}

type EqualityFlag string

const EqualityFlagExcludeID EqualityFlag = "eo-exclude-id"
const EqualityFlagExcludeTimestamps EqualityFlag = "flag-exclude-timestamps"

type EquatableRecord interface {
	clearID()
	clearTimestamps()
}

func RecordEqual(p, pp EquatableRecord, flags ...EqualityFlag) (bool, error) {
	changelog, err := RecordDiff(p, pp, flags...)
	return len(changelog) == 0, err
}

func RecordDiff(p, pp EquatableRecord, flags ...EqualityFlag) (diff.Changelog, error) {
	for idx := range flags {
		switch flags[idx] {
		case EqualityFlagExcludeTimestamps:
			p.clearTimestamps()
			pp.clearTimestamps()
		case EqualityFlagExcludeID:
			p.clearID()
			pp.clearID()
		}
	}
	return diff.Diff(p, pp)
}

// NewRecordID -
func NewRecordID() string {
	uuidByte, _ := uuid.NewRandom()
	uuidString := uuidByte.String()
	return uuidString
}

// NewRecordTimestamp -
func NewRecordTimestamp() time.Time {
	return timestamp()
}

// NewRecordNullTimestamp -
func NewRecordNullTimestamp() sql.NullTime {
	return null.NullTimeFromTime(timestamp())
}

func timestamp() time.Time {
	return time.Now().UTC()
}
