package runner

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/modeller"
	"gitlab.com/alienspaces/go-boilerplate/server/schema"
	"gitlab.com/alienspaces/go-boilerplate/server/service/player/internal/model"
	"gitlab.com/alienspaces/go-boilerplate/server/service/player/internal/record"
)

// GetPlayersHandler -
func (rnr *Runner) GetPlayersHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) {

	l.Info("** Get players handler ** p >%#v< m >%#v<", pp, m)

	var recs []*record.Player
	var err error

	id := pp.ByName("player_id")

	// single resource
	if id != "" {

		l.Info("Getting player record ID >%s<", id)

		rec, err := m.(*model.Model).GetPlayerRec(id, false)
		if err != nil {
			rnr.WriteModelError(l, w, err)
			return
		}

		// resource not found
		if rec == nil {
			rnr.WriteNotFoundError(l, w, id)
			return
		}

		recs = append(recs, rec)

	} else {

		l.Info("Querying player records")

		params := make(map[string]interface{})
		for paramName, paramValue := range qp {
			params[paramName] = paramValue
		}

		recs, err = m.(*model.Model).GetPlayerRecs(params, nil, false)
		if err != nil {
			rnr.WriteModelError(l, w, err)
			return
		}
	}

	// assign response properties
	data := []schema.PlayerData{}
	for _, rec := range recs {

		// response data
		responseData, err := rnr.RecordToPlayerResponseData(rec)
		if err != nil {
			rnr.WriteSystemError(l, w, err)
			return
		}

		data = append(data, responseData)
	}

	res := schema.PlayerResponse{
		Data: data,
	}

	err = rnr.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return
	}
}

// PostPlayersHandler -
func (rnr *Runner) PostPlayersHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) {

	l.Info("** Post players handler ** p >%#v< m >#%v<", pp, m)

	// parameters
	id := pp.ByName("player_id")

	req := schema.PlayerRequest{}

	err := rnr.ReadRequest(l, r, &req)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	rec := record.Player{}

	// assign request properties
	rec.ID = id

	// record data
	err = rnr.PlayerRequestDataToRecord(req.Data, &rec)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	err = m.(*model.Model).CreatePlayerRec(&rec)
	if err != nil {
		rnr.WriteModelError(l, w, err)
		return
	}

	// response data
	responseData, err := rnr.RecordToPlayerResponseData(&rec)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	// assign response properties
	res := schema.PlayerResponse{
		Data: []schema.PlayerData{
			responseData,
		},
	}

	err = rnr.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return
	}
}

// PutPlayersHandler -
func (rnr *Runner) PutPlayersHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) {

	l.Info("** Put players handler ** p >%#v< m >#%v<", pp, m)

	// parameters
	id := pp.ByName("player_id")

	l.Info("Updating resource ID >%s<", id)

	rec, err := m.(*model.Model).GetPlayerRec(id, false)
	if err != nil {
		rnr.WriteModelError(l, w, err)
		return
	}

	// resource not found
	if rec == nil {
		rnr.WriteNotFoundError(l, w, id)
		return
	}

	req := schema.PlayerRequest{}

	err = rnr.ReadRequest(l, r, &req)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	// record data
	err = rnr.PlayerRequestDataToRecord(req.Data, rec)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	err = m.(*model.Model).UpdatePlayerRec(rec)
	if err != nil {
		rnr.WriteModelError(l, w, err)
		return
	}

	// response data
	responseData, err := rnr.RecordToPlayerResponseData(rec)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	// assign response properties
	res := schema.PlayerResponse{
		Data: []schema.PlayerData{
			responseData,
		},
	}

	err = rnr.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return
	}
}

// PlayerRequestDataToRecord -
func (rnr *Runner) PlayerRequestDataToRecord(data schema.PlayerData, rec *record.Player) error {

	rec.Name = data.Name
	rec.Email = data.Email
	rec.Provider = data.Provider
	rec.ProviderAccountID = data.ProviderAccountID

	return nil
}

// RecordToPlayerResponseData -
func (rnr *Runner) RecordToPlayerResponseData(playerRec *record.Player) (schema.PlayerData, error) {

	data := schema.PlayerData{
		ID:                playerRec.ID,
		Name:              playerRec.Name,
		Email:             playerRec.Email,
		Provider:          playerRec.Provider,
		ProviderAccountID: playerRec.ProviderAccountID,
		CreatedAt:         playerRec.CreatedAt,
		UpdatedAt:         playerRec.UpdatedAt.Time,
	}

	return data, nil
}
