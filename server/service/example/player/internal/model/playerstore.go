package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-boilerplate/server/service/player/internal/record"
)

// GetPlayerRecs -
func (m *Model) GetPlayerRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.Player, error) {

	m.Log.Info("Getting player records params >%s<", params)

	r := m.PlayerRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetPlayerRec -
func (m *Model) GetPlayerRec(recID string, forUpdate bool) (*record.Player, error) {

	m.Log.Info("Getting player rec ID >%s<", recID)

	r := m.PlayerRepository()

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

// CreatePlayerRec -
func (m *Model) CreatePlayerRec(rec *record.Player) error {

	m.Log.Info("Creating player rec >%v<", rec)

	r := m.PlayerRepository()

	err := m.ValidatePlayerRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdatePlayerRec -
func (m *Model) UpdatePlayerRec(rec *record.Player) error {

	m.Log.Info("Updating player rec >%v<", rec)

	r := m.PlayerRepository()

	err := m.ValidatePlayerRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeletePlayerRec -
func (m *Model) DeletePlayerRec(recID string) error {

	m.Log.Info("Deleting player rec ID >%s<", recID)

	r := m.PlayerRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeletePlayerRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemovePlayerRec -
func (m *Model) RemovePlayerRec(recID string) error {

	m.Log.Info("Removing player rec ID >%s<", recID)

	r := m.PlayerRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeletePlayerRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
