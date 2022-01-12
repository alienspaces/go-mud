package record

import (
	"database/sql"

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
	DungeonID                  string         `db:"dungeon_id"`
	Name                       string         `db:"name"`
	Description                string         `db:"description"`
	Default                    bool           `db:"default"`
	NorthDungeonLocationID     sql.NullString `db:"north_dungeon_location_id"`
	NortheastDungeonLocationID sql.NullString `db:"northeast_dungeon_location_id"`
	EastDungeonLocationID      sql.NullString `db:"east_dungeon_location_id"`
	SoutheastDungeonLocationID sql.NullString `db:"southeast_dungeon_location_id"`
	SouthDungeonLocationID     sql.NullString `db:"south_dungeon_location_id"`
	SouthwestDungeonLocationID sql.NullString `db:"southwest_dungeon_location_id"`
	WestDungeonLocationID      sql.NullString `db:"west_dungeon_location_id"`
	NorthwestDungeonLocationID sql.NullString `db:"northwest_dungeon_location_id"`
	UpDungeonLocationID        sql.NullString `db:"up_dungeon_location_id"`
	DownDungeonLocationID      sql.NullString `db:"down_dungeon_location_id"`
	repository.Record
}

type DungeonCharacter struct {
	DungeonID           string `db:"dungeon_id"`
	DungeonLocationID   string `db:"dungeon_location_id"`
	Name                string `db:"name"`
	Strength            int    `db:"strength"`
	Dexterity           int    `db:"dexterity"`
	Intelligence        int    `db:"intelligence"`
	CurrentStrength     int    `db:"current_strength"`
	CurrentDexterity    int    `db:"current_dexterity"`
	CurrentIntelligence int    `db:"current_intelligence"`
	Health              int    `db:"health"`
	Fatigue             int    `db:"fatigue"`
	CurrentHealth       int    `db:"current_health"`
	CurrentFatigue      int    `db:"current_fatigue"`
	Coins               int    `db:"coins"`
	ExperiencePoints    int    `db:"experience_points"`
	AttributePoints     int    `db:"attribute_points"`
	repository.Record
}

type DungeonMonster struct {
	DungeonID           string `db:"dungeon_id"`
	DungeonLocationID   string `db:"dungeon_location_id"`
	Name                string `db:"name"`
	Strength            int    `db:"strength"`
	Dexterity           int    `db:"dexterity"`
	Intelligence        int    `db:"intelligence"`
	CurrentStrength     int    `db:"current_strength"`
	CurrentDexterity    int    `db:"current_dexterity"`
	CurrentIntelligence int    `db:"current_intelligence"`
	Health              int    `db:"health"`
	Fatigue             int    `db:"fatigue"`
	CurrentHealth       int    `db:"current_health"`
	CurrentFatigue      int    `db:"current_fatigue"`
	Coins               int    `db:"coins"`
	ExperiencePoints    int    `db:"experience_points"`
	AttributePoints     int    `db:"attribute_points"`
	repository.Record
}

type DungeonObject struct {
	DungeonID           string         `db:"dungeon_id"`
	DungeonLocationID   sql.NullString `db:"dungeon_location_id"`
	DungeonCharacterID  sql.NullString `db:"dungeon_character_id"`
	DungeonMonsterID    sql.NullString `db:"dungeon_monster_id"`
	Name                string         `db:"name"`
	Description         string         `db:"description"`
	DescriptionDetailed string         `db:"description_detailed"`
	IsStashed           bool           `db:"is_stashed"`
	IsEquipped          bool           `db:"is_equipped"`
	repository.Record
}

