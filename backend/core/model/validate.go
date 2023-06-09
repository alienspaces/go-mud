package model

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/core/null"
)

// IsUUID - tests whether provided string is a valid UUID
func IsUUID(s string) bool {
	if _, err := uuid.Parse(s); err != nil {
		return false
	}

	return true
}

// This remains here until all references instead call the above package function
func (m *Model) IsUUID(s string) bool {
	if !IsUUID(s) {
		m.Log.Warn("UUID >%s< is not valid", s)
		return false
	}

	return true
}

func (m *Model) ValidateStringField(field string, fieldName string) error {
	if field == "" {
		err := coreerror.NewInvalidDataError("%s should not be empty >%s<", fieldName, field)
		m.Log.Warn(err.Error())
		return err
	}

	return nil
}

func (m *Model) ValidateNullStringField(field sql.NullString, fieldName string) error {
	if !null.NullStringIsValid(field) {
		err := coreerror.NewInvalidDataError("%s should not be empty >%s<", fieldName, field.String)
		m.Log.Warn(err.Error())
		return err
	}

	return nil
}

func (m *Model) ValidateNullBoolField(field sql.NullBool, fieldName string) error {
	if !null.NullBoolIsValid(field) {
		err := coreerror.NewInvalidDataError("%s should not be empty", fieldName)
		m.Log.Warn(err.Error())
		return err
	}

	return nil
}

func (m *Model) ValidateStringArrayField(field pq.StringArray, fieldName string) error {
	if len(field) == 0 {
		err := coreerror.NewInvalidDataError("%s should not be empty >%#v<", fieldName, field)
		m.Log.Warn(err.Error())
		return err
	}

	return nil
}
