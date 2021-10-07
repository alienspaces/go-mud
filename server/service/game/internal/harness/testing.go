package harness

import (
	"gitlab.com/alienspaces/go-mud/server/core/harness"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// Testing -
type Testing struct {
	harness.Testing
	Data       *Data
	DataConfig DataConfig
}

// DataConfig -
type DataConfig struct {
	GameConfig []GameConfig
}

// GameConfig -
type GameConfig struct {
	Record record.Game
}

// Data -
type Data struct {
	GameRecs []record.Game
}

// NewTesting -
func NewTesting(config DataConfig) (t *Testing, err error) {

	// harness
	t = &Testing{}

	// modeller
	t.ModellerFunc = t.Modeller

	// data
	t.CreateDataFunc = t.CreateData
	t.RemoveDataFunc = t.RemoveData

	t.DataConfig = config
	t.Data = &Data{}

	err = t.Init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Modeller -
func (t *Testing) Modeller() (modeller.Modeller, error) {

	m, err := model.NewModel(t.Config, t.Log, t.Store)
	if err != nil {
		t.Log.Warn("Failed new model >%v<", err)
		return nil, err
	}

	return m, nil
}

// CreateData - Custom data
func (t *Testing) CreateData() error {

	for _, gameConfig := range t.DataConfig.GameConfig {

		gameRec, err := t.createGameRec(gameConfig)
		if err != nil {
			t.Log.Warn("Failed creating game record >%v<", err)
			return err
		}
		t.Data.GameRecs = append(t.Data.GameRecs, gameRec)
	}

	return nil
}

// RemoveData -
func (t *Testing) RemoveData() error {

GAME_RECS:
	for {
		if len(t.Data.GameRecs) == 0 {
			break GAME_RECS
		}
		rec := record.Game{}
		rec, t.Data.GameRecs = t.Data.GameRecs[0], t.Data.GameRecs[1:]

		err := t.Model.(*model.Model).RemoveGameRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing game record >%v<", err)
			return err
		}
	}

	return nil
}
