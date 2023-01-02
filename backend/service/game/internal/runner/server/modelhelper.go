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
	l = loggerWithContext(l, "getInstanceViewRecordSetByCharacterID")

	characterInstanceViewRec, err := m.(*model.Model).GetCharacterInstanceViewRecByCharacterID(characterID)
	if err != nil {
		return nil, err
	}

	if characterInstanceViewRec == nil {
		l.Warn("character instance record is nil")
		err := coreerror.NewInternalError()
		return nil, err
	}

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
