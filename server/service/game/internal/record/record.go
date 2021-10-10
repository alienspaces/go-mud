package record

import (
	"gitlab.com/alienspaces/go-mud/server/core/repository"
)

// Dungeon -
type Dungeon struct {
	Name        string `db:"name"`
	Description string `db:"description"`
	repository.Record
}

// DungeonLocation -
type DungeonLocation struct {
	DungeonID                  string `db:"dungeon_id"`
	Name                       string `db:"name"`
	Description                string `db:"description"`
	Default                    bool   `db:"default"`
	NorthDungeonLocationID     string `db:"north_dungeon_location_id"`
	NortheastDungeonLocationID string `db:"northeast_dungeon_location_id"`
	EastDungeonLocationID      string `db:"east_dungeon_location_id"`
	SoutheastDungeonLocationID string `db:"southeast_dungeon_location_id"`
	SouthDungeonLocationID     string `db:"south_dungeon_location_id"`
	SouthwestDungeonLocationID string `db:"southwest_dungeon_location_id"`
	WestDungeonLocationID      string `db:"west_dungeon_location_id"`
	NorthwestDungeonLocationID string `db:"northwest_dungeon_location_id"`
	UpDungeonLocationID        string `db:"up_dungeon_location_id"`
	DownDungeonLocationID      string `db:"down_dungeon_location_id"`
	repository.Record
}

type DungeonCharacter struct {
	DungeonID         string `db:"dungeon_id"`
	DungeonLocationID string `db:"dungeon_location_id"`
	Name              string `db:"name"`
	Strength          int    `db:"strength"`
	Dexterity         int    `db:"dexterity"`
	Intelligence      int    `db:"intelligence"`
	Health            int    `db:"health"`
	Fatigue           int    `db:"fatigue"`
	Coins             int    `db:"coins"`
	ExperiencePoints  int    `db:"experience_points"`
	AttributePoints   int    `db:"attribute_points"`
	repository.Record
}

type DungeonMonster struct {
	DungeonID         string `db:"dungeon_id"`
	DungeonLocationID string `db:"dungeon_location_id"`
	Name              string `db:"name"`
	Strength          int    `db:"strength"`
	Dexterity         int    `db:"dexterity"`
	Intelligence      int    `db:"intelligence"`
	Health            int    `db:"health"`
	Fatigue           int    `db:"fatigue"`
	Coins             int    `db:"coins"`
	ExperiencePoints  int    `db:"experience_points"`
	AttributePoints   int    `db:"attribute_points"`
	repository.Record
}

type DungeonObject struct {
	DungeonID          string `db:"dungeon_id"`
	DungeonLocationID  string `db:"dungeon_location_id"`
	DungeonCharacterID string `db:"dungeon_character_id"`
	DungeonMonsterID   string `db:"dungeon_monster_id"`
	Name               string `db:"name"`
	Description        string `db:"description"`
	DescriptionLong    string `db:"description_long"`
	IsStashed          bool   `db:"is_stashed"`
	IsEquipped         bool   `db:"is_equipped"`
	repository.Record
}

type DungeonAction struct {
	DungeonID                              string `db:"dungeon_id"`
	DungeonLocationID                      string `db:"dungeon_location_id"`
	DungeonCharacterID                     string `db:"dungeon_character_id"`
	DungeonMonsterID                       string `db:"dungeon_monster_id"`
	SerialID                               string `db:"serial_id"`
	ResolvedCommand                        string `db:"resolved_command"`
	ResolvedEquippedDungeonObjectID        string `db:"resolved_equipped_dungeon_object_id"`
	ResolvedEquippedDungeonObjectName      string `db:"resolved_equipped_dungeon_object_name"`
	ResolvedStashedDungeonObjectID         string `db:"resolved_stashed_dungeon_object_id"`
	ResolvedStashedDungeonObjectName       string `db:"resolved_stashed_dungeon_object_name"`
	ResolvedTargetDungeonObjectID          string `db:"resolved_target_dungeon_object_id"`
	ResolvedTargetDungeonObjectName        string `db:"resolved_target_dungeon_object_name"`
	ResolvedTargetDungeonCharacterID       string `db:"resolved_target_dungeon_character_id"`
	ResolvedTargetDungeonCharacterName     string `db:"resolved_target_dungeon_character_name"`
	ResolvedTargetDungeonMonsterID         string `db:"resolved_target_dungeon_monster_id"`
	ResolvedTargetDungeonMonsterName       string `db:"resolved_target_dungeon_monster_name"`
	ResolvedTargetDungeonLocationDirection string `db:"resolved_target_dungeon_location_direction"`
	ResolvedTargetDungeonLocationID        string `db:"resolved_target_dungeon_location_id"`
	ResolvedTargetDungeonLocationName      string `db:"resolved_target_dungeon_location_name"`
	repository.Record
}

type DungeonActionCharacter struct {
	DungeonActionID    string `db:"dungeon_action_id"`
	DungeonCharacterID string `db:"dungeon_character_id"`
	repository.Record
}

type DungeonActionMonster struct {
	DungeonActionID  string `db:"dungeon_action_id"`
	DungeonMonsterID string `db:"dungeon_monster_id"`
	repository.Record
}

type DungeonActionObject struct {
	DungeonActionID string `db:"dungeon_action_id"`
	DungeonObjectID string `db:"dungeon_object_id"`
	repository.Record
}
