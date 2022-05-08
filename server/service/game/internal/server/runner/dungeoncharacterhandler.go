package runner

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	coreerror "gitlab.com/alienspaces/go-mud/server/core/error"
	"gitlab.com/alienspaces/go-mud/server/core/server"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/server/schema"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetCharacterHandler -
func (rnr *Runner) GetCharacterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {

	l.Info("** Get dungeons handler ** p >%#v< m >%#v<", pp, m)

	var recs []*record.Character
	var err error

	// Path parameters
	characterID := pp.ByName("character_id")

	if characterID == "" {
		err := coreerror.NewNotFoundError("character", characterID)
		server.WriteError(l, w, err)
		return err
	}
	if !m.(*model.Model).IsUUID(characterID) {
		err := coreerror.NewNotFoundError("character", characterID)
		server.WriteError(l, w, err)
		return err
	}

	l.Info("Getting character record ID >%s<", characterID)

	rec, err := m.(*model.Model).GetCharacterRec(characterID, false)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	// Resource not found
	if rec == nil {
		err := coreerror.NewNotFoundError("character", characterID)
		server.WriteError(l, w, err)
		return err
	}

	recs = append(recs, rec)

	// Assign response properties
	data := []schema.CharacterData{}
	for _, rec := range recs {

		// Response data
		responseData, err := rnr.RecordToCharacterResponseData(*rec)
		if err != nil {
			server.WriteError(l, w, err)
			return err
		}

		data = append(data, responseData)
	}

	res := schema.CharacterResponse{
		Data: data,
	}

	err = server.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return err
	}

	return nil
}

// GetCharactersHandler -
func (rnr *Runner) GetCharactersHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {

	l.Info("** Get characters handler ** p >%#v< m >%#v<", pp, m)

	var recs []*record.Character
	var err error

	l.Info("Querying character records")

	// Add query parameters
	params := make(map[string]interface{})
	for paramName, paramValue := range qp {
		l.Info("Querying dungeon records with param name >%s< value >%v<", paramName, paramValue)
		params[paramName] = paramValue
	}

	recs, err = m.(*model.Model).GetCharacterRecs(params, nil, false)
	if err != nil {
		l.Warn("Failed getting dungeon character records >%v<", err)
		server.WriteError(l, w, err)
		return err
	}

	// Assign response properties
	data := []schema.CharacterData{}
	for _, rec := range recs {

		// Response data
		responseData, err := rnr.RecordToCharacterResponseData(*rec)
		if err != nil {
			server.WriteError(l, w, err)
			return err
		}

		data = append(data, responseData)
	}

	res := schema.CharacterResponse{
		Data: data,
	}

	l.Info("Responding with >%#v<", res)

	err = server.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return err
	}

	return nil
}

// PostCharactersHandler -
func (rnr *Runner) PostCharactersHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {

	l.Info("** Post characters handler ** p >%#v< m >#%v<", pp, m)

	req := schema.CharacterRequest{}
	err := server.ReadRequest(l, r, &req)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	rec := record.Character{}

	// Record data
	err = rnr.CharacterRequestDataToRecord(req.Data, &rec)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	l.Info("Creating character record >%#v<", rec)

	err = m.(*model.Model).CreateCharacterRec(&rec)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	// Response data
	responseData, err := rnr.RecordToCharacterResponseData(rec)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	// Assign response properties
	res := schema.CharacterResponse{
		Data: []schema.CharacterData{
			responseData,
		},
	}

	err = server.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return err
	}

	return nil
}

// PutCharactersHandler -
func (rnr *Runner) PutCharacterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {

	l.Info("** Put characters handler ** p >%#v< m >#%v<", pp, m)

	// Path parameters
	id := pp.ByName("character_id")

	l.Info("Updating character ID >%s<", id)

	rec, err := m.(*model.Model).GetCharacterRec(id, false)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	// Resource not found
	if rec == nil {
		err := coreerror.NewNotFoundError("character", id)
		server.WriteError(l, w, err)
		return err
	}

	req := schema.CharacterRequest{}

	err = server.ReadRequest(l, r, &req)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	// Record data
	err = rnr.CharacterRequestDataToRecord(req.Data, rec)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	err = m.(*model.Model).UpdateCharacterRec(rec)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	// Response data
	responseData, err := rnr.RecordToCharacterResponseData(*rec)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	// Assign response properties
	res := schema.CharacterResponse{
		Data: []schema.CharacterData{
			responseData,
		},
	}

	err = server.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return err
	}

	return err
}

// CharacterRequestDataToRecord -
func (rnr *Runner) CharacterRequestDataToRecord(data schema.CharacterData, rec *record.Character) error {

	rec.Name = data.Name
	rec.Strength = data.Strength
	rec.Dexterity = data.Dexterity
	rec.Intelligence = data.Intelligence

	return nil
}

// RecordToCharacterResponseData -
func (rnr *Runner) RecordToCharacterResponseData(dungeonCharacterRec record.Character) (schema.CharacterData, error) {

	data := schema.CharacterData{
		ID:               dungeonCharacterRec.ID,
		Name:             dungeonCharacterRec.Name,
		Strength:         dungeonCharacterRec.Strength,
		Dexterity:        dungeonCharacterRec.Dexterity,
		Intelligence:     dungeonCharacterRec.Intelligence,
		Health:           dungeonCharacterRec.Health,
		Fatigue:          dungeonCharacterRec.Fatigue,
		Coins:            dungeonCharacterRec.Coins,
		AttributePoints:  dungeonCharacterRec.AttributePoints,
		ExperiencePoints: dungeonCharacterRec.ExperiencePoints,
		CreatedAt:        dungeonCharacterRec.CreatedAt,
		UpdatedAt:        dungeonCharacterRec.UpdatedAt.Time,
	}

	return data, nil
}
