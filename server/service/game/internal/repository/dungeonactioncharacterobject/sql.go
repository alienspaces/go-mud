package dungeonactioncharacterobject

var createOneSQL = `
INSERT INTO dungeon_action_character_object (
	id,
	record_type,
	dungeon_action_id,
	dungeon_character_id,
	dungeon_object_id,
	created_at
) VALUES (
	:id,
	:record_type,
	:dungeon_action_id,
	:dungeon_character_id,
	:dungeon_object_id,
	:created_at
)
RETURNING *
`

var updateOneSQL = `
UPDATE dungeon_action_character_object SET
	record_type          = :record_type,
	dungeon_action_id    = :dungeon_action_id,
	dungeon_character_id = :dungeon_character_id,
	dungeon_object_id    = :dungeon_object_id,
	updated_at 		     = :updated_at
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
