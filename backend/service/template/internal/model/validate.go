package model

import (
	"fmt"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/service/template/internal/record"
)

// ValidateTemplateRec - validates creating and updating a template record
func (m *Model) ValidateTemplateRec(rec *record.Template) error {

	// validate UUID
	if rec.ID != "" {
		if !m.IsUUID(rec.ID) {
			return coreerror.NewInvalidError("template", fmt.Sprintf("ID >%s< is not a valid UUID", rec.ID))
		}
	}

	return nil
}

// ValidateDeleteTemplateRec - validates it is okay to delete a template record
func (m *Model) ValidateDeleteTemplateRec(recID string) error {

	// validate UUID
	if !m.IsUUID(recID) {
		return coreerror.NewInvalidError("template", fmt.Sprintf("ID >%s< is not a valid UUID", recID))
	}

	return nil
}
