package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// TODO: Complete processing character action
func (m *Model) ProcessDungeonCharacterAction(dungeonCharacterID string, sentence string) (*record.DungeonAction, error) {

	m.Log.Info("Processing dungeon character ID >%s< action command >%s<", dungeonCharacterID, sentence)

	// Get current dungeon location record set

	// Resolve character action

	// Perform character action

	// Refetch current dungeon location record set

	// Create dungeon action event records

	return nil, nil
}
