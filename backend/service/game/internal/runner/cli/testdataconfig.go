package runner

import (
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/data/cabin"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/data/cave"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/harness"
)

// The following data is test data used to seed a test game server
func TestDataConfig() harness.DataConfig {

	var objectConfig = cave.ObjectConfig()
	objectConfig = append(objectConfig, cabin.ObjectConfig()...)

	var monsterConfig = cave.MonsterConfig()
	monsterConfig = append(monsterConfig, cabin.MonsterConfig()...)

	var characterConfig = cave.CharacterConfig()
	characterConfig = append(characterConfig, cabin.CharacterConfig()...)

	var dungeonConfig = []harness.DungeonConfig{
		cave.DungeonConfig(),
		cabin.DungeonConfig(),
	}

	d := harness.DataConfig{
		ObjectConfig:    objectConfig,
		MonsterConfig:   monsterConfig,
		CharacterConfig: characterConfig,
		DungeonConfig:   dungeonConfig,
	}

	return d
}
