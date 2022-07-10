package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetCharacterEquippedObjectRecs -
func (m *Model) GetCharacterEquippedObjectRecs(characterID string) ([]*record.CharacterObject, error) {
	l := m.Logger("GetCharacterEquippedObjectRecs")

	l.Debug("Getting character ID >%s< equipped object records", characterID)

	r := m.CharacterObjectRepository()

	return r.GetMany(map[string]interface{}{
		"character_id": characterID,
		"is_equipped":  true,
	}, nil, false)
}

// GetCharacterStashedObjectRecs -
func (m *Model) GetCharacterStashedObjectRecs(characterID string) ([]*record.CharacterObject, error) {

	l := m.Logger("GetCharacterStashedObjectRecs")

	l.Debug("Getting character ID >%s< stashed object records", characterID)

	r := m.CharacterObjectRepository()

	return r.GetMany(map[string]interface{}{
		"character_id": characterID,
		"is_stashed":   true,
	}, nil, false)
}
