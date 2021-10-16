package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// TODO: Complete processing character action
func (m *Model) ProcessDungeonCharacterAction(dungeonCharacterID string, sentence string) (*record.DungeonAction, error) {

	m.Log.Info("Processing dungeon character ID >%s< action command >%s<", dungeonCharacterID, sentence)

	return nil, nil
}
