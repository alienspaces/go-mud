package schema

import "gitlab.com/alienspaces/go-mud/backend/schema"

// CharacterResponse -
type CharacterResponse struct {
	schema.Response
	Data []DungeonCharacterData `json:"data"`
}

// CharacterRequest -
type CharacterRequest struct {
	schema.Request
	Data DungeonCharacterData `json:"data"`
}
