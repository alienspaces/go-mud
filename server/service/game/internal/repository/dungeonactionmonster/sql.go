package dungeonactionmonster

var createOneSQL = `
INSERT INTO dungeon_action_monster (
	id,
	dungeon_action_id,
	dungeon_monster_id,
	created_at
) VALUES (
	:id,
	:dungeon_action_id,
	:dungeon_monster_id,
	:created_at
)
RETURNING *
`

var updateOneSQL = `
UPDATE dungeon_action_monster SET
	dungeon_action_id  = :dungeon_action_id,
	dungeon_monster_id = :dungeon_monster_id,
	updated_at 		   = :updated_at
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
