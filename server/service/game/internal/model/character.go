package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetCharacterObjectRecs -
func (m *Model) GetCharacterObjectRecs(characterID string) ([]*record.Object, error) {

	m.Log.Info("Getting character ID >%s< object records", characterID)

	r := m.ObjectRepository()

	return r.GetMany(map[string]interface{}{
		"character_id": characterID,
	}, nil, false)
}

// GetCharacterEquippedObjectRecs -
func (m *Model) GetCharacterEquippedObjectRecs(characterID string) ([]*record.Object, error) {

	m.Log.Info("Getting character ID >%s< equipped object records", characterID)

	r := m.ObjectRepository()

	return r.GetMany(map[string]interface{}{
		"character_id": characterID,
		"is_equipped":  true,
	}, nil, false)
}

// GetCharacterStashedObjectRecs -
func (m *Model) GetCharacterStashedObjectRecs(characterID string) ([]*record.Object, error) {

	m.Log.Info("Getting character ID >%s< stashed object records", characterID)

	r := m.ObjectRepository()

	return r.GetMany(map[string]interface{}{
		"character_id": characterID,
		"is_stashed":   true,
	}, nil, false)
}
