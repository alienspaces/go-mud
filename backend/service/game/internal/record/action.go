package record

import (
	"database/sql"

	"gitlab.com/alienspaces/go-mud/backend/core/repository"
)

const (
	ActionCommandMove   string = "move"
	ActionCommandLook   string = "look"
	ActionCommandUse    string = "use"
	ActionCommandEquip  string = "equip"
	ActionCommandStash  string = "stash"
	ActionCommandDrop   string = "drop"
	ActionCommandAttack string = "attack"
)

type Action struct {
	DungeonInstanceID                 string         `db:"dungeon_instance_id"`
	LocationInstanceID                string         `db:"location_instance_id"`
	CharacterInstanceID               sql.NullString `db:"character_instance_id"`
	MonsterInstanceID                 sql.NullString `db:"monster_instance_id"`
	SerialNumber                      sql.NullInt16  `db:"serial_number,readonly"`
	TurnNumber                        int            `db:"turn_number"`
	ResolvedCommand                   string         `db:"resolved_command"`
	ResolvedEquippedObjectInstanceID  sql.NullString `db:"resolved_equipped_object_instance_id"`
	ResolvedStashedObjectInstanceID   sql.NullString `db:"resolved_stashed_object_instance_id"`
	ResolvedDroppedObjectInstanceID   sql.NullString `db:"resolved_dropped_object_instance_id"`
	ResolvedTargetObjectInstanceID    sql.NullString `db:"resolved_target_object_instance_id"`
	ResolvedTargetCharacterInstanceID sql.NullString `db:"resolved_target_character_instance_id"`
	ResolvedTargetMonsterInstanceID   sql.NullString `db:"resolved_target_monster_instance_id"`
	ResolvedTargetLocationDirection   sql.NullString `db:"resolved_target_location_direction"`
	ResolvedTargetLocationInstanceID  sql.NullString `db:"resolved_target_location_instance_id"`
	repository.Record
}

const (
	ActionCharacterRecordTypeSource          string = "source"
	ActionCharacterRecordTypeTarget          string = "target"
	ActionCharacterRecordTypeCurrentLocation string = "current_location"
	ActionCharacterRecordTypeTargetLocation  string = "target_location"
)

type ActionCharacter struct {
	RecordType          string `db:"record_type"`
	ActionID            string `db:"action_id"`
	LocationInstanceID  string `db:"location_instance_id"`
	CharacterInstanceID string `db:"character_instance_id"`
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

type ActionCharacterObject struct {
	ActionCharacterID string `db:"action_character_id"`
	ObjectInstanceID  string `db:"object_instance_id"`
	Name              string `db:"name"`
	IsStashed         bool   `db:"is_stashed"`
	IsEquipped        bool   `db:"is_equipped"`
	repository.Record
}

const (
	ActionMonsterRecordTypeSource          string = "source"
	ActionMonsterRecordTypeTarget          string = "target"
	ActionMonsterRecordTypeCurrentLocation string = "current_location"
	ActionMonsterRecordTypeTargetLocation  string = "target_location"
)

type ActionMonster struct {
	RecordType          string `db:"record_type"`
	ActionID            string `db:"action_id"`
	LocationInstanceID  string `db:"location_instance_id"`
	MonsterInstanceID   string `db:"monster_instance_id"`
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

type ActionMonsterObject struct {
	ActionMonsterID  string `db:"action_monster_id"`
	ObjectInstanceID string `db:"object_instance_id"`
	Name             string `db:"name"`
	IsStashed        bool   `db:"is_stashed"`
	IsEquipped       bool   `db:"is_equipped"`
	repository.Record
}

const (
	// Equipped objects are being worn or held
	ActionObjectRecordTypeEquipped string = "equipped"
	// Stashed objects are packed in a bag or backback
	ActionObjectRecordTypeStashed string = "stashed"
	// Dropped objects are present at a location
	ActionObjectRecordTypeDropped string = "dropped"
	// Target objects are are being actively looked at, used, equipped or stashed
	ActionObjectRecordTypeTarget string = "target"
	// Current location objects
	ActionObjectRecordTypeCurrentLocation string = "current_location"
	// Target location objects
	ActionObjectRecordTypeTargetLocation string = "target_location"
)

type ActionObject struct {
	RecordType         string `db:"record_type"`
	ActionID           string `db:"action_id"`
	LocationInstanceID string `db:"location_instance_id"`
	ObjectInstanceID   string `db:"object_instance_id"`
	Name               string `db:"name"`
	// Description could be either the object `description` or `description_detailed`
	// depending on the characters `look` action result.
	Description string `db:"description"`
	IsStashed   bool   `db:"is_stashed"`
	IsEquipped  bool   `db:"is_equipped"`
	repository.Record
}
