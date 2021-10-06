package record

import (
	"gitlab.com/alienspaces/go-boilerplate/server/core/repository"
)

// Character -
type Character struct {
	repository.Record
	PlayerID         string `db:"player_id"`
	Name             string `db:"name"`
	Avatar           string `db:"avatar"`
	Strength         int    `db:"strength"`
	Dexterity        int    `db:"dexterity"`
	Intelligence     int    `db:"intelligence"`
	AttributePoints  int64  `db:"attribute_points"`
	ExperiencePoints int64  `db:"experience_points"`
	Coins            int64  `db:"coins"`
}
