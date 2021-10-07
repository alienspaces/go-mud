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

// GetGamesHandler -
func (rnr *Runner) GetGamesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) {

	l.Info("** Get games handler ** p >%#v< m >%#v<", pp, m)

	var recs []*record.Game
	var err error

	// Path parameters
	id := pp.ByName("game_id")

	// Single resource
	if id != "" {

		l.Info("Getting game record ID >%s<", id)

		rec, err := m.(*model.Model).GetGameRec(id, false)
		if err != nil {
			rnr.WriteModelError(l, w, err)
			return
		}

		// Resource not found
		if rec == nil {
			rnr.WriteNotFoundError(l, w, id)
			return
		}

		recs = append(recs, rec)

	} else {

		l.Info("Querying game records")

		params := make(map[string]interface{})
		for paramName, paramValue := range qp {
			l.Info("Querying game records with param name >%s< value >%v<", paramName, paramValue)
			params[paramName] = paramValue
		}

		recs, err = m.(*model.Model).GetGameRecs(params, nil, false)
		if err != nil {
			rnr.WriteModelError(l, w, err)
			return
		}
	}

	// Assign response properties
	data := []schema.GameData{}
	for _, rec := range recs {

		// Response data
		responseData, err := rnr.RecordToGameResponseData(rec)
		if err != nil {
			rnr.WriteSystemError(l, w, err)
			return
		}

		data = append(data, responseData)
	}

	res := schema.GameResponse{
		Data: data,
	}

	err = rnr.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return
	}
}

// PostGamesHandler -
func (rnr *Runner) PostGamesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) {

	l.Info("** Post games handler ** p >%#v< m >#%v<", pp, m)

	// Path parameters
	id := pp.ByName("game_id")

	req := schema.GameRequest{}

	err := rnr.ReadRequest(l, r, &req)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	rec := record.Game{}

	// Assign request properties
	rec.ID = id

	// Record data
	err = rnr.GameRequestDataToRecord(req.Data, &rec)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	err = m.(*model.Model).CreateGameRec(&rec)
	if err != nil {
		rnr.WriteModelError(l, w, err)
		return
	}

	// Response data
	responseData, err := rnr.RecordToGameResponseData(&rec)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	// Assign response properties
	res := schema.GameResponse{
		Data: []schema.GameData{
			responseData,
		},
	}

	err = rnr.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return
	}
}

// PutGamesHandler -
func (rnr *Runner) PutGamesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) {

	l.Info("** Put games handler ** p >%#v< m >#%v<", pp, m)

	// Path parameters
	id := pp.ByName("game_id")

	l.Info("Updating resource ID >%s<", id)

	rec, err := m.(*model.Model).GetGameRec(id, false)
	if err != nil {
		rnr.WriteModelError(l, w, err)
		return
	}

	// Resource not found
	if rec == nil {
		rnr.WriteNotFoundError(l, w, id)
		return
	}

	req := schema.GameRequest{}

	err = rnr.ReadRequest(l, r, &req)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	// Record data
	err = rnr.GameRequestDataToRecord(req.Data, rec)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	err = m.(*model.Model).UpdateGameRec(rec)
	if err != nil {
		rnr.WriteModelError(l, w, err)
		return
	}

	// Response data
	responseData, err := rnr.RecordToGameResponseData(rec)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	// Assign response properties
	res := schema.GameResponse{
		Data: []schema.GameData{
			responseData,
		},
	}

	err = rnr.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return
	}
}

// GameRequestDataToRecord -
func (rnr *Runner) GameRequestDataToRecord(data schema.GameData, rec *record.Game) error {

	return nil
}

// RecordToGameResponseData -
func (rnr *Runner) RecordToGameResponseData(gameRec *record.Game) (schema.GameData, error) {

	data := schema.GameData{
		ID:        gameRec.ID,
		CreatedAt: gameRec.CreatedAt,
		UpdatedAt: gameRec.UpdatedAt.Time,
	}

	return data, nil
}
