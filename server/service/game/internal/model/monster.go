package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetMonsterEquippedObjectRecs -
func (m *Model) GetMonsterEquippedObjectRecs(monsterID string) ([]*record.MonsterObject, error) {

	l := m.Logger("GetMonsterEquippedObjectRecs")

	l.Debug("Getting monster ID >%s< equipped object records", monsterID)

	r := m.MonsterObjectRepository()

	return r.GetMany(map[string]interface{}{
		"monster_id":  monsterID,
		"is_equipped": true,
	}, nil, false)
}

// GetMonsterStashedObjectRecs -
func (m *Model) GetMonsterStashedObjectRecs(monsterID string) ([]*record.MonsterObject, error) {

	l := m.Logger("GetMonsterStashedObjectRecs")

	l.Debug("Getting monster ID >%s< stashed object records", monsterID)

	r := m.MonsterObjectRepository()

	return r.GetMany(map[string]interface{}{
		"monster_id": monsterID,
		"is_stashed": true,
	}, nil, false)
}
