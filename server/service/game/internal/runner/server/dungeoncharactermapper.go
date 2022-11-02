package runner

import (
	"gitlab.com/alienspaces/go-mud/server/schema"
)

// InstanceViewRecordSetToDungeonCharacterResponse
func (rnr *Runner) InstanceViewRecordSetToDungeonCharacterResponseData(instanceRecordSet *InstanceViewRecordSet) (schema.DungeonCharacterData, error) {

	data := schema.DungeonCharacterData{
		ID:                  instanceRecordSet.CharacterInstanceViewRec.CharacterID,
		DungeonID:           instanceRecordSet.DungeonInstanceViewRec.DungeonID,
		DungeonName:         instanceRecordSet.DungeonInstanceViewRec.Name,
		DungeonDescription:  instanceRecordSet.DungeonInstanceViewRec.Description,
		LocationID:          instanceRecordSet.LocationInstanceViewRec.LocationID,
		LocationName:        instanceRecordSet.LocationInstanceViewRec.Name,
		LocationDescription: instanceRecordSet.LocationInstanceViewRec.Description,
		Name:                instanceRecordSet.CharacterInstanceViewRec.Name,
		Strength:            instanceRecordSet.CharacterInstanceViewRec.Strength,
		Dexterity:           instanceRecordSet.CharacterInstanceViewRec.Dexterity,
		Intelligence:        instanceRecordSet.CharacterInstanceViewRec.Intelligence,
		CurrentStrength:     instanceRecordSet.CharacterInstanceViewRec.CurrentStrength,
		CurrentDexterity:    instanceRecordSet.CharacterInstanceViewRec.CurrentDexterity,
		CurrentIntelligence: instanceRecordSet.CharacterInstanceViewRec.CurrentIntelligence,
		Health:              instanceRecordSet.CharacterInstanceViewRec.Health,
		Fatigue:             instanceRecordSet.CharacterInstanceViewRec.Fatigue,
		CurrentHealth:       instanceRecordSet.CharacterInstanceViewRec.CurrentHealth,
		CurrentFatigue:      instanceRecordSet.CharacterInstanceViewRec.CurrentFatigue,
		Coins:               instanceRecordSet.CharacterInstanceViewRec.Coins,
		ExperiencePoints:    instanceRecordSet.CharacterInstanceViewRec.ExperiencePoints,
		AttributePoints:     instanceRecordSet.CharacterInstanceViewRec.AttributePoints,
		CreatedAt:           instanceRecordSet.CharacterInstanceViewRec.CreatedAt,
		UpdatedAt:           instanceRecordSet.CharacterInstanceViewRec.UpdatedAt.Time,
	}

	return data, nil
}
