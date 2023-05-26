package null

import (
	"database/sql"
)

func NullInt16FromInt16(s int16) sql.NullInt16 {
	return sql.NullInt16{
		Int16: s,
		Valid: true,
	}
}

func NullInt16FromInt16Ptr(s *int16) sql.NullInt16 {
	if s == nil {
		return sql.NullInt16{
			Int16: 0,
			Valid: false,
		}
	}

	return sql.NullInt16{
		Int16: *s,
		Valid: true,
	}
}

func NullInt16ToInt16(ns sql.NullInt16) int16 {
	if !ns.Valid {
		return 0
	}
	return ns.Int16
}

func NullInt16ToInt16Ptr(ns sql.NullInt16) *int16 {
	if !ns.Valid {
		return nil
	}
	return &ns.Int16
}

func NullInt16IsValid(ns sql.NullInt16) bool {
	return ns.Valid
}
