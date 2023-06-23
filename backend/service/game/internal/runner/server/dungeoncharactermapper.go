package runner

import (
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	schema "gitlab.com/alienspaces/go-mud/backend/schema/game"
)

// dungeonCharacterResponseData
func dungeonCharacterResponseData(l logger.Logger, rs *InstanceViewRecordSet) (schema.DungeonCharacterData, error) {

	data := schema.DungeonCharacterData{
		ID:                  rs.CharacterInstanceViewRec.CharacterID,
		Name:                rs.CharacterInstanceViewRec.Name,
		Strength:            rs.CharacterInstanceViewRec.Strength,
		Dexterity:           rs.CharacterInstanceViewRec.Dexterity,
		Intelligence:        rs.CharacterInstanceViewRec.Intelligence,
		CurrentStrength:     rs.CharacterInstanceViewRec.CurrentStrength,
		CurrentDexterity:    rs.CharacterInstanceViewRec.CurrentDexterity,
		CurrentIntelligence: rs.CharacterInstanceViewRec.CurrentIntelligence,
		Health:              rs.CharacterInstanceViewRec.Health,
		Fatigue:             rs.CharacterInstanceViewRec.Fatigue,
		CurrentHealth:       rs.CharacterInstanceViewRec.CurrentHealth,
		CurrentFatigue:      rs.CharacterInstanceViewRec.CurrentFatigue,
		Coins:               rs.CharacterInstanceViewRec.Coins,
		ExperiencePoints:    rs.CharacterInstanceViewRec.ExperiencePoints,
		AttributePoints:     rs.CharacterInstanceViewRec.AttributePoints,
		CreatedAt:           rs.CharacterInstanceViewRec.CreatedAt,
		UpdatedAt:           rs.CharacterInstanceViewRec.UpdatedAt.Time,
	}

	data.Dungeon = &schema.DungeonCharacterDungeonData{
		ID:          rs.DungeonInstanceViewRec.DungeonID,
		Name:        rs.DungeonInstanceViewRec.Name,
		Description: rs.DungeonInstanceViewRec.Description,
	}

	data.Location = &schema.DungeonCharacterLocationData{
		ID:          rs.LocationInstanceViewRec.LocationID,
		Name:        rs.LocationInstanceViewRec.Name,
		Description: rs.LocationInstanceViewRec.Description,
	}

	return data, nil
}
