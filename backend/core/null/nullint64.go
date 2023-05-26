package null

import (
	"database/sql"
)

func NullInt64FromInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}

func NullInt64FromInt64Ptr(i *int64) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{
			Int64: 0,
			Valid: false,
		}
	}

	return sql.NullInt64{
		Int64: *i,
		Valid: true,
	}
}

func NullInt64ToInt64(n sql.NullInt64) int64 {
	if !n.Valid {
		return 0
	}
	return n.Int64
}

func NullInt64ToInt64Ptr(n sql.NullInt64) *int64 {
	if !n.Valid {
		return nil
	}
	return &n.Int64
}

func NullInt64IsValid(n sql.NullInt64) bool {
	return n.Valid
}
