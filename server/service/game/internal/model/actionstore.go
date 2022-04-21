package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetActionRecs -
func (m *Model) GetActionRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.Action, error) {

	l := m.Logger("GetActionRecs")

	l.Debug("Getting dungeon action records params >%s<", params)

	r := m.ActionRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetActionRec -
func (m *Model) GetActionRec(recID string, forUpdate bool) (*record.Action, error) {

	l := m.Logger("GetActionRec")

	l.Debug("Getting dungeon action rec ID >%s<", recID)

	r := m.ActionRepository()

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

// CreateActionRec -
func (m *Model) CreateActionRec(rec *record.Action) error {

	l := m.Logger("CreateActionRec")

	l.Debug("Creating dungeon action rec >%#v<", rec)

	r := m.ActionRepository()

	err := m.ValidateActionRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateActionRec -
func (m *Model) UpdateActionRec(rec *record.Action) error {

	l := m.Logger("UpdateActionRec")

	l.Debug("Updating dungeon action rec >%#v<", rec)

	r := m.ActionRepository()

	err := m.ValidateActionRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteActionRec -
func (m *Model) DeleteActionRec(recID string) error {

	l := m.Logger("DeleteActionRec")

	l.Debug("Deleting dungeon action rec ID >%s<", recID)

	r := m.ActionRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteActionRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveActionRec -
func (m *Model) RemoveActionRec(recID string) error {

	l := m.Logger("RemoveActionRec")

	l.Debug("Removing dungeon action rec ID >%s<", recID)

	r := m.ActionRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteActionRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
