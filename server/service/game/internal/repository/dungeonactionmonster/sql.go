package dungeonactionmonster

var createOneSQL = `
INSERT INTO dungeon_action_monster (
	id,
	record_type,
	dungeon_action_id,
	dungeon_location_id,
	dungeon_monster_id,
	name,
	strength,
	dexterity,
	intelligence,
	health,
	fatigue,
	created_at
) VALUES (
	:id,
	:record_type,
	:dungeon_action_id,
	:dungeon_location_id,
	:dungeon_monster_id,
	:name,
	:strength,
	:dexterity,
	:intelligence,
	:health,
	:fatigue,
	:created_at
)
RETURNING *
`

var updateOneSQL = `
UPDATE dungeon_action_monster SET
	record_type         = :record_type,
	dungeon_action_id   = :dungeon_action_id,
	dungeon_location_id = :dungeon_location_id,
	dungeon_monster_id  = :dungeon_monster_id,
	name 				= :name,
	strength 			= :strength,
	dexterity 			= :dexterity,
	intelligence 		= :intelligence,
	health 				= :health,
	fatigue 			= :fatigue,
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
