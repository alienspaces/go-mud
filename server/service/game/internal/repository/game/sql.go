package game

var createOneSQL = `
INSERT INTO game (
	id,
	created_at
) VALUES (
	:id,
	:created_at
)
RETURNING *
`

var updateOneSQL = `
UPDATE game SET
    updated_at = :updated_at
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
