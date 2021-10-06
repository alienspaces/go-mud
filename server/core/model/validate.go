package model

import (
	"github.com/google/uuid"
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
