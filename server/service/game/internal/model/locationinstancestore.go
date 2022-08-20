package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetLocationInstanceRecs -
func (m *Model) GetLocationInstanceRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.LocationInstance, error) {

	l := m.Logger("GetLocationInstanceRecs")

	l.Debug("Getting dungeon location records params >%s<", params)

	r := m.LocationInstanceRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetLocationInstanceRec -
func (m *Model) GetLocationInstanceRec(recID string, forUpdate bool) (*record.LocationInstance, error) {

	l := m.Logger("GetLocationInstanceRec")

	l.Debug("Getting dungeon location rec ID >%s<", recID)

	r := m.LocationInstanceRepository()

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

// GetLocationInstanceViewRecs -
func (m *Model) GetLocationInstanceViewRecs(params map[string]interface{}, operators map[string]string) ([]*record.LocationInstanceView, error) {

	l := m.Logger("GetLocationInstanceViewRecs")

	l.Debug("Getting dungeon location records params >%s<", params)

	r := m.LocationInstanceViewRepository()

	return r.GetMany(params, operators, false)
}

// GetLocationInstanceViewRec -
func (m *Model) GetLocationInstanceViewRec(recID string) (*record.LocationInstanceView, error) {

	l := m.Logger("GetLocationInstanceViewRec")

	l.Debug("Getting dungeon location rec ID >%s<", recID)

	r := m.LocationInstanceViewRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return nil, fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	rec, err := r.GetOne(recID, false)
	if err == sql.ErrNoRows {
		l.Warn("No record found ID >%s<", recID)
		return nil, nil
	}

	return rec, err
}

// CreateLocationInstanceRec -
func (m *Model) CreateLocationInstanceRec(rec *record.LocationInstance) error {

	l := m.Logger("CreateLocationInstanceRec")

	l.Debug("Creating dungeon location record >%#v<", rec)

	r := m.LocationInstanceRepository()

	err := m.validateLocationInstanceRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateLocationInstanceRec -
func (m *Model) UpdateLocationInstanceRec(rec *record.LocationInstance) error {

	l := m.Logger("UpdateLocationInstanceRec")

	l.Debug("Updating dungeon location record >%#v<", rec)

	r := m.LocationInstanceRepository()

	err := m.validateLocationInstanceRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteLocationInstanceRec -
func (m *Model) DeleteLocationInstanceRec(recID string) error {

	l := m.Logger("DeleteLocationInstanceRec")

	l.Debug("Deleting dungeon location rec ID >%s<", recID)

	r := m.LocationInstanceRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteLocationInstanceRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveLocationInstanceRec -
func (m *Model) RemoveLocationInstanceRec(recID string) error {

	l := m.Logger("RemoveLocationInstanceRec")

	l.Debug("Removing dungeon location rec ID >%s<", recID)

	r := m.LocationInstanceRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteLocationInstanceRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
