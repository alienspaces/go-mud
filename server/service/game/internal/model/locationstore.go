package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetLocationRecs -
func (m *Model) GetLocationRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.Location, error) {

	m.Log.Debug("Getting dungeon location records params >%s<", params)

	r := m.LocationRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetLocationRec -
func (m *Model) GetLocationRec(recID string, forUpdate bool) (*record.Location, error) {

	m.Log.Debug("Getting dungeon location rec ID >%s<", recID)

	r := m.LocationRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return nil, fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	rec, err := r.GetOne(recID, forUpdate)
	if err == sql.ErrNoRows {
		m.Log.Warn("No record found ID >%s<", recID)
		return nil, nil
	}

	return rec, err
}

// CreateLocationRec -
func (m *Model) CreateLocationRec(rec *record.Location) error {

	m.Log.Debug("Creating dungeon location rec >%#v<", rec)

	r := m.LocationRepository()

	err := m.ValidateLocationRec(rec)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateLocationRec -
func (m *Model) UpdateLocationRec(rec *record.Location) error {

	m.Log.Debug("Updating dungeon location rec >%#v<", rec)

	r := m.LocationRepository()

	err := m.ValidateLocationRec(rec)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteLocationRec -
func (m *Model) DeleteLocationRec(recID string) error {

	m.Log.Debug("Deleting dungeon location rec ID >%s<", recID)

	r := m.LocationRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteLocationRec(recID)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveLocationRec -
func (m *Model) RemoveLocationRec(recID string) error {

	m.Log.Debug("Removing dungeon location rec ID >%s<", recID)

	r := m.LocationRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteLocationRec(recID)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
