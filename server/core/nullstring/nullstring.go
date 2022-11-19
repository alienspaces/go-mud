package nullstring

import (
	"database/sql"
)

func FromString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func FromStringPtr(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{
			String: "",
			Valid:  false,
		}
	}

	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}

func ToString(ns sql.NullString) string {
	if !ns.Valid {
		return ""
	}
	return ns.String
}

func ToStringPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

func IsValid(ns sql.NullString) bool {
	return ns.Valid && ns.String != ""
}
