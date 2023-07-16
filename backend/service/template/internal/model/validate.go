package model

import (
	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/service/template/internal/record"
)

// ValidateTemplateRec - validates creating and updating a template record
func (m *Model) ValidateTemplateRec(rec *record.Template) error {

	if rec.ID != "" {
		if !m.IsUUID(rec.ID) {
			return coreerror.NewInvalidDataError("ID >%s< is not a valid UUID", rec.ID)
		}
	}

	return nil
}

// ValidateDeleteTemplateRec - validates it is okay to delete a template record
func (m *Model) ValidateDeleteTemplateRec(recID string) error {

	if !m.IsUUID(recID) {
		return coreerror.NewInvalidDataError("ID >%s< is not a valid UUID", recID)
	}

	return nil
}
