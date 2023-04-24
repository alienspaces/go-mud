package record

import (
	"database/sql"

	"gitlab.com/alienspaces/go-mud/backend/core/repository"
)

type Monster struct {
	Name             string `db:"name"`
	Description      string `db:"description"`
	Strength         int    `db:"strength"`
	Dexterity        int    `db:"dexterity"`
	Intelligence     int    `db:"intelligence"`
	Health           int    `db:"health"`
	Fatigue          int    `db:"fatigue"`
	Coins            int    `db:"coins"`
	ExperiencePoints int    `db:"experience_points"`
	AttributePoints  int    `db:"attribute_points"`
	repository.Record
}

type MonsterObject struct {
	MonsterID  string `db:"monster_id"`
	ObjectID   string `db:"object_id"`
	IsStashed  bool   `db:"is_stashed"`
	IsEquipped bool   `db:"is_equipped"`
	repository.Record
}

type MonsterInstance struct {
	MonsterID          string `db:"monster_id"`
	DungeonInstanceID  string `db:"dungeon_instance_id"`
	LocationInstanceID string `db:"location_instance_id"`
	Strength           int    `db:"strength"`
	Dexterity          int    `db:"dexterity"`
	Intelligence       int    `db:"intelligence"`
	Health             int    `db:"health"`
	Fatigue            int    `db:"fatigue"`
	Coins              int    `db:"coins"`
	ExperiencePoints   int    `db:"experience_points"`
	AttributePoints    int    `db:"attribute_points"`
	repository.Record
}

type MonsterInstanceMemory struct {
	MonsterInstanceID         string         `db:"monster_instance_id"`
	MemoryCommand             string         `db:"memory_command"`
	MemoryType                string         `db:"memory_type"`
	TurnNumber                int            `db:"turn_number"`
	MemoryLocationInstanceID  sql.NullString `db:"memory_location_instance_id"`
	MemoryCharacterInstanceID sql.NullString `db:"memory_character_instance_id"`
	MemoryMonsterInstanceID   sql.NullString `db:"memory_monster_instance_id"`
	MemoryObjectInstanceID    sql.NullString `db:"memory_object_instance_id"`
	repository.Record
}

type MonsterInstanceView struct {
	MonsterID           string `db:"monster_id"`
	DungeonInstanceID   string `db:"dungeon_instance_id"`
	LocationInstanceID  string `db:"location_instance_id"`
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
