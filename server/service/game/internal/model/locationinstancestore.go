package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetLocationInstanceRecs -
func (m *Model) GetLocationInstanceRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.LocationInstance, error) {

	m.Log.Debug("Getting dungeon location records params >%s<", params)

	r := m.LocationInstanceRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetLocationInstanceRec -
func (m *Model) GetLocationInstanceRec(recID string, forUpdate bool) (*record.LocationInstance, error) {

	m.Log.Debug("Getting dungeon location rec ID >%s<", recID)

	r := m.LocationInstanceRepository()

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

// GetLocationInstanceViewRecs -
func (m *Model) GetLocationInstanceViewRecs(params map[string]interface{}, operators map[string]string) ([]*record.LocationInstanceView, error) {

	m.Log.Debug("Getting dungeon location records params >%s<", params)

	r := m.LocationInstanceViewRepository()

	return r.GetMany(params, operators, false)
}

// GetLocationInstanceViewRec -
func (m *Model) GetLocationInstanceViewRec(recID string) (*record.LocationInstanceView, error) {

	m.Log.Debug("Getting dungeon location rec ID >%s<", recID)

	r := m.LocationInstanceViewRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return nil, fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	rec, err := r.GetOne(recID, false)
	if err == sql.ErrNoRows {
		m.Log.Warn("No record found ID >%s<", recID)
		return nil, nil
	}

	return rec, err
}

// CreateLocationInstanceRec -
func (m *Model) CreateLocationInstanceRec(rec *record.LocationInstance) error {

	m.Log.Debug("Creating dungeon location rec >%#v<", rec)

	r := m.LocationInstanceRepository()

	err := m.ValidateLocationInstanceRec(rec)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateLocationInstanceRec -
func (m *Model) UpdateLocationInstanceRec(rec *record.LocationInstance) error {

	m.Log.Debug("Updating dungeon location rec >%#v<", rec)

	r := m.LocationInstanceRepository()

	err := m.ValidateLocationInstanceRec(rec)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteLocationInstanceRec -
func (m *Model) DeleteLocationInstanceRec(recID string) error {

	m.Log.Debug("Deleting dungeon location rec ID >%s<", recID)

	r := m.LocationInstanceRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteLocationInstanceRec(recID)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveLocationInstanceRec -
func (m *Model) RemoveLocationInstanceRec(recID string) error {

	m.Log.Debug("Removing dungeon location rec ID >%s<", recID)

	r := m.LocationInstanceRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteLocationInstanceRec(recID)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
