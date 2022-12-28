package runner

import (
	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

type InstanceViewRecordSet struct {
	CharacterInstanceViewRec *record.CharacterInstanceView
	DungeonInstanceViewRec   *record.DungeonInstanceView
	LocationInstanceViewRec  *record.LocationInstanceView
}

func (rnr *Runner) getInstanceViewRecordSetByCharacterID(l logger.Logger, m modeller.Modeller, characterID string) (*InstanceViewRecordSet, error) {
	l = HTTPLogger(l, "getInstanceViewRecordSetByCharacterID")

	characterInstanceViewRecs, err := m.(*model.Model).GetCharacterInstanceViewRecs(
		map[string]interface{}{
			"character_id": characterID,
		}, nil,
	)
	if err != nil {
		return nil, err
	}

	if len(characterInstanceViewRecs) == 0 {
		l.Warn("Character with ID %s has not entered a dungeon", characterID)
		return nil, nil
	}

	if len(characterInstanceViewRecs) > 1 {
		l.Warn("Unexpected number of character instance records returned >%d<", len(characterInstanceViewRecs))
		err := coreerror.NewInternalError()
		return nil, err
	}

	characterInstanceViewRec := characterInstanceViewRecs[0]

	dungeonInstanceViewRec, err := m.(*model.Model).GetDungeonInstanceViewRec(characterInstanceViewRec.DungeonInstanceID)
	if err != nil {
		l.Warn("Failed to get dungeon instance view record ID >%s<", characterInstanceViewRec.DungeonInstanceID)
		return nil, err
	}

	if dungeonInstanceViewRec == nil {
		l.Warn("Dungeon instance record ID >%s< does not exist", characterInstanceViewRec.DungeonInstanceID)
		err := coreerror.NewInternalError()
		return nil, err
	}

	locationInstanceViewRec, err := m.(*model.Model).GetLocationInstanceViewRec(characterInstanceViewRec.LocationInstanceID)
	if err != nil {
		l.Warn("Failed to get location instance record ID >%s<", characterInstanceViewRec.LocationInstanceID)
		return nil, err
	}

	if locationInstanceViewRec == nil {
		l.Warn("Location instance record ID >%s< does not exist", characterInstanceViewRec.LocationInstanceID)
		err := coreerror.NewInternalError()
		return nil, err
	}

	rs := InstanceViewRecordSet{
		CharacterInstanceViewRec: characterInstanceViewRec,
		DungeonInstanceViewRec:   dungeonInstanceViewRec,
		LocationInstanceViewRec:  locationInstanceViewRec,
	}

	return &rs, nil
}
