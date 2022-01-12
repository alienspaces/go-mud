package dungeonactioncharacter

var createOneSQL = `
INSERT INTO dungeon_action_character (
	id,
	record_type,
	dungeon_action_id,
	dungeon_location_id,
	dungeon_character_id,
	name,
	strength,
	dexterity,
	intelligence,
	current_strength,
	current_dexterity,
	current_intelligence,
	health,
	fatigue,
	current_health,
	current_fatigue,
	created_at
) VALUES (
	:id,
	:record_type,
	:dungeon_action_id,
	:dungeon_location_id,
	:dungeon_character_id,
	:name,
	:strength,
	:dexterity,
	:intelligence,
	:current_strength,
	:current_dexterity,
	:current_intelligence,
	:health,
	:fatigue,
	:current_health,
	:current_fatigue,
	:created_at
)
RETURNING *
`

var updateOneSQL = `
UPDATE dungeon_action_character SET
	record_type          = :record_type,
	dungeon_action_id    = :dungeon_action_id,
	dungeon_location_id  = :dungeon_location_id,
	dungeon_character_id = :dungeon_character_id,
	name 				 = :name,
	strength 			 = :strength,
	dexterity 			 = :dexterity,
	intelligence 		 = :intelligence,
	current_strength 	 = :current_strength,
	current_dexterity 	 = :current_dexterity,
	current_intelligence = :current_intelligence,
	health 				 = :health,
	fatigue 			 = :fatigue,
	current_health       = :current_health,
	current_fatigue      = :current_fatigue,
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
