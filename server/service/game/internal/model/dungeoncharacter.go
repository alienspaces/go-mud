package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetCharacterObjectRecs -
func (m *Model) GetCharacterObjectRecs(characterID string) ([]*record.DungeonObject, error) {

	m.Log.Info("Getting character ID >%s< object records", characterID)

	r := m.DungeonObjectRepository()

	return r.GetMany(map[string]interface{}{
		"dungeon_character_id": characterID,
	}, nil, false)
}

// GetCharacterEquippedDungeonObjectRecs -
func (m *Model) GetCharacterEquippedDungeonObjectRecs(characterID string) ([]*record.DungeonObject, error) {

	m.Log.Info("Getting character ID >%s< equipped object records", characterID)

	r := m.DungeonObjectRepository()

	return r.GetMany(map[string]interface{}{
		"dungeon_character_id": characterID,
		"is_equipped":          true,
	}, nil, false)
}

// GetCharacterStashedDungeonObjectRecs -
func (m *Model) GetCharacterStashedDungeonObjectRecs(characterID string) ([]*record.DungeonObject, error) {

	m.Log.Info("Getting character ID >%s< stashed object records", characterID)

	r := m.DungeonObjectRepository()

	return r.GetMany(map[string]interface{}{
		"dungeon_character_id": characterID,
		"is_stashed":           true,
	}, nil, false)
}
