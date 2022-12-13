package runner

import (
	"gitlab.com/alienspaces/go-mud/backend/schema"
)

// InstanceViewRecordSetToDungeonCharacterResponse
func (rnr *Runner) InstanceViewRecordSetToDungeonCharacterResponseData(instanceRecordSet *InstanceViewRecordSet) (schema.DungeonCharacterData, error) {

	data := schema.DungeonCharacterData{
		DungeonID:                    instanceRecordSet.DungeonInstanceViewRec.DungeonID,
		DungeonName:                  instanceRecordSet.DungeonInstanceViewRec.Name,
		DungeonDescription:           instanceRecordSet.DungeonInstanceViewRec.Description,
		LocationID:                   instanceRecordSet.LocationInstanceViewRec.LocationID,
		LocationName:                 instanceRecordSet.LocationInstanceViewRec.Name,
		LocationDescription:          instanceRecordSet.LocationInstanceViewRec.Description,
		CharacterID:                  instanceRecordSet.CharacterInstanceViewRec.CharacterID,
		CharacterName:                instanceRecordSet.CharacterInstanceViewRec.Name,
		CharacterStrength:            instanceRecordSet.CharacterInstanceViewRec.Strength,
		CharacterDexterity:           instanceRecordSet.CharacterInstanceViewRec.Dexterity,
		CharacterIntelligence:        instanceRecordSet.CharacterInstanceViewRec.Intelligence,
		CharacterCurrentStrength:     instanceRecordSet.CharacterInstanceViewRec.CurrentStrength,
		CharacterCurrentDexterity:    instanceRecordSet.CharacterInstanceViewRec.CurrentDexterity,
		CharacterCurrentIntelligence: instanceRecordSet.CharacterInstanceViewRec.CurrentIntelligence,
		CharacterHealth:              instanceRecordSet.CharacterInstanceViewRec.Health,
		CharacterFatigue:             instanceRecordSet.CharacterInstanceViewRec.Fatigue,
		CharacterCurrentHealth:       instanceRecordSet.CharacterInstanceViewRec.CurrentHealth,
		CharacterCurrentFatigue:      instanceRecordSet.CharacterInstanceViewRec.CurrentFatigue,
		CharacterCoins:               instanceRecordSet.CharacterInstanceViewRec.Coins,
		CharacterExperiencePoints:    instanceRecordSet.CharacterInstanceViewRec.ExperiencePoints,
		CharacterAttributePoints:     instanceRecordSet.CharacterInstanceViewRec.AttributePoints,
		CharacterCreatedAt:           instanceRecordSet.CharacterInstanceViewRec.CreatedAt,
		CharacterUpdatedAt:           instanceRecordSet.CharacterInstanceViewRec.UpdatedAt.Time,
	}

	return data, nil
}
