package model

import (
	"database/sql"
	"fmt"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetLocationObjectRecs -
func (m *Model) GetLocationObjectRecs(opts *coresql.Options) ([]*record.LocationObject, error) {

	l := m.loggerWithFunctionContext("GetLocationObjectRecs")

	l.Debug("Getting dungeon location records opts >%#v<", opts)

	r := m.LocationObjectRepository()

	return r.GetMany(opts)
}

// GetLocationObjectRec -
func (m *Model) GetLocationObjectRec(recID string, lock *coresql.Lock) (*record.LocationObject, error) {

	l := m.loggerWithFunctionContext("GetLocationObjectRec")

	l.Debug("Getting dungeon location rec ID >%s<", recID)

	r := m.LocationObjectRepository()

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

// CreateLocationObjectRec -
func (m *Model) CreateLocationObjectRec(rec *record.LocationObject) error {

	l := m.loggerWithFunctionContext("CreateLocationObjectRec")

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

	l := m.loggerWithFunctionContext("UpdateLocationObjectRec")

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

	l := m.loggerWithFunctionContext("DeleteLocationObjectRec")

	l.Debug("Deleting dungeon location rec ID >%s<", recID)

	r := m.LocationObjectRepository()

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

	l := m.loggerWithFunctionContext("RemoveLocationObjectRec")

	l.Debug("Removing dungeon location rec ID >%s<", recID)

	r := m.LocationObjectRepository()

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
