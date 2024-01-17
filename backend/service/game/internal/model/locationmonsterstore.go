package model

import (
	"database/sql"
	"fmt"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetLocationMonsterRecs -
func (m *Model) GetLocationMonsterRecs(opts *coresql.Options) ([]*record.LocationMonster, error) {

	l := m.loggerWithFunctionContext("GetLocationMonsterRecs")

	l.Debug("Getting dungeon location records opts >%#v<", opts)

	r := m.LocationMonsterRepository()

	return r.GetMany(opts)
}

// GetLocationMonsterRec -
func (m *Model) GetLocationMonsterRec(recID string, lock *coresql.Lock) (*record.LocationMonster, error) {

	l := m.loggerWithFunctionContext("GetLocationMonsterRec")

	l.Debug("Getting dungeon location rec ID >%s<", recID)

	r := m.LocationMonsterRepository()

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

// CreateLocationMonsterRec -
func (m *Model) CreateLocationMonsterRec(rec *record.LocationMonster) error {

	l := m.loggerWithFunctionContext("CreateLocationMonsterRec")

	l.Debug("Creating dungeon location record >%#v<", rec)

	r := m.LocationMonsterRepository()

	err := m.validateLocationMonsterRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateLocationMonsterRec -
func (m *Model) UpdateLocationMonsterRec(rec *record.LocationMonster) error {

	l := m.loggerWithFunctionContext("UpdateLocationMonsterRec")

	l.Debug("Updating dungeon location record >%#v<", rec)

	r := m.LocationMonsterRepository()

	err := m.validateLocationMonsterRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteLocationMonsterRec -
func (m *Model) DeleteLocationMonsterRec(recID string) error {

	l := m.loggerWithFunctionContext("DeleteLocationMonsterRec")

	l.Debug("Deleting dungeon location rec ID >%s<", recID)

	r := m.LocationMonsterRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteLocationMonsterRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveLocationMonsterRec -
func (m *Model) RemoveLocationMonsterRec(recID string) error {

	l := m.loggerWithFunctionContext("RemoveLocationMonsterRec")

	l.Debug("Removing dungeon location rec ID >%s<", recID)

	r := m.LocationMonsterRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteLocationMonsterRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
