package nullbool

import (
	"database/sql"
)

func FromBool(b bool) sql.NullBool {
	return sql.NullBool{
		Bool:  b,
		Valid: true,
	}
}

func FromBoolPtr(b *bool) sql.NullBool {
	if b == nil {
		return sql.NullBool{
			Bool:  false,
			Valid: false,
		}
	}

	return sql.NullBool{
		Bool:  *b,
		Valid: true,
	}
}

func ToBoolPtr(ns sql.NullBool) *bool {
	if !ns.Valid {
		return nil
	}

	return &ns.Bool
}

func IsValid(ns sql.NullBool) bool {
	return ns.Valid
}
