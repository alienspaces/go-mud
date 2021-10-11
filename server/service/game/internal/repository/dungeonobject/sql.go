package dungeonobject

var createOneSQL = `
INSERT INTO dungeon_object (
	id,
	dungeon_id,
	dungeon_location_id,
	dungeon_character_id,
	dungeon_monster_id,
	name,
	description,
	description_long,
	is_stashed,
	is_equipped,
	created_at
) VALUES (
	:id,
	:dungeon_id,
	:dungeon_location_id,
	:dungeon_character_id,
	:dungeon_monster_id,
	:name,
	:description,
	:description_long,
	:is_stashed,
	:is_equipped,
	:created_at
)
RETURNING *
`

var updateOneSQL = `
UPDATE dungeon_object SET
	dungeon_id 			 = :dungeon_id,
	dungeon_location_id  = :dungeon_location_id,
	dungeon_character_id = :dungeon_character_id,
	dungeon_monster_id 	 = :dungeon_monster_id,
	name 				 = :name,
	description 		 = :description,
	description_long 	 = :description_long,
	is_stashed 			 = :is_stashed,
	is_equipped 		 = :is_equipped,
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
