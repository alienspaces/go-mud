package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetLocationRecs -
func (m *Model) GetLocationRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.Location, error) {

	l := m.Logger("GetLocationRecs")

	l.Debug("Getting dungeon location records params >%s<", params)

	r := m.LocationRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetLocationRec -
func (m *Model) GetLocationRec(recID string, forUpdate bool) (*record.Location, error) {

	l := m.Logger("GetLocationRec")

	l.Debug("Getting dungeon location rec ID >%s<", recID)

	r := m.LocationRepository()

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

// CreateLocationRec -
func (m *Model) CreateLocationRec(rec *record.Location) error {

	l := m.Logger("CreateLocationRec")

	l.Debug("Creating dungeon location rec >%#v<", rec)

	r := m.LocationRepository()

	err := m.ValidateLocationRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateLocationRec -
func (m *Model) UpdateLocationRec(rec *record.Location) error {

	l := m.Logger("UpdateLocationRec")

	l.Debug("Updating dungeon location rec >%#v<", rec)

	r := m.LocationRepository()

	err := m.ValidateLocationRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteLocationRec -
func (m *Model) DeleteLocationRec(recID string) error {

	l := m.Logger("DeleteLocationRec")

	l.Debug("Deleting dungeon location rec ID >%s<", recID)

	r := m.LocationRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteLocationRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveLocationRec -
func (m *Model) RemoveLocationRec(recID string) error {

	l := m.Logger("RemoveLocationRec")

	l.Debug("Removing dungeon location rec ID >%s<", recID)

	r := m.LocationRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteLocationRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
