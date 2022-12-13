package runner

import (
	"gitlab.com/alienspaces/go-mud/backend/schema"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// DungeonRequestDataToRecord -
func (rnr *Runner) DungeonRequestDataToRecord(data schema.DungeonData, rec *record.Dungeon) error {

	return nil
}

// RecordToDungeonResponseData -
func (rnr *Runner) RecordToDungeonResponseData(dungeonRec record.Dungeon) (schema.DungeonData, error) {

	data := schema.DungeonData{
		DungeonID:          dungeonRec.ID,
		DungeonName:        dungeonRec.Name,
		DungeonDescription: dungeonRec.Description,
		DungeonCreatedAt:   dungeonRec.CreatedAt,
		DungeonUpdatedAt:   dungeonRec.UpdatedAt.Time,
	}

	return data, nil
}