type DungeonAction struct {
	DungeonID                              string         `db:"dungeon_id"`
	DungeonLocationID                      string         `db:"dungeon_location_id"`
	DungeonCharacterID                     sql.NullString `db:"dungeon_character_id"`
	DungeonMonsterID                       sql.NullString `db:"dungeon_monster_id"`
	SerialID                               sql.NullInt16  `db:"serial_id"`
	ResolvedCommand                        string         `db:"resolved_command"`
	ResolvedEquippedDungeonObjectID        sql.NullString `db:"resolved_equipped_dungeon_object_id"`
	ResolvedStashedDungeonObjectID         sql.NullString `db:"resolved_stashed_dungeon_object_id"`
	ResolvedTargetDungeonObjectID          sql.NullString `db:"resolved_target_dungeon_object_id"`
	ResolvedTargetDungeonCharacterID       sql.NullString `db:"resolved_target_dungeon_character_id"`
	ResolvedTargetDungeonMonsterID         sql.NullString `db:"resolved_target_dungeon_monster_id"`
	ResolvedTargetDungeonLocationDirection sql.NullString `db:"resolved_target_dungeon_location_direction"`
	ResolvedTargetDungeonLocationID        sql.NullString `db:"resolved_target_dungeon_location_id"`
	repository.Record
}

const (
	DungeonActionCharacterRecordTypeSource   string = "source"
	DungeonActionCharacterRecordTypeTarget   string = "target"
	DungeonActionCharacterRecordTypeOccupant string = "occupant"
)

type DungeonActionCharacter struct {
	RecordType          string `db:"record_type"`
	DungeonActionID     string `db:"dungeon_action_id"`
	DungeonLocationID   string `db:"dungeon_location_id"`
	DungeonCharacterID  string `db:"dungeon_character_id"`
	Name                string `db:"name"`
	Strength            int    `db:"strength"`
	Dexterity           int    `db:"dexterity"`
	Intelligence        int    `db:"intelligence"`
	CurrentStrength     int    `db:"current_strength"`
	CurrentDexterity    int    `db:"current_dexterity"`
	CurrentIntelligence int    `db:"current_intelligence"`
	Health              int    `db:"health"`
	Fatigue             int    `db:"fatigue"`
	CurrentHealth       int    `db:"current_health"`
	CurrentFatigue      int    `db:"current_fatigue"`
	repository.Record
}

const (
	DungeonActionMonsterRecordTypeSource   string = "source"
	DungeonActionMonsterRecordTypeTarget   string = "target"
	DungeonActionMonsterRecordTypeOccupant string = "occupant"
)

type DungeonActionMonster struct {
	RecordType          string `db:"record_type"`
	DungeonActionID     string `db:"dungeon_action_id"`
	DungeonLocationID   string `db:"dungeon_location_id"`
	DungeonMonsterID    string `db:"dungeon_monster_id"`
	Name                string `db:"name"`
	Strength            int    `db:"strength"`
	Dexterity           int    `db:"dexterity"`
	Intelligence        int    `db:"intelligence"`
	CurrentStrength     int    `db:"current_strength"`
	CurrentDexterity    int    `db:"current_dexterity"`
	CurrentIntelligence int    `db:"current_intelligence"`
	Health              int    `db:"health"`
	Fatigue             int    `db:"fatigue"`
	CurrentHealth       int    `db:"current_health"`
	CurrentFatigue      int    `db:"current_fatigue"`
	repository.Record
}

const (
	// Equipped objects are being worn or held
	DungeonActionObjectRecordTypeEquipped string = "equipped"
	// Stashed objects are packed in a bag or backback
	DungeonActionObjectRecordTypeStashed string = "stashed"
	// Target objects are are being actively looked at, used, equipped or stashed
	DungeonActionObjectRecordTypeTarget string = "target"
	// Occupant objects are present at a location
	DungeonActionObjectRecordTypeOccupant string = "occupant"
)

type DungeonActionObject struct {
	RecordType        string `db:"record_type"`
	DungeonActionID   string `db:"dungeon_action_id"`
	DungeonLocationID string `db:"dungeon_location_id"`
	DungeonObjectID   string `db:"dungeon_object_id"`
	Name              string `db:"name"`
	// Description could be either the object `description` or `description_detailed`
	// depending on the characters `look` action result.
	Description string `db:"description"`
	IsStashed   bool   `db:"is_stashed"`
	IsEquipped  bool   `db:"is_equipped"`
	repository.Record
}
