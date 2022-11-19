package schema

// CharacterResponse -
type CharacterResponse struct {
	Response
	Data []DungeonCharacterData `json:"data"`
}

// CharacterRequest -
type CharacterRequest struct {
	Request
	Data DungeonCharacterData `json:"data"`
}
