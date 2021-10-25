package runner

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/server/schema"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetCharacterHandler -
func (rnr *Runner) GetDungeonCharacterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) {

	l.Info("** Get dungeons handler ** p >%#v< m >%#v<", pp, m)

	var recs []*record.DungeonCharacter
	var err error

	// Path parameters
	dungeonID := pp.ByName("dungeon_id")
	characterID := pp.ByName("character_id")

	if dungeonID == "" {
		rnr.WriteNotFoundError(l, w, dungeonID)
		return
	}

	if characterID == "" {
		rnr.WriteNotFoundError(l, w, characterID)
		return
	}

	l.Info("Getting character record ID >%s<", characterID)

	rec, err := m.(*model.Model).GetDungeonCharacterRec(characterID, false)
	if err != nil {
		rnr.WriteModelError(l, w, err)
		return
	}

	// Resource not found
	if rec == nil {
		rnr.WriteNotFoundError(l, w, characterID)
		return
	}

	// Character record dungeon does not match parameter dungeon
	if rec.DungeonID != dungeonID {
		rnr.WriteNotFoundError(l, w, characterID)
		return
	}

	recs = append(recs, rec)

	// Assign response properties
	data := []schema.DungeonCharacterData{}
	for _, rec := range recs {

		// Response data
		responseData, err := rnr.RecordToDungeonCharacterResponseData(rec)
		if err != nil {
			rnr.WriteSystemError(l, w, err)
			return
		}

		data = append(data, responseData)
	}

	res := schema.DungeonCharacterResponse{
		Data: data,
	}

	err = rnr.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return
	}
}

// GetCharactersHandler -
func (rnr *Runner) GetDungeonCharactersHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) {

	l.Info("** Get dungeons handler ** p >%#v< m >%#v<", pp, m)

	var recs []*record.DungeonCharacter
	var err error

	// Path parameters
	dungeonID := pp.ByName("dungeon_id")

	if dungeonID == "" {
		rnr.WriteNotFoundError(l, w, dungeonID)
		return
	}

	l.Info("Querying dungeon records")

	// Add query parameters
	params := make(map[string]interface{})
	for paramName, paramValue := range qp {
		l.Info("Querying dungeon records with param name >%s< value >%v<", paramName, paramValue)
		params[paramName] = paramValue
	}

	// Add path parameters
	params["dungeon_id"] = dungeonID

	recs, err = m.(*model.Model).GetDungeonCharacterRecs(params, nil, false)
	if err != nil {
		rnr.WriteModelError(l, w, err)
		return
	}

	// Assign response properties
	data := []schema.DungeonCharacterData{}
	for _, rec := range recs {

		// Response data
		responseData, err := rnr.RecordToDungeonCharacterResponseData(rec)
		if err != nil {
			rnr.WriteSystemError(l, w, err)
			return
		}

		data = append(data, responseData)
	}

	res := schema.DungeonCharacterResponse{
		Data: data,
	}

	l.Info("Responding with >%#v<", res)

	err = rnr.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return
	}
}

// PostDungeonCharactersHandler -
func (rnr *Runner) PostDungeonCharactersHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) {

	l.Info("** Post characters handler ** p >%#v< m >#%v<", pp, m)

	// Path parameters
	dungeonID := pp.ByName("dungeon_id")

	req := schema.DungeonCharacterRequest{}

	err := rnr.ReadRequest(l, r, &req)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	rec := record.DungeonCharacter{DungeonID: dungeonID}

	// Record data
	err = rnr.DungeonCharacterRequestDataToRecord(req.Data, &rec)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	err = m.(*model.Model).CreateDungeonCharacterRec(&rec)
	if err != nil {
		rnr.WriteModelError(l, w, err)
		return
	}

	// Response data
	responseData, err := rnr.RecordToDungeonCharacterResponseData(&rec)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	// Assign response properties
	res := schema.DungeonCharacterResponse{
		Data: []schema.DungeonCharacterData{
			responseData,
		},
	}

	err = rnr.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return
	}
}

// PutDungeonCharactersHandler -
func (rnr *Runner) PutDungeonCharacterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) {

	l.Info("** Put characters handler ** p >%#v< m >#%v<", pp, m)

	// Path parameters
	dungeonID := pp.ByName("dungeon_id")
	characterID := pp.ByName("character_id")

	l.Info("Updating character ID >%s<", characterID)

	rec, err := m.(*model.Model).GetDungeonCharacterRec(characterID, false)
	if err != nil {
		rnr.WriteModelError(l, w, err)
		return
	}

	// Resource not found
	if rec == nil {
		rnr.WriteNotFoundError(l, w, characterID)
		return
	}

	// Character record dungeon does not match parameter dungeon
	if rec.DungeonID != dungeonID {
		rnr.WriteNotFoundError(l, w, characterID)
		return
	}

	req := schema.DungeonCharacterRequest{}

	err = rnr.ReadRequest(l, r, &req)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	// Record data
	err = rnr.DungeonCharacterRequestDataToRecord(req.Data, rec)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	err = m.(*model.Model).UpdateDungeonCharacterRec(rec)
	if err != nil {
		rnr.WriteModelError(l, w, err)
		return
	}

	// Response data
	responseData, err := rnr.RecordToDungeonCharacterResponseData(rec)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	// Assign response properties
	res := schema.DungeonCharacterResponse{
		Data: []schema.DungeonCharacterData{
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
func (rnr *Runner) DungeonCharacterRequestDataToRecord(data schema.DungeonCharacterData, rec *record.DungeonCharacter) error {

	rec.Name = data.Name
	rec.Strength = data.Strength
	rec.Dexterity = data.Dexterity
	rec.Intelligence = data.Intelligence

	return nil
}

// RecordToCharacterResponseData -
func (rnr *Runner) RecordToDungeonCharacterResponseData(dungeonCharacterRec *record.DungeonCharacter) (schema.DungeonCharacterData, error) {

	data := schema.DungeonCharacterData{
		ID:        dungeonCharacterRec.ID,
		Name:      dungeonCharacterRec.Name,
		CreatedAt: dungeonCharacterRec.CreatedAt,
		UpdatedAt: dungeonCharacterRec.UpdatedAt.Time,
	}

	return data, nil
}
