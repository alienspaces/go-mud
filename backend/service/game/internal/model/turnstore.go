package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetTurnRecs -
func (m *Model) GetTurnRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.Turn, error) {

	l := m.Logger("GetTurnRecs")

	l.Debug("Getting turn records params >%s<", params)

	r := m.TurnRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetTurnRec -
func (m *Model) GetTurnRec(recID string, forUpdate bool) (*record.Turn, error) {

	l := m.Logger("GetTurnRec")

	l.Debug("Getting turn rec ID >%s<", recID)

	r := m.TurnRepository()

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

// CreateTurnRec -
func (m *Model) CreateTurnRec(rec *record.Turn) error {
	l := m.Logger("CreateTurnRec")

	// All turns start at turn
	rec.Turn = 1

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

	l := m.Logger("UpdateTurnRec")

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

	l := m.Logger("DeleteTurnRec")

	l.Debug("Deleting turn rec ID >%s<", recID)

	r := m.TurnRepository()

	// validate UUID
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

	l := m.Logger("RemoveTurnRec")

	l.Debug("Removing turn rec ID >%s<", recID)

	r := m.TurnRepository()

	// validate UUID
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
