package dungeonmonster

var createOneSQL = `
INSERT INTO dungeon_monster (
	id,
	dungeon_id,
	dungeon_location_id,
	name,
	strength,
	dexterity,
	intelligence,
	health,
	fatigue,
	coins,
	experience_points,
	attribute_points,
	created_at
) VALUES (
	:id,
	:dungeon_id,
	:dungeon_location_id,
	:name,
	:strength,
	:dexterity,
	:intelligence,
	:health,
	:fatigue,
	:coins,
	:experience_points,
	:attribute_points,
	:created_at
)
RETURNING *
`

var updateOneSQL = `
UPDATE dungeon_monster SET
	dungeon_id 			= :dungeon_id,
	dungeon_location_id = :dungeon_location_id,
	name 				= :name,
	strength 			= :strength,
	dexterity 			= :dexterity,
	intelligence 		= :intelligence,
	health 				= :health,
	fatigue 			= :fatigue,
	coins 				= :coins,
	experience_points 	= :experience_points,
	attribute_points 	= :attribute_points,
    updated_at          = :updated_at
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
