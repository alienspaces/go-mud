package record

import "gitlab.com/alienspaces/go-mud/server/core/repository"

type Character struct {
	Name             string `db:"name"`
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

type CharacterInstance struct {
	CharacterID               string `db:"character_id"`
	DungeonInstanceID         string `db:"dungeon_instance_id"`
	DungeonLocationInstanceID string `db:"dungeon_location_instance_id"`
	Strength                  int    `db:"strength"`
	Dexterity                 int    `db:"dexterity"`
	Intelligence              int    `db:"intelligence"`
	Health                    int    `db:"health"`
	Fatigue                   int    `db:"fatigue"`
	Coins                     int    `db:"coins"`
	ExperiencePoints          int    `db:"experience_points"`
	AttributePoints           int    `db:"attribute_points"`
	repository.Record
}
