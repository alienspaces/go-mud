package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetCharacterInstanceObjectInstanceRecs -
func (m *Model) GetCharacterInstanceObjectInstanceRecs(characterID string) ([]*record.ObjectInstance, error) {

	m.Log.Info("Getting character ID >%s< object records", characterID)

	r := m.ObjectInstanceRepository()

	return r.GetMany(map[string]interface{}{
		"character_instance_id": characterID,
	}, nil, false)
}

// GetCharacterInstanceEquippedObjectInstanceRecs -
func (m *Model) GetCharacterInstanceEquippedObjectInstanceRecs(characterID string) ([]*record.ObjectInstance, error) {

	m.Log.Info("Getting character ID >%s< equipped object records", characterID)

	r := m.ObjectInstanceRepository()

	return r.GetMany(map[string]interface{}{
		"character_instance_id": characterID,
		"is_equipped":           true,
	}, nil, false)
}

// GetCharacterInstanceStashedObjectInstanceRecs -
func (m *Model) GetCharacterInstanceStashedObjectInstanceRecs(characterID string) ([]*record.ObjectInstance, error) {

	m.Log.Info("Getting character ID >%s< stashed object records", characterID)

	r := m.ObjectInstanceRepository()

	return r.GetMany(map[string]interface{}{
		"character_instance_id": characterID,
		"is_stashed":            true,
	}, nil, false)
}

// GetCharacterInstanceObjectInstanceViewRecs -
func (m *Model) GetCharacterInstanceObjectInstanceViewRecs(characterID string) ([]*record.ObjectInstanceView, error) {

	m.Log.Info("Getting character ID >%s< object records", characterID)

	r := m.ObjectInstanceViewRepository()

	return r.GetMany(map[string]interface{}{
		"character_instance_id": characterID,
	}, nil, false)
}

// GetCharacterInstanceEquippedObjectInstanceViewRecs -
func (m *Model) GetCharacterInstanceEquippedObjectInstanceViewRecs(characterID string) ([]*record.ObjectInstanceView, error) {

	m.Log.Info("Getting character ID >%s< equipped object records", characterID)

	r := m.ObjectInstanceViewRepository()

	return r.GetMany(map[string]interface{}{
		"character_instance_id": characterID,
		"is_equipped":           true,
	}, nil, false)
}

// GetCharacterInstanceStashedObjectInstanceViewRecs -
func (m *Model) GetCharacterInstanceStashedObjectInstanceViewRecs(characterID string) ([]*record.ObjectInstanceView, error) {

	m.Log.Info("Getting character ID >%s< stashed object records", characterID)

	r := m.ObjectInstanceViewRepository()

	return r.GetMany(map[string]interface{}{
		"character_instance_id": characterID,
		"is_stashed":            true,
	}, nil, false)
}
