package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetLocationMonsterRecs -
func (m *Model) GetLocationMonsterRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.LocationMonster, error) {

	l := m.Logger("GetLocationMonsterRecs")

	l.Debug("Getting dungeon location records params >%s<", params)

	r := m.LocationMonsterRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetLocationMonsterRec -
func (m *Model) GetLocationMonsterRec(recID string, forUpdate bool) (*record.LocationMonster, error) {

	l := m.Logger("GetLocationMonsterRec")

	l.Debug("Getting dungeon location rec ID >%s<", recID)

	r := m.LocationMonsterRepository()

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

// CreateLocationMonsterRec -
func (m *Model) CreateLocationMonsterRec(rec *record.LocationMonster) error {

	l := m.Logger("CreateLocationMonsterRec")

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

	l := m.Logger("UpdateLocationMonsterRec")

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

	l := m.Logger("DeleteLocationMonsterRec")

	l.Debug("Deleting dungeon location rec ID >%s<", recID)

	r := m.LocationMonsterRepository()

	// validate UUID
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

	l := m.Logger("RemoveLocationMonsterRec")

	l.Debug("Removing dungeon location rec ID >%s<", recID)

	r := m.LocationMonsterRepository()

	// validate UUID
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
