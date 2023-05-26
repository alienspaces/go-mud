package null

import (
	"database/sql"
	"time"
)

func NullTimeFromTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{}
	}
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}

func NullTimeToTime(n sql.NullTime) time.Time {
	if !n.Valid {
		return time.Time{}
	}
	return n.Time
}

func NullTimeIsValid(nt sql.NullTime) bool {
	return nt.Valid && !nt.Time.IsZero()
}
