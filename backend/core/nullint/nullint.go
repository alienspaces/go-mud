package nullint

import (
	"database/sql"
)

func FromInt16(i int16) sql.NullInt16 {
	return sql.NullInt16{
		Int16: i,
		Valid: true,
	}
}

func FromInt16Ptr(i *int16) sql.NullInt16 {
	if i == nil {
		return sql.NullInt16{
			Int16: 0,
			Valid: false,
		}
	}

	return sql.NullInt16{
		Int16: *i,
		Valid: true,
	}
}

func ToInt16(ni sql.NullInt16) int16 {
	if !ni.Valid {
		return 0
	}
	return ni.Int16
}

func ToInt16Ptr(ni sql.NullInt16) *int16 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int16
}

func IsValid(ni sql.NullInt16) bool {
	return ni.Valid
}
