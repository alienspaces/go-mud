package runner

import (
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/data/cabin"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/data/cave"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/harness"
)

// The following data is actual game data used to seed a new game server
func SeedDataConfig() harness.DataConfig {
	var objectConfig = cave.ObjectConfig()
	objectConfig = append(objectConfig, cabin.ObjectConfig()...)

	var monsterConfig = cave.MonsterConfig()
	monsterConfig = append(monsterConfig, cabin.MonsterConfig()...)

	caveDungeonConfig := cave.DungeonConfig()
	caveDungeonConfig.DungeonInstanceConfig = nil

	cabinDungeonConfig := cabin.DungeonConfig()
	cabinDungeonConfig.DungeonInstanceConfig = nil

	var dungeonConfig = []harness.DungeonConfig{
		caveDungeonConfig,
		cabinDungeonConfig,
	}

	d := harness.DataConfig{
		ObjectConfig:  objectConfig,
		MonsterConfig: monsterConfig,
		DungeonConfig: dungeonConfig,
	}

	return d
}
