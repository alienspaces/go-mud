package harness

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

func (t *Testing) createGameRec(gameConfig GameConfig) (record.Game, error) {

	rec := gameConfig.Record

	// NOTE: Add default values for required properties here

	t.Log.Info("Creating testing record >%#v<", rec)

	err := t.Model.(*model.Model).CreateGameRec(&rec)
	if err != nil {
		t.Log.Warn("Failed creating testing game record >%v<", err)
		return rec, err
	}
	return rec, nil
}
