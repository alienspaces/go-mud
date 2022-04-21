package mapper

import (
	"gitlab.com/alienspaces/go-mud/server/core/repository"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// CharacterInstanceViewToCharacterInstance -
func CharacterInstanceViewToCharacterInstance(l logger.Logger, characterInstanceViewRec *record.CharacterInstanceView) (*record.CharacterInstance, error) {

	characterInstanceRec := record.CharacterInstance{
		Record: repository.Record{
			ID:        characterInstanceViewRec.ID,
			CreatedAt: characterInstanceViewRec.CreatedAt,
			UpdatedAt: characterInstanceViewRec.UpdatedAt,
		},
		CharacterID:        characterInstanceViewRec.CharacterID,
		DungeonInstanceID:  characterInstanceViewRec.DungeonInstanceID,
		LocationInstanceID: characterInstanceViewRec.LocationInstanceID,
		Strength:           characterInstanceViewRec.Strength,
		Dexterity:          characterInstanceViewRec.Dexterity,
		Intelligence:       characterInstanceViewRec.Intelligence,
		Health:             characterInstanceViewRec.Health,
		Fatigue:            characterInstanceViewRec.Fatigue,
		Coins:              characterInstanceViewRec.Coins,
		ExperiencePoints:   characterInstanceViewRec.ExperiencePoints,
		AttributePoints:    characterInstanceViewRec.AttributePoints,
	}

	return &characterInstanceRec, nil
}
