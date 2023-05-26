package null

import (
	"database/sql"
)

func NullBoolFromBool(b bool) sql.NullBool {
	return sql.NullBool{
		Bool:  b,
		Valid: true,
	}
}

func NullBoolFromBoolPtr(b *bool) sql.NullBool {
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

func NullBoolToBool(ns sql.NullBool) bool {
	if !ns.Valid {
		return false
	}

	return ns.Bool
}

func NullBoolToBoolPtr(ns sql.NullBool) *bool {
	if !ns.Valid {
		return nil
	}

	return &ns.Bool
}

func NullBoolIsValid(ns sql.NullBool) bool {
	return ns.Valid
}
