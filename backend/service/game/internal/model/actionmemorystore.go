package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetActionMemoryRecs -
func (m *Model) GetActionMemoryRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.ActionMemory, error) {

	l := m.Logger("GetActionMemoryRecs")

	l.Debug("Getting monster instance records params >%s<", params)

	r := m.ActionMemoryRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetActionMemoryRec -
func (m *Model) GetActionMemoryRec(recID string, forUpdate bool) (*record.ActionMemory, error) {

	l := m.Logger("GetActionMemoryRec")

	l.Debug("Getting monster instance record ID >%s<", recID)

	r := m.ActionMemoryRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return nil, fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	rec, err := r.GetOne(recID, forUpdate)
	if err == sql.ErrNoRows {
		l.Warn("No record found ID >%s<", recID)
		return nil, nil
	}

	return rec, err
}

// CreateActionMemoryRec -
func (m *Model) CreateActionMemoryRec(rec *record.ActionMemory) error {

	l := m.Logger("CreateActionMemoryRec")

	l.Debug("Creating monster record >%#v<", rec)

	r := m.ActionMemoryRepository()

	err := m.validateActionMemoryRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateActionMemoryRec -
func (m *Model) UpdateActionMemoryRec(rec *record.ActionMemory) error {

	l := m.Logger("UpdateActionMemoryRec")

	l.Debug("Updating monster record >%#v<", rec)

	r := m.ActionMemoryRepository()

	err := m.validateActionMemoryRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteActionMemoryRec -
func (m *Model) DeleteActionMemoryRec(recID string) error {

	l := m.Logger("DeleteActionMemoryRec")

	l.Debug("Deleting monster rec ID >%s<", recID)

	r := m.ActionMemoryRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteActionMemoryRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveActionMemoryRec -
func (m *Model) RemoveActionMemoryRec(recID string) error {

	l := m.Logger("RemoveActionMemoryRec")

	l.Debug("Removing monster rec ID >%s<", recID)

	r := m.ActionMemoryRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteActionMemoryRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
