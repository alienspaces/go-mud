package harness

import (
	"gitlab.com/alienspaces/go-boilerplate/server/core/harness"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/modeller"
	"gitlab.com/alienspaces/go-boilerplate/server/service/player/internal/model"
	"gitlab.com/alienspaces/go-boilerplate/server/service/player/internal/record"
)

// Testing -
type Testing struct {
	harness.Testing
	Data       *Data
	DataConfig DataConfig
}

// DataConfig -
type DataConfig struct {
	PlayerConfig []PlayerConfig
}

// PlayerConfig -
type PlayerConfig struct {
	Record            record.Player
	PlayerRoleConfig []PlayerRoleConfig
}

// PlayerRoleConfig -
type PlayerRoleConfig struct {
	Record record.PlayerRole
}

// Data -
type Data struct {
	PlayerRecs     []record.Player
	PlayerRoleRecs []record.PlayerRole
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

	for _, accountConfig := range t.DataConfig.PlayerConfig {

		accountRec, err := t.createPlayerRec(accountConfig)
		if err != nil {
			t.Log.Warn("Failed creating account record >%v<", err)
			return err
		}
		t.Data.PlayerRecs = append(t.Data.PlayerRecs, accountRec)

		for _, accountRoleConfig := range accountConfig.PlayerRoleConfig {
			accountRoleRec, err := t.createPlayerRoleRec(accountRec, accountRoleConfig)
			if err != nil {
				t.Log.Warn("Failed creating account role record >%v<", err)
				return err
			}
			t.Data.PlayerRoleRecs = append(t.Data.PlayerRoleRecs, accountRoleRec)
		}
	}

	return nil
}

// RemoveData -
func (t *Testing) RemoveData() error {

ACCOUNT_ROLE_RECS:
	for {
		if len(t.Data.PlayerRoleRecs) == 0 {
			break ACCOUNT_ROLE_RECS
		}
		rec := record.PlayerRole{}
		rec, t.Data.PlayerRoleRecs = t.Data.PlayerRoleRecs[0], t.Data.PlayerRoleRecs[1:]

		err := t.Model.(*model.Model).RemovePlayerRoleRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing account role record >%v<", err)
			return err
		}
	}

ACCOUNT_RECS:
	for {
		if len(t.Data.PlayerRecs) == 0 {
			break ACCOUNT_RECS
		}
		rec := record.Player{}
		rec, t.Data.PlayerRecs = t.Data.PlayerRecs[0], t.Data.PlayerRecs[1:]

		err := t.Model.(*model.Model).RemovePlayerRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing account record >%v<", err)
			return err
		}
	}

	return nil
}
