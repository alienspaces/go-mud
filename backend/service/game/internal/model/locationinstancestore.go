package model

import (
	"database/sql"
	"fmt"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetLocationInstanceRecs -
func (m *Model) GetLocationInstanceRecs(opts *coresql.Options) ([]*record.LocationInstance, error) {

	l := m.loggerWithContext("GetLocationInstanceRecs")

	l.Debug("Getting dungeon location records opts >%#v<", opts)

	r := m.LocationInstanceRepository()

	return r.GetMany(opts)
}

// GetLocationInstanceRec -
func (m *Model) GetLocationInstanceRec(recID string, lock *coresql.Lock) (*record.LocationInstance, error) {

	l := m.loggerWithContext("GetLocationInstanceRec")

	l.Debug("Getting dungeon location rec ID >%s<", recID)

	r := m.LocationInstanceRepository()

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

// GetLocationInstanceViewRecs -
func (m *Model) GetLocationInstanceViewRecs(opts *coresql.Options) ([]*record.LocationInstanceView, error) {

	l := m.loggerWithContext("GetLocationInstanceViewRecs")

	l.Debug("Getting dungeon location records opts >%#v<", opts)

	r := m.LocationInstanceViewRepository()

	return r.GetMany(opts)
}

// GetLocationInstanceViewRec -
func (m *Model) GetLocationInstanceViewRec(recID string) (*record.LocationInstanceView, error) {

	l := m.loggerWithContext("GetLocationInstanceViewRec")

	l.Debug("Getting dungeon location rec ID >%s<", recID)

	r := m.LocationInstanceViewRepository()

	if !m.IsUUID(recID) {
		return nil, fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	rec, err := r.GetOne(recID, nil)
	if err == sql.ErrNoRows {
		l.Warn("No record found ID >%s<", recID)
		return nil, nil
	}

	return rec, err
}

// CreateLocationInstanceRec -
func (m *Model) CreateLocationInstanceRec(rec *record.LocationInstance) error {

	l := m.loggerWithContext("CreateLocationInstanceRec")

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

	l := m.loggerWithContext("UpdateLocationInstanceRec")

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

	l := m.loggerWithContext("DeleteLocationInstanceRec")

	l.Debug("Deleting dungeon location rec ID >%s<", recID)

	r := m.LocationInstanceRepository()

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

	l := m.loggerWithContext("RemoveLocationInstanceRec")

	l.Debug("Removing dungeon location rec ID >%s<", recID)

	r := m.LocationInstanceRepository()

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
