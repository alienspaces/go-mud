package dungeonactionobject

var createOneSQL = `
INSERT INTO dungeon_action_object (
	id,
	record_type,
	dungeon_action_id,
	dungeon_location_id,
	dungeon_object_id,
	name,
	description,
	is_stashed,
	is_equipped,
	created_at
) VALUES (
	:id,
	:record_type,
	:dungeon_action_id,
	:dungeon_location_id,
	:dungeon_object_id,
	:name,
	:description,
	:is_stashed,
	:is_equipped,
	:created_at
)
RETURNING *
`

var updateOneSQL = `
UPDATE dungeon_action_object SET
	record_type         = :record_type,
	dungeon_action_id   = :dungeon_action_id,
	dungeon_location_id = :dungeon_location_id,
	dungeon_object_id   = :dungeon_object_id,
	name                = :name,
	description         = :description,
	is_stashed          = :is_stashed,
	is_equipped         = :is_equipped,
	updated_at 		    = :updated_at
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
