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

// GetDungeonHandler -
func (rnr *Runner) GetDungeonHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) {

	l.Info("** Get games handler ** p >%#v< m >%#v<", pp, m)

	var recs []*record.Game
	var err error

	// Path parameters
	id := pp.ByName("dungeon_id")

	if id == "" {
		rnr.WriteNotFoundError(l, w, id)
		return
	}

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

	// Assign response properties
	data := []schema.DungeonData{}
	for _, rec := range recs {

		// Response data
		responseData, err := rnr.RecordToDungeonResponseData(rec)
		if err != nil {
			rnr.WriteSystemError(l, w, err)
			return
		}

		data = append(data, responseData)
	}

	res := schema.DungeonResponse{
		Data: data,
	}

	err = rnr.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return
	}
}

// GetDungeonsHandler -
func (rnr *Runner) GetDungeonsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) {

	l.Info("** Get games handler ** p >%#v< m >%#v<", pp, m)

	var recs []*record.Game
	var err error

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

	// Assign response properties
	data := []schema.DungeonData{}
	for _, rec := range recs {

		// Response data
		responseData, err := rnr.RecordToDungeonResponseData(rec)
		if err != nil {
			rnr.WriteSystemError(l, w, err)
			return
		}

		data = append(data, responseData)
	}

	res := schema.DungeonResponse{
		Data: data,
	}

	err = rnr.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return
	}
}

// DungeonRequestDataToRecord -
func (rnr *Runner) DungeonRequestDataToRecord(data schema.DungeonData, rec *record.Dungeon) error {

	return nil
}

// RecordToDungeonResponseData -
func (rnr *Runner) RecordToDungeonResponseData(gameRec *record.Dungeon) (schema.DungeonData, error) {

	data := schema.DungeonData{
		ID:        gameRec.ID,
		CreatedAt: gameRec.CreatedAt,
		UpdatedAt: gameRec.UpdatedAt.Time,
	}

	return data, nil
}
