package dungeonlocation

var createOneSQL = `
INSERT INTO dungeon_location (
	id,
	dungeon_id,
	name,
	description,
	default,
	north_dungeon_location_id,
	northeast_dungeon_location_id,
	east_dungeon_location_id,
	southeast_dungeon_location_id,
	south_dungeon_location_id,
	southwest_dungeon_location_id,
	west_dungeon_location_id,
	northwest_dungeon_location_id,
	up_dungeon_location_id,
	down_dungeon_location_id,
	created_at
) VALUES (
	:id,
	:dungeon_id,
	:name,
	:description,
	:default,
	:north_dungeon_location_id,
	:northeast_dungeon_location_id,
	:east_dungeon_location_id,
	:southeast_dungeon_location_id,
	:south_dungeon_location_id,
	:southwest_dungeon_location_id,
	:west_dungeon_location_id,
	:northwest_dungeon_location_id,
	:up_dungeon_location_id,
	:down_dungeon_location_id,
	:created_at
)
RETURNING *
`

var updateOneSQL = `
UPDATE dungeon_location SET
	dungeon_id 			          = :dungeon_id,
	name 				          = :name,
	description 		          = :description,
	default 					  = :default,
	north_dungeon_location_id 	  = :north_dungeon_location_id,
	northeast_dungeon_location_id = :northeast_dungeon_location_id,
	east_dungeon_location_id      = :east_dungeon_location_id,
	southeast_dungeon_location_id = :southeast_dungeon_location_id,
	south_dungeon_location_id     = :south_dungeon_location_id,
	southwest_dungeon_location_id = :southwest_dungeon_location_id,
	west_dungeon_location_id      = :west_dungeon_location_id,
	northwest_dungeon_location_id = :northwest_dungeon_location_id,
	up_dungeon_location_id        = :up_dungeon_location_id,
	down_dungeon_location_id      = :down_dungeon_location_id,
    updated_at                    = :updated_at
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
