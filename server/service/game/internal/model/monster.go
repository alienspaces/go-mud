package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetMonsterObjectRecs -
func (m *Model) GetMonsterObjectRecs(monsterID string) ([]*record.Object, error) {

	m.Log.Info("Getting monster ID >%s< object records", monsterID)

	r := m.ObjectRepository()

	return r.GetMany(map[string]interface{}{
		"dungeon_monster_id": monsterID,
	}, nil, false)
}

// GetMonsterEquippedObjectRecs -
func (m *Model) GetMonsterEquippedObjectRecs(monsterID string) ([]*record.Object, error) {

	m.Log.Info("Getting monster ID >%s< equipped object records", monsterID)

	r := m.ObjectRepository()

	return r.GetMany(map[string]interface{}{
		"dungeon_monster_id": monsterID,
		"is_equipped":        true,
	}, nil, false)
}

// GetMonsterStashedObjectRecs -
func (m *Model) GetMonsterStashedObjectRecs(monsterID string) ([]*record.Object, error) {

	m.Log.Info("Getting monster ID >%s< stashed object records", monsterID)

	r := m.ObjectRepository()

	return r.GetMany(map[string]interface{}{
		"dungeon_monster_id": monsterID,
		"is_stashed":         true,
	}, nil, false)
}
