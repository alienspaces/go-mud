package record

import (
	"database/sql"

	"gitlab.com/alienspaces/go-mud/server/core/repository"
)

const (
	DungeonActionCommandMove   string = "move"
	DungeonActionCommandLook   string = "look"
	DungeonActionCommandEquip  string = "equip"
	DungeonActionCommandStash  string = "stash"
	DungeonActionCommandDrop   string = "drop"
	DungeonActionCommandAttack string = "attack"
)

type DungeonAction struct {
	DungeonInstanceID                      string         `db:"dungeon_instance_id"`
	DungeonLocationInstanceID              string         `db:"dungeon_location_instance_id"`
	CharacterInstanceID                    sql.NullString `db:"character_instance_id"`
	MonsterInstanceID                      sql.NullString `db:"monster_instance_id"`
	SerialID                               sql.NullInt16  `db:"serial_id"`
	ResolvedCommand                        string         `db:"resolved_command"`
	ResolvedEquippedDungeonObjectID        sql.NullString `db:"resolved_equipped_dungeon_object_id"`
	ResolvedStashedDungeonObjectID         sql.NullString `db:"resolved_stashed_dungeon_object_id"`
	ResolvedDroppedDungeonObjectID         sql.NullString `db:"resolved_dropped_dungeon_object_id"`
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
	RecordType                string `db:"record_type"`
	DungeonActionID           string `db:"dungeon_action_id"`
	DungeonLocationInstanceID string `db:"dungeon_location_instance_id"`
	CharacterInstanceID       string `db:"character_instance_id"`
	Name                      string `db:"name"`
	Strength                  int    `db:"strength"`
	Dexterity                 int    `db:"dexterity"`
	Intelligence              int    `db:"intelligence"`
	CurrentStrength           int    `db:"current_strength"`
	CurrentDexterity          int    `db:"current_dexterity"`
	CurrentIntelligence       int    `db:"current_intelligence"`
	Health                    int    `db:"health"`
	Fatigue                   int    `db:"fatigue"`
	CurrentHealth             int    `db:"current_health"`
	CurrentFatigue            int    `db:"current_fatigue"`
	repository.Record
}

type DungeonActionCharacterObject struct {
	DungeonActionID     string `db:"dungeon_action_id"`
	CharacterInstanceID string `db:"character_instance_id"`
	ObjectInstanceID    string `db:"object_instance_id"`
	Name                string `db:"name"`
	IsStashed           bool   `db:"is_stashed"`
	IsEquipped          bool   `db:"is_equipped"`
	repository.Record
}

const (
	DungeonActionMonsterRecordTypeSource   string = "source"
	DungeonActionMonsterRecordTypeTarget   string = "target"
	DungeonActionMonsterRecordTypeOccupant string = "occupant"
)

type DungeonActionMonster struct {
	RecordType                string `db:"record_type"`
	DungeonActionID           string `db:"dungeon_action_id"`
	DungeonLocationInstanceID string `db:"dungeon_location_instance_id"`
	MonsterInstanceID         string `db:"monster_instance_id"`
	Name                      string `db:"name"`
	Strength                  int    `db:"strength"`
	Dexterity                 int    `db:"dexterity"`
	Intelligence              int    `db:"intelligence"`
	CurrentStrength           int    `db:"current_strength"`
	CurrentDexterity          int    `db:"current_dexterity"`
	CurrentIntelligence       int    `db:"current_intelligence"`
	Health                    int    `db:"health"`
	Fatigue                   int    `db:"fatigue"`
	CurrentHealth             int    `db:"current_health"`
	CurrentFatigue            int    `db:"current_fatigue"`
	repository.Record
}

type DungeonActionMonsterObject struct {
	DungeonActionID   string `db:"dungeon_action_id"`
	MonsterInstanceID string `db:"dungeon_monster_instance_id"`
	ObjectInstanceID  string `db:"dungeon_object_instance_id"`
	Name              string `db:"name"`
	IsStashed         bool   `db:"is_stashed"`
	IsEquipped        bool   `db:"is_equipped"`
	repository.Record
}

const (
	// Equipped objects are being worn or held
	DungeonActionObjectRecordTypeEquipped string = "equipped"
	// Stashed objects are packed in a bag or backback
	DungeonActionObjectRecordTypeStashed string = "stashed"
	// Dropped objects are present at a location
	DungeonActionObjectRecordTypeDropped string = "dropped"
	// Target objects are are being actively looked at, used, equipped or stashed
	DungeonActionObjectRecordTypeTarget string = "target"
	// Occupant objects are present at a location
	DungeonActionObjectRecordTypeOccupant string = "occupant"
)

type DungeonActionObject struct {
	RecordType                string `db:"record_type"`
	DungeonActionID           string `db:"dungeon_action_id"`
	DungeonLocationInstanceID string `db:"dungeon_location_instance_id"`
	ObjectInstanceID          string `db:"object_instance_id"`
	Name                      string `db:"name"`
	// Description could be either the object `description` or `description_detailed`
	// depending on the characters `look` action result.
	Description string `db:"description"`
	IsStashed   bool   `db:"is_stashed"`
	IsEquipped  bool   `db:"is_equipped"`
	repository.Record
}
