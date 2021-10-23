package store

import "database/sql"

func NullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{
			String: "",
			Valid:  false,
		}
	}
	return sql.NullString{
		String: value,
		Valid:  true,
	}
}

func NullInt(value int) sql.NullInt32 {
	if value == 0 {
		return sql.NullInt32{
			Valid: false,
		}
	}
	return sql.NullInt32{
		Int32: int32(value),
		Valid: true,
	}
}
