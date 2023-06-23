package runner

import (
	schema "gitlab.com/alienspaces/go-mud/backend/schema/game"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// DungeonRequestDataToRecord -
func (rnr *Runner) DungeonRequestDataToRecord(data schema.DungeonData, rec *record.Dungeon) error {

	return nil
}

// RecordToDungeonResponseData -
func (rnr *Runner) RecordToDungeonResponseData(dungeonRec record.Dungeon) (schema.DungeonData, error) {

	data := schema.DungeonData{
		ID:          dungeonRec.ID,
		Name:        dungeonRec.Name,
		Description: dungeonRec.Description,
		CreatedAt:   dungeonRec.CreatedAt,
		UpdatedAt:   dungeonRec.UpdatedAt.Time,
	}

	return data, nil
}
