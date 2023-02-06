package record

// Query record types

type DungeonInstanceCapacity struct {
	DungeonInstanceID             string `db:"dungeon_instance_id"`
	DungeonInstanceCharacterCount int    `db:"dungeon_instance_character_count"`
	DungeonID                     string `db:"dungeon_id"`
	DungeonLocationCount          int    `db:"dungeon_location_count"`
}

type DungeonEntityInstanceTurn struct {
	DungeonInstanceID         string `db:"dungeon_instance_id"`
	DungeonName               string `db:"dungeon_name"`
	DungeonInstanceTurnNumber int    `db:"dungeon_instance_turn_number"`
	EntityType                string `db:"entity_type"`
	EntityInstanceID          string `db:"entity_instance_id"`
	EntityName                string `db:"entity_name"`
	EntityInstanceTurnNumber  int    `db:"entity_instance_turn_number"`
}
