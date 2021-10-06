package player

var createOneSQL = `
INSERT INTO player (
    id,
    name,
    email,
    provider,
    provider_account_id,
    created_at
) VALUES (
    :id,
    :name,
    :email,
    :provider,
    :provider_account_id,
    :created_at
)
RETURNING *
`

var updateOneSQL = `
UPDATE player SET
    name                = :name,
    email               = :email,
    provider            = :provider,
    provider_account_id = :provider_account_id,
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
