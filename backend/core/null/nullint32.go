package null

import (
	"database/sql"
)

func NullInt32FromInt32(s int32) sql.NullInt32 {
	return sql.NullInt32{
		Int32: s,
		Valid: true,
	}
}

func NullInt32FromInt32Ptr(s *int32) sql.NullInt32 {
	if s == nil {
		return sql.NullInt32{
			Int32: 0,
			Valid: false,
		}
	}

	return sql.NullInt32{
		Int32: *s,
		Valid: true,
	}
}

func NullInt32ToInt32(ns sql.NullInt32) int32 {
	if !ns.Valid {
		return 0
	}
	return ns.Int32
}

func NullInt32ToInt32Ptr(ns sql.NullInt32) *int32 {
	if !ns.Valid {
		return nil
	}
	return &ns.Int32
}

func NullInt32IsValid(ns sql.NullInt32) bool {
	return ns.Valid
}
