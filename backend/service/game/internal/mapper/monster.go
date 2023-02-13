package mapper

import (
	"gitlab.com/alienspaces/go-mud/backend/core/repository"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// MonsterInstanceViewToMonsterInstance -
func MonsterInstanceViewToMonsterInstance(l logger.Logger, characterInstanceViewRec *record.MonsterInstanceView) (*record.MonsterInstance, error) {

	characterInstanceRec := record.MonsterInstance{
		Record: repository.Record{
			ID:        characterInstanceViewRec.ID,
			CreatedAt: characterInstanceViewRec.CreatedAt,
			UpdatedAt: characterInstanceViewRec.UpdatedAt,
		},
		MonsterID:          characterInstanceViewRec.MonsterID,
		DungeonInstanceID:  characterInstanceViewRec.DungeonInstanceID,
		LocationInstanceID: characterInstanceViewRec.LocationInstanceID,
		Strength:           characterInstanceViewRec.CurrentStrength,
		Dexterity:          characterInstanceViewRec.CurrentDexterity,
		Intelligence:       characterInstanceViewRec.CurrentIntelligence,
		Health:             characterInstanceViewRec.CurrentHealth,
		Fatigue:            characterInstanceViewRec.CurrentFatigue,
		Coins:              characterInstanceViewRec.Coins,
		ExperiencePoints:   characterInstanceViewRec.ExperiencePoints,
		AttributePoints:    characterInstanceViewRec.AttributePoints,
	}

	return &characterInstanceRec, nil
}
