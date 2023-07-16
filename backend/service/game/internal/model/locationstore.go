package model

import (
	"database/sql"
	"fmt"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetLocationRecs -
func (m *Model) GetLocationRecs(opts *coresql.Options) ([]*record.Location, error) {

	l := m.Logger("GetLocationRecs")

	l.Debug("Getting location records opts >%#v<", opts)

	r := m.LocationRepository()

	return r.GetMany(opts)
}

// GetLocationRec -
func (m *Model) GetLocationRec(recID string, lock *coresql.Lock) (*record.Location, error) {

	l := m.Logger("GetLocationRec")

	l.Debug("Getting location rec ID >%s<", recID)

	r := m.LocationRepository()

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

// CreateLocationRec -
func (m *Model) CreateLocationRec(rec *record.Location) error {
	l := m.Logger("CreateLocationRec")

	l.Debug("Creating location record >%#v<", rec)

	r := m.LocationRepository()

	err := m.validateLocationRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateLocationRec -
func (m *Model) UpdateLocationRec(rec *record.Location) error {

	l := m.Logger("UpdateLocationRec")

	l.Debug("Updating location record >%#v<", rec)

	r := m.LocationRepository()

	err := m.validateLocationRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteLocationRec -
func (m *Model) DeleteLocationRec(recID string) error {

	l := m.Logger("DeleteLocationRec")

	l.Debug("Deleting location rec ID >%s<", recID)

	r := m.LocationRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteLocationRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveLocationRec -
func (m *Model) RemoveLocationRec(recID string) error {

	l := m.Logger("RemoveLocationRec")

	l.Debug("Removing location rec ID >%s<", recID)

	r := m.LocationRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteLocationRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
