package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetMonsterInstanceObjectInstanceRecs -
func (m *Model) GetMonsterInstanceObjectInstanceRecs(characterID string) ([]*record.ObjectInstance, error) {

	m.Log.Info("Getting character ID >%s< object records", characterID)

	r := m.ObjectInstanceRepository()

	return r.GetMany(map[string]interface{}{
		"character_instance_id": characterID,
	}, nil, false)
}

// GetMonsterInstanceEquippedObjectInstanceRecs -
func (m *Model) GetMonsterInstanceEquippedObjectInstanceRecs(characterID string) ([]*record.ObjectInstance, error) {

	m.Log.Info("Getting character ID >%s< equipped object records", characterID)

	r := m.ObjectInstanceRepository()

	return r.GetMany(map[string]interface{}{
		"character_instance_id": characterID,
		"is_equipped":           true,
	}, nil, false)
}

// GetMonsterInstanceStashedObjectInstanceRecs -
func (m *Model) GetMonsterInstanceStashedObjectInstanceRecs(characterID string) ([]*record.ObjectInstance, error) {

	m.Log.Info("Getting character ID >%s< stashed object records", characterID)

	r := m.ObjectInstanceRepository()

	return r.GetMany(map[string]interface{}{
		"character_instance_id": characterID,
		"is_stashed":            true,
	}, nil, false)
}

// GetMonsterInstanceObjectInstanceViewRecs -
func (m *Model) GetMonsterInstanceObjectInstanceViewRecs(characterID string) ([]*record.ObjectInstanceView, error) {

	m.Log.Info("Getting character ID >%s< object records", characterID)

	r := m.ObjectInstanceViewRepository()

	return r.GetMany(map[string]interface{}{
		"character_instance_id": characterID,
	}, nil, false)
}

// GetMonsterInstanceEquippedObjectInstanceViewRecs -
func (m *Model) GetMonsterInstanceEquippedObjectInstanceViewRecs(characterID string) ([]*record.ObjectInstanceView, error) {

	m.Log.Info("Getting character ID >%s< equipped object records", characterID)

	r := m.ObjectInstanceViewRepository()

	return r.GetMany(map[string]interface{}{
		"character_instance_id": characterID,
		"is_equipped":           true,
	}, nil, false)
}

// GetMonsterInstanceStashedObjectInstanceViewRecs -
func (m *Model) GetMonsterInstanceStashedObjectInstanceViewRecs(characterID string) ([]*record.ObjectInstanceView, error) {

	m.Log.Info("Getting character ID >%s< stashed object records", characterID)

	r := m.ObjectInstanceViewRepository()

	return r.GetMany(map[string]interface{}{
		"character_instance_id": characterID,
		"is_stashed":            true,
	}, nil, false)
}
