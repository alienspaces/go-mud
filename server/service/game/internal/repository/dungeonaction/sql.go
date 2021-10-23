package dungeonaction

var createOneSQL = `
INSERT INTO dungeon_action (
	id,
	dungeon_id,
	dungeon_location_id,
	dungeon_character_id,
	dungeon_monster_id,
	resolved_command,
	resolved_equipped_dungeon_object_name,
	resolved_equipped_dungeon_object_id,
	resolved_stashed_dungeon_object_name,
	resolved_stashed_dungeon_object_id,
	resolved_target_dungeon_object_name,
	resolved_target_dungeon_object_id,
	resolved_target_dungeon_character_name,
	resolved_target_dungeon_character_id,
	resolved_target_dungeon_monster_name,
	resolved_target_dungeon_monster_id,
	resolved_target_dungeon_location_direction,
	resolved_target_dungeon_location_name,
	resolved_target_dungeon_location_id,
	created_at
) VALUES (
	:id,
	:dungeon_id,
	:dungeon_location_id,
	:dungeon_character_id,
	:dungeon_monster_id,
	:resolved_command,
	:resolved_equipped_dungeon_object_name,
	:resolved_equipped_dungeon_object_id,
	:resolved_stashed_dungeon_object_name,
	:resolved_stashed_dungeon_object_id,
	:resolved_target_dungeon_object_name,
	:resolved_target_dungeon_object_id,
	:resolved_target_dungeon_character_name,
	:resolved_target_dungeon_character_id,
	:resolved_target_dungeon_monster_name,
	:resolved_target_dungeon_monster_id,
	:resolved_target_dungeon_location_direction,
	:resolved_target_dungeon_location_name,
	:resolved_target_dungeon_location_id,
	:created_at
)
RETURNING *
`

var updateOneSQL = `
UPDATE dungeon_action SET
	dungeon_id 									= :dungeon_id,
	dungeon_location_id 						= :dungeon_location_id,
	dungeon_character_id 						= :dungeon_character_id,
	dungeon_monster_id 							= :dungeon_monster_id,
	serial_id 									= :serial_id,
	resolved_command 							= :resolved_command,
	resolved_equipped_dungeon_object_name 		= :resolved_equipped_dungeon_object_name,
	resolved_equipped_dungeon_object_id 		= :resolved_equipped_dungeon_object_id,
	resolved_stashed_dungeon_object_name 		= :resolved_stashed_dungeon_object_name,
	resolved_stashed_dungeon_object_id 			= :resolved_stashed_dungeon_object_id,
	resolved_target_dungeon_object_name 		= :resolved_target_dungeon_object_name,
	resolved_target_dungeon_object_id 			= :resolved_target_dungeon_object_id,
	resolved_target_dungeon_character_name 		= :resolved_target_dungeon_character_name,
	resolved_target_dungeon_character_id 		= :resolved_target_dungeon_character_id,
	resolved_target_dungeon_monster_name 		= :resolved_target_dungeon_monster_name,
	resolved_target_dungeon_monster_id 			= :resolved_target_dungeon_monster_id,
	resolved_target_dungeon_location_direction 	= :resolved_target_dungeon_location_direction,
	resolved_target_dungeon_location_name 		= :resolved_target_dungeon_location_name,
	resolved_target_dungeon_location_id 		= :resolved_target_dungeon_location_id,
	updated_at 		     						= :updated_at
WHERE id = :id
AND   deleted_at IS NULL
RETURNING *
`

// CreateOneSQL -
func (r *Repository) CreateOneSQL() string {
	return createOneSQL
}

// UpdateOneSQL -
func (r *Repository) UpdateOneSQL() string {
	return updateOneSQL
}

// OrderBy -
func (r *Repository) OrderBy() string {
	return "created_at desc"
}
