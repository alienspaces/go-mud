package playerrole

var createOneSQL = `
INSERT INTO player_role (
    id,
    player_id,
    role,
    created_at
) VALUES (
    :id,
    :player_id,
    :role,
    :created_at
)
RETURNING *
`

var updateOneSQL = `
UPDATE player_role SET
    player_id   = :player_id,
    role         = :role,
    updated_at   = :updated_at
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
