package model

import (
	"database/sql"
	"fmt"
	"time"

	"gitlab.com/alienspaces/go-mud/backend/core/null"
	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetTurnRecs -
func (m *Model) GetTurnRecs(opts *coresql.Options) ([]*record.Turn, error) {

	l := m.loggerWithContext("GetTurnRecs")

	l.Debug("Getting turn records opts >%#v<", opts)

	r := m.TurnRepository()

	return r.GetMany(opts)
}

// GetTurnRec -
func (m *Model) GetTurnRec(recID string, lock *coresql.Lock) (*record.Turn, error) {

	l := m.loggerWithContext("GetTurnRec")

	l.Debug("Getting turn rec ID >%s<", recID)

	r := m.TurnRepository()

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

// CreateTurnRec -
func (m *Model) CreateTurnRec(rec *record.Turn) error {
	l := m.loggerWithContext("CreateTurnRec")

	// Initial defaults
	rec.TurnNumber = 1
	rec.IncrementedAt = null.NullTimeFromTime(time.Now().UTC())

	l.Debug("Creating turn record >%#v<", rec)

	r := m.TurnRepository()

	err := m.validateTurnRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateTurnRec -
func (m *Model) UpdateTurnRec(rec *record.Turn) error {
	l := m.loggerWithContext("UpdateTurnRec")

	l.Debug("Updating turn record >%#v<", rec)

	r := m.TurnRepository()

	err := m.validateTurnRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteTurnRec -
func (m *Model) DeleteTurnRec(recID string) error {
	l := m.loggerWithContext("DeleteTurnRec")

	l.Debug("Deleting turn rec ID >%s<", recID)

	r := m.TurnRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteTurnRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveTurnRec -
func (m *Model) RemoveTurnRec(recID string) error {

	l := m.loggerWithContext("RemoveTurnRec")

	l.Debug("Removing turn rec ID >%s<", recID)

	r := m.TurnRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteTurnRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
