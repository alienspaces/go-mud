package runner

import (
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/schema"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// characterResponseData
func characterResponseData(l logger.Logger, characterRec *record.Character, instanceRecordSet *InstanceViewRecordSet) (schema.DungeonCharacterData, error) {

	data := schema.DungeonCharacterData{
		CharacterID:               characterRec.ID,
		CharacterName:             characterRec.Name,
		CharacterStrength:         characterRec.Strength,
		CharacterDexterity:        characterRec.Dexterity,
		CharacterIntelligence:     characterRec.Intelligence,
		CharacterHealth:           characterRec.Health,
		CharacterFatigue:          characterRec.Fatigue,
		CharacterCoins:            characterRec.Coins,
		CharacterExperiencePoints: characterRec.ExperiencePoints,
		CharacterAttributePoints:  characterRec.AttributePoints,
		CharacterCreatedAt:        characterRec.CreatedAt,
		CharacterUpdatedAt:        characterRec.UpdatedAt.Time,
	}

	if instanceRecordSet != nil {
		data.DungeonID = instanceRecordSet.DungeonInstanceViewRec.DungeonID
		data.DungeonName = instanceRecordSet.DungeonInstanceViewRec.Name
		data.DungeonDescription = instanceRecordSet.DungeonInstanceViewRec.Description
		data.LocationID = instanceRecordSet.LocationInstanceViewRec.LocationID
		data.LocationName = instanceRecordSet.LocationInstanceViewRec.Name
		data.LocationDescription = instanceRecordSet.LocationInstanceViewRec.Description
		data.CharacterCurrentStrength = instanceRecordSet.CharacterInstanceViewRec.CurrentStrength
		data.CharacterCurrentDexterity = instanceRecordSet.CharacterInstanceViewRec.CurrentDexterity
		data.CharacterCurrentIntelligence = instanceRecordSet.CharacterInstanceViewRec.CurrentIntelligence
		data.CharacterCurrentHealth = instanceRecordSet.CharacterInstanceViewRec.CurrentHealth
		data.CharacterCurrentFatigue = instanceRecordSet.CharacterInstanceViewRec.CurrentFatigue
	}

	return data, nil
}
