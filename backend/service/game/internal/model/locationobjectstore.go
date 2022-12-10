package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetLocationObjectRecs -
func (m *Model) GetLocationObjectRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.LocationObject, error) {

	l := m.Logger("GetLocationObjectRecs")

	l.Debug("Getting dungeon location records params >%s<", params)

	r := m.LocationObjectRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetLocationObjectRec -
func (m *Model) GetLocationObjectRec(recID string, forUpdate bool) (*record.LocationObject, error) {

	l := m.Logger("GetLocationObjectRec")

	l.Debug("Getting dungeon location rec ID >%s<", recID)

	r := m.LocationObjectRepository()

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

// CreateLocationObjectRec -
func (m *Model) CreateLocationObjectRec(rec *record.LocationObject) error {

	l := m.Logger("CreateLocationObjectRec")

	l.Debug("Creating dungeon location record >%#v<", rec)

	r := m.LocationObjectRepository()

	err := m.validateLocationObjectRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateLocationObjectRec -
func (m *Model) UpdateLocationObjectRec(rec *record.LocationObject) error {

	l := m.Logger("UpdateLocationObjectRec")

	l.Debug("Updating dungeon location record >%#v<", rec)

	r := m.LocationObjectRepository()

	err := m.validateLocationObjectRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteLocationObjectRec -
func (m *Model) DeleteLocationObjectRec(recID string) error {

	l := m.Logger("DeleteLocationObjectRec")

	l.Debug("Deleting dungeon location rec ID >%s<", recID)

	r := m.LocationObjectRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteLocationObjectRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveLocationObjectRec -
func (m *Model) RemoveLocationObjectRec(recID string) error {

	l := m.Logger("RemoveLocationObjectRec")

	l.Debug("Removing dungeon location rec ID >%s<", recID)

	r := m.LocationObjectRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteLocationObjectRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
