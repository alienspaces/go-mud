package model

import (
	"database/sql"
	"fmt"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetActionObjectRecs -
func (m *Model) GetActionObjectRecs(opts *coresql.Options) ([]*record.ActionObject, error) {

	l := m.loggerWithContext("GetActionObjectRecs")

	l.Debug("Getting dungeon action object records opts >%#v<", opts)

	r := m.ActionObjectRepository()

	return r.GetMany(opts)
}

// GetActionObjectRec -
func (m *Model) GetActionObjectRec(recID string, lock *coresql.Lock) (*record.ActionObject, error) {

	l := m.loggerWithContext("GetActionObjectRec")

	l.Debug("Getting dungeon action object rec ID >%s<", recID)

	r := m.ActionObjectRepository()

	if !m.IsUUID(recID) {
		return nil, fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	rec, err := r.GetOne(recID, lock)
	if err == sql.ErrNoRows {
		l.Warn("No record found ID >%s<", recID)
		return nil, nil
	}

	return rec, err
}

// CreateActionObjectRec -
func (m *Model) CreateActionObjectRec(rec *record.ActionObject) error {

	l := m.loggerWithContext("CreateActionObjectRec")

	l.Debug("Creating dungeon action object record >%#v<", rec)

	r := m.ActionObjectRepository()

	err := m.validateActionObjectRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateActionObjectRec -
func (m *Model) UpdateActionObjectRec(rec *record.ActionObject) error {

	l := m.loggerWithContext("UpdateActionObjectRec")

	l.Debug("Updating dungeon action object record >%#v<", rec)

	r := m.ActionObjectRepository()

	err := m.validateActionObjectRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteActionObjectRec -
func (m *Model) DeleteActionObjectRec(recID string) error {

	l := m.loggerWithContext("DeleteActionObjectRec")

	l.Debug("Deleting dungeon action object rec ID >%s<", recID)

	r := m.ActionObjectRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteActionObjectRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveActionObjectRec -
func (m *Model) RemoveActionObjectRec(recID string) error {

	l := m.loggerWithContext("RemoveActionObjectRec")

	l.Debug("Removing dungeon action object rec ID >%s<", recID)

	r := m.ActionObjectRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteActionObjectRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
