package harness

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// Data -
type Data struct {
	// Object
	ObjectRecs []*record.Object

	// Monster
	MonsterRecs       []*record.Monster
	MonsterObjectRecs []*record.MonsterObject

	// Character
	CharacterRecs       []*record.Character
	CharacterObjectRecs []*record.CharacterObject

	// Dungeon
	DungeonRecs         []*record.Dungeon
	LocationRecs        []*record.Location
	LocationObjectRecs  []*record.LocationObject
	LocationMonsterRecs []*record.LocationMonster

	// Dungeon Instance
	DungeonInstanceRecs   []*record.DungeonInstance
	LocationInstanceRecs  []*record.LocationInstance
	CharacterInstanceRecs []*record.CharacterInstance
	MonsterInstanceRecs   []*record.MonsterInstance
	ObjectInstanceRecs    []*record.ObjectInstance

	// Action
	ActionRecs                []*record.Action
	ActionCharacterRecs       []*record.ActionCharacter
	ActionCharacterObjectRecs []*record.ActionCharacterObject
	ActionMonsterRecs         []*record.ActionMonster
	ActionMonsterObjectRecs   []*record.ActionMonsterObject
	ActionObjectRecs          []*record.ActionObject
}

// teardownData -
type teardownData struct {
	// Object
	ObjectRecs []record.Object

	// Monster
	MonsterRecs       []record.Monster
	MonsterObjectRecs []record.MonsterObject

	// Character
	CharacterRecs       []record.Character
	CharacterObjectRecs []record.CharacterObject

	// Dungeon
	DungeonRecs         []record.Dungeon
	LocationRecs        []record.Location
	LocationObjectRecs  []record.LocationObject
	LocationMonsterRecs []record.LocationMonster

	// Dungeon Instance
	DungeonInstanceRecs   []*record.DungeonInstance
	LocationInstanceRecs  []*record.LocationInstance
	CharacterInstanceRecs []*record.CharacterInstance
	MonsterInstanceRecs   []*record.MonsterInstance
	ObjectInstanceRecs    []*record.ObjectInstance

	// Action
	ActionRecs                []*record.Action
	ActionCharacterRecs       []*record.ActionCharacter
	ActionCharacterObjectRecs []*record.ActionCharacterObject
	ActionMonsterRecs         []*record.ActionMonster
	ActionMonsterObjectRecs   []*record.ActionMonsterObject
	ActionObjectRecs          []*record.ActionObject
}

func (d *Data) GetObjectRecByName(objectName string) (*record.Object, error) {

	for idx := range d.ObjectRecs {
		if d.ObjectRecs[idx].Name == objectName {
			return d.ObjectRecs[idx], nil
		}
	}

	return nil, fmt.Errorf("unknown object name >%s<", objectName)
}

func (d *Data) GetMonsterRecByName(monsterName string) (*record.Monster, error) {

	for idx := range d.MonsterRecs {
		if d.MonsterRecs[idx].Name == monsterName {
			return d.MonsterRecs[idx], nil
		}
	}

	return nil, fmt.Errorf("unknown monster name >%s<", monsterName)
}
