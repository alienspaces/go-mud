package model

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/core/nullbool"
	"gitlab.com/alienspaces/go-mud/backend/core/nullstring"
)

// IsUUID - tests whether provided string is a valid UUID
func (m *Model) IsUUID(s string) bool {

	_, err := uuid.Parse(s)
	if err != nil {
		m.Log.Warn("UUID >%s< is not valid >%v<", s, err)
		return false
	}
	return true
}

func (m *Model) ValidateStringField(field string, fieldName string) error {
	if field == "" {
		errMsg := fmt.Sprintf("%s should not be empty >%s<", fieldName, field)
		m.Log.Warn("failed validating %s >%s<", fieldName, errMsg)
		return coreerror.NewInvalidError(fieldName, errMsg)
	}

	return nil
}

func (m *Model) ValidateNullStringField(field sql.NullString, fieldName string) error {
	if !nullstring.IsValid(field) {
		errMsg := fmt.Sprintf("%s should not be empty >%s<", fieldName, field.String)
		m.Log.Warn("failed validating %s >%s<", fieldName, errMsg)
		return coreerror.NewInvalidError(fieldName, errMsg)
	}

	return nil
}

func (m *Model) ValidateNullBoolField(field sql.NullBool, fieldName string) error {
	if !nullbool.IsValid(field) {
		errMsg := fmt.Sprintf("%s should not be empty", fieldName)
		m.Log.Warn("failed validating %s >%s<", fieldName, errMsg)
		return coreerror.NewInvalidError(fieldName, errMsg)
	}

	return nil
}

func (m *Model) ValidateStringArrayField(field pq.StringArray, fieldName string) error {
	if len(field) == 0 {
		errMsg := fmt.Sprintf("%s should not be empty >%#v<", fieldName, field)
		m.Log.Warn("failed validating %s >%s<", fieldName, errMsg)
		return coreerror.NewInvalidError(fieldName, errMsg)
	}

	return nil
}
