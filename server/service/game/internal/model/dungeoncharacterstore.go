package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetDungeonCharacterRecs -
func (m *Model) GetDungeonCharacterRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.DungeonCharacter, error) {

	m.Log.Info("Getting dungeon character records params >%s<", params)

	r := m.DungeonCharacterRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetDungeonCharacterRec -
func (m *Model) GetDungeonCharacterRec(recID string, forUpdate bool) (*record.DungeonCharacter, error) {

	m.Log.Info("Getting dungeon character rec ID >%s<", recID)

	r := m.DungeonCharacterRepository()

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

// CreateDungeonCharacterRec -
func (m *Model) CreateDungeonCharacterRec(rec *record.DungeonCharacter) error {

	m.Log.Info("Creating dungeon character rec >%#v<", rec)

	characterRepo := m.DungeonCharacterRepository()

	// Get default location record for the character dungeon
	locationRecs, err := m.GetDungeonLocationRecs(map[string]interface{}{
		"dungeon_id": rec.DungeonID,
		"default":    true,
	}, nil, false)
	if err != nil {
		msg := fmt.Sprintf("failed to get default dungeon location record for dungeon ID >%s< >%v<", rec.DungeonID, err)
		m.Log.Info(msg)
		return err
	}

	if len(locationRecs) != 1 {
		msg := fmt.Sprintf("unexpected number of dungeon location records returned for dungeon ID >%s<", rec.DungeonID)
		m.Log.Warn(msg)
		return fmt.Errorf(msg)
	}

	rec.DungeonLocationID = locationRecs[0].ID

	err = m.ValidateDungeonCharacterRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return characterRepo.CreateOne(rec)
}

// UpdateDungeonCharacterRec -
func (m *Model) UpdateDungeonCharacterRec(rec *record.DungeonCharacter) error {

	m.Log.Info("Updating dungeon character rec >%#v<", rec)

	r := m.DungeonCharacterRepository()

	err := m.ValidateDungeonCharacterRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteDungeonCharacterRec -
func (m *Model) DeleteDungeonCharacterRec(recID string) error {

	m.Log.Info("Deleting dungeon character rec ID >%s<", recID)

	r := m.DungeonCharacterRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonCharacterRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveDungeonCharacterRec -
func (m *Model) RemoveDungeonCharacterRec(recID string) error {

	m.Log.Info("Removing dungeon character rec ID >%s<", recID)

	r := m.DungeonCharacterRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonCharacterRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
