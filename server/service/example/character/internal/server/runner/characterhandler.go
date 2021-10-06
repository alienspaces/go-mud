package runner

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/modeller"
	"gitlab.com/alienspaces/go-boilerplate/server/schema"
	"gitlab.com/alienspaces/go-boilerplate/server/service/character/internal/model"
	"gitlab.com/alienspaces/go-boilerplate/server/service/character/internal/record"
)

// GetCharactersHandler -
// Admininstrator or default role
// No restrictions on character types
func (rnr *Runner) GetCharactersHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) {

	l.Info("** Get characters handler ** p >%#v< m >%#v<", pp, m)

	var characterRecs []*record.Character

	var err error

	id := pp.ByName("character_id")

	// single resource
	if id != "" {

		l.Info("Getting character record ID >%s<", id)

		rec, err := m.(*model.Model).GetCharacterRec(id, false)
		if err != nil {
			rnr.WriteModelError(l, w, err)
			return
		}

		// resource not found
		if rec == nil {
			rnr.WriteNotFoundError(l, w, id)
			return
		}

		characterRecs = append(characterRecs, rec)

	} else {

		l.Info("Querying character records")

		// query parameters
		params := make(map[string]interface{})
		for paramName, paramValue := range qp {
			params[paramName] = paramValue
		}

		characterRecs, err = m.(*model.Model).GetCharacterRecs(params, nil, false)
		if err != nil {
			rnr.WriteModelError(l, w, err)
			return
		}
	}

	// assign response properties
	data := []schema.CharacterData{}
	for _, characterRec := range characterRecs {

		// response data
		responseData, err := rnr.RecordToCharacterResponseData(characterRec)
		if err != nil {
			rnr.WriteSystemError(l, w, err)
			return
		}

		data = append(data, responseData)
	}

	res := schema.CharacterResponse{
		Data: data,
	}

	err = rnr.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return
	}
}

// GetPlayerCharactersHandler -
// - Administrator or default role
// - Player ID required
// - Player ID in path must match idcharacter
// - Restricted to character types `player-mage` and `player-familliar`
func (rnr *Runner) GetPlayerCharactersHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) {

	l.Info("** Get characters handler ** p >%#v< m >%#v<", pp, m)

	var characterRecs []*record.Character

	var err error

	playerID := pp.ByName("player_id")
	characterID := pp.ByName("character_id")

	// single resource
	if playerID != "" && characterID != "" {

		l.Info("Getting player ID >%s< character ID >%s< record ", playerID, characterID)

		characterRec, err := m.(*model.Model).GetCharacterRec(characterID, false)
		if err != nil {
			rnr.WriteModelError(l, w, err)
			return
		}

		// Character records
		characterRecs = append(characterRecs, characterRec)

	} else {

		l.Info("Querying character records")

		// query parameters
		params := make(map[string]interface{})
		for paramName, paramValue := range qp {
			params[paramName] = paramValue
		}

		characterRecs, err = m.(*model.Model).GetCharacterRecs(params, nil, false)
		if err != nil {
			rnr.WriteModelError(l, w, err)
			return
		}

	}

	// assign response properties
	data := []schema.CharacterData{}
	for _, characterRec := range characterRecs {

		// response data
		responseData, err := rnr.RecordToCharacterResponseData(characterRec)
		if err != nil {
			rnr.WriteSystemError(l, w, err)
			return
		}

		data = append(data, responseData)
	}

	res := schema.CharacterResponse{
		Data: data,
	}

	err = rnr.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return
	}
}

// PostPlayerCharactersHandler -
func (rnr *Runner) PostPlayerCharactersHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) {

	l.Info("** Post characters handler ** p >%#v< m >#%v<", pp, m)

	// Path parameters
	characterID := pp.ByName("character_id")
	playerID := pp.ByName("player_id")

	req := schema.CharacterRequest{}

	err := rnr.ReadRequest(l, r, &req)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	characterRec := record.Character{}

	// Assign request properties
	characterRec.ID = characterID
	characterRec.PlayerID = playerID

	// Record data
	err = rnr.CharacterRequestDataToRecord(req.Data, &characterRec)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	err = m.(*model.Model).CreateCharacterRec(&characterRec)
	if err != nil {
		rnr.WriteModelError(l, w, err)
		return
	}

	// Response data
	responseData, err := rnr.RecordToCharacterResponseData(&characterRec)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	// Assign response properties
	res := schema.CharacterResponse{
		Data: []schema.CharacterData{
			responseData,
		},
	}

	err = rnr.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return
	}
}

// PutPlayerCharactersHandler -
func (rnr *Runner) PutPlayerCharactersHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) {

	l.Info("** Put characters handler ** p >%#v< m >#%v<", pp, m)

	characterID := pp.ByName("character_id")
	playerID := pp.ByName("player_id")

	l.Info("Updating resource player ID >%s< character ID >%s<", playerID, characterID)

	// Player character record
	l.Info("Getting player ID >%s< character ID >%s< record ", playerID, characterID)

	characterRec, err := m.(*model.Model).GetCharacterRec(characterID, false)
	if err != nil {
		rnr.WriteModelError(l, w, err)
		return
	}

	// Resource not found
	if characterRec == nil {
		rnr.WriteNotFoundError(l, w, characterID)
		return
	}

	req := schema.CharacterRequest{}

	err = rnr.ReadRequest(l, r, &req)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	// Record data
	err = rnr.CharacterRequestDataToRecord(req.Data, characterRec)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	err = m.(*model.Model).UpdateCharacterRec(characterRec)
	if err != nil {
		rnr.WriteModelError(l, w, err)
		return
	}

	// Response data
	responseData, err := rnr.RecordToCharacterResponseData(characterRec)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	// Assign response properties
	res := schema.CharacterResponse{
		Data: []schema.CharacterData{
			responseData,
		},
	}

	err = rnr.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return
	}
}

// CharacterRequestDataToRecord -
func (rnr *Runner) CharacterRequestDataToRecord(data schema.CharacterData, rec *record.Character) error {

	// NOTE: PlayerID is sourced from path parameters

	rec.Name = data.Name
	rec.Avatar = data.Avatar
	rec.Strength = data.Strength
	rec.Dexterity = data.Dexterity
	rec.Intelligence = data.Intelligence

	return nil
}

// RecordToCharacterResponseData -
func (rnr *Runner) RecordToCharacterResponseData(characterRec *record.Character) (schema.CharacterData, error) {

	data := schema.CharacterData{
		ID:               characterRec.ID,
		PlayerID:         characterRec.PlayerID,
		Name:             characterRec.Name,
		Avatar:           characterRec.Avatar,
		Strength:         characterRec.Strength,
		Dexterity:        characterRec.Dexterity,
		Intelligence:     characterRec.Intelligence,
		AttributePoints:  characterRec.AttributePoints,
		ExperiencePoints: characterRec.ExperiencePoints,
		Coins:            characterRec.Coins,
		CreatedAt:        characterRec.CreatedAt,
		UpdatedAt:        characterRec.UpdatedAt.Time,
	}

	return data, nil
}
