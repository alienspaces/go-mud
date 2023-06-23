package runner

import (
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	schema "gitlab.com/alienspaces/go-mud/backend/schema/game"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// characterResponseData
func characterResponseData(l logger.Logger, characterRec *record.Character, rs *InstanceViewRecordSet) (schema.DungeonCharacterData, error) {

	data := schema.DungeonCharacterData{
		ID:               characterRec.ID,
		Name:             characterRec.Name,
		Strength:         characterRec.Strength,
		Dexterity:        characterRec.Dexterity,
		Intelligence:     characterRec.Intelligence,
		Health:           characterRec.Health,
		Fatigue:          characterRec.Fatigue,
		Coins:            characterRec.Coins,
		ExperiencePoints: characterRec.ExperiencePoints,
		AttributePoints:  characterRec.AttributePoints,
		CreatedAt:        characterRec.CreatedAt,
		UpdatedAt:        characterRec.UpdatedAt.Time,
	}

	if rs != nil {
		dungeonData := schema.DungeonCharacterDungeonData{
			ID:          rs.DungeonInstanceViewRec.DungeonID,
			Name:        rs.DungeonInstanceViewRec.Name,
			Description: rs.DungeonInstanceViewRec.Description,
		}
		data.Dungeon = &dungeonData

		locationData := schema.DungeonCharacterLocationData{
			ID:          rs.LocationInstanceViewRec.LocationID,
			Name:        rs.LocationInstanceViewRec.Name,
			Description: rs.LocationInstanceViewRec.Description,
		}
		data.Location = &locationData

		data.CurrentStrength = rs.CharacterInstanceViewRec.CurrentStrength
		data.CurrentDexterity = rs.CharacterInstanceViewRec.CurrentDexterity
		data.CurrentIntelligence = rs.CharacterInstanceViewRec.CurrentIntelligence
		data.CurrentHealth = rs.CharacterInstanceViewRec.CurrentHealth
		data.CurrentFatigue = rs.CharacterInstanceViewRec.CurrentFatigue
	}

	return data, nil
}
