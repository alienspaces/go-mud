package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetGameRecs -
func (m *Model) GetGameRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.Game, error) {

	m.Log.Info("Getting game records params >%s<", params)

	r := m.GameRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetGameRec -
func (m *Model) GetGameRec(recID string, forUpdate bool) (*record.Game, error) {

	m.Log.Info("Getting game rec ID >%s<", recID)

	r := m.GameRepository()

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

// CreateGameRec -
func (m *Model) CreateGameRec(rec *record.Game) error {

	m.Log.Info("Creating game rec >%#v<", rec)

	r := m.GameRepository()

	err := m.ValidateGameRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateGameRec -
func (m *Model) UpdateGameRec(rec *record.Game) error {

	m.Log.Info("Updating game rec >%#v<", rec)

	r := m.GameRepository()

	err := m.ValidateGameRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteGameRec -
func (m *Model) DeleteGameRec(recID string) error {

	m.Log.Info("Deleting game rec ID >%s<", recID)

	r := m.GameRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteGameRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveGameRec -
func (m *Model) RemoveGameRec(recID string) error {

	m.Log.Info("Removing game rec ID >%s<", recID)

	r := m.GameRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteGameRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
