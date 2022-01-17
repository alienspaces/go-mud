package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetCharacterEquippedDungeonObjectRecs -
func (m *Model) GetCharacterEquippedDungeonObjectRecs(characterID string) ([]*record.DungeonObject, error) {

	m.Log.Info("Getting character ID >%s< equipped dungeon object records", characterID)

	r := m.DungeonObjectRepository()

	return r.GetMany(map[string]interface{}{
		"dungeon_character_id": characterID,
		"is_equipped":          true,
	}, nil, false)
}

// GetCharacterStashedDungeonObjectRecs -
func (m *Model) GetCharacterStashedDungeonObjectRecs(characterID string) ([]*record.DungeonObject, error) {

	m.Log.Info("Getting character ID >%s< stashed dungeon object records", characterID)

	r := m.DungeonObjectRepository()

	return r.GetMany(map[string]interface{}{
		"dungeon_character_id": characterID,
		"is_stashed":           true,
	}, nil, false)
}
