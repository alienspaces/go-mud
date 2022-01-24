package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetMonsterObjectRecs -
func (m *Model) GetMonsterObjectRecs(monsterID string) ([]*record.DungeonObject, error) {

	m.Log.Info("Getting monster ID >%s< object records", monsterID)

	r := m.DungeonObjectRepository()

	return r.GetMany(map[string]interface{}{
		"dungeon_monster_id": monsterID,
	}, nil, false)
}

// GetMonsterEquippedDungeonObjectRecs -
func (m *Model) GetMonsterEquippedDungeonObjectRecs(monsterID string) ([]*record.DungeonObject, error) {

	m.Log.Info("Getting monster ID >%s< equipped object records", monsterID)

	r := m.DungeonObjectRepository()

	return r.GetMany(map[string]interface{}{
		"dungeon_monster_id": monsterID,
		"is_equipped":        true,
	}, nil, false)
}

// GetMonsterStashedDungeonObjectRecs -
func (m *Model) GetMonsterStashedDungeonObjectRecs(monsterID string) ([]*record.DungeonObject, error) {

	m.Log.Info("Getting monster ID >%s< stashed object records", monsterID)

	r := m.DungeonObjectRepository()

	return r.GetMany(map[string]interface{}{
		"dungeon_monster_id": monsterID,
		"is_stashed":         true,
	}, nil, false)
}
