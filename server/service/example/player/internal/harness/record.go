package harness

import (
	"github.com/brianvoe/gofakeit"

	"gitlab.com/alienspaces/go-boilerplate/server/service/player/internal/model"
	"gitlab.com/alienspaces/go-boilerplate/server/service/player/internal/record"
)

func (t *Testing) createPlayerRec(accountConfig PlayerConfig) (record.Player, error) {

	rec := accountConfig.Record

	t.Log.Info("Creating test account record >%#v<", rec)

	// NOTE: Add default values for required properties here
	if rec.Name == "" {
		rec.Name = gofakeit.Name()
	}

	if rec.Email == "" {
		rec.Email = gofakeit.Email()
	}

	if rec.Provider == "" {
		rec.Provider = record.AccountProviderGoogle
	}

	if rec.ProviderAccountID == "" {
		rec.ProviderAccountID = gofakeit.UUID()
	}

	err := t.Model.(*model.Model).CreatePlayerRec(&rec)
	if err != nil {
		t.Log.Warn("Failed creating testing account record >%v<", err)
		return rec, err
	}
	return rec, nil
}

func (t *Testing) createPlayerRoleRec(accountRec record.Player, accountRoleConfig PlayerRoleConfig) (record.PlayerRole, error) {

	rec := accountRoleConfig.Record

	t.Log.Info("Creating test account role record >%#v<", rec)

	// NOTE: Add default values for required properties here
	rec.PlayerID = accountRec.ID

	if rec.Role == "" {
		rec.Role = record.PlayerRoleDefault
	}

	err := t.Model.(*model.Model).CreatePlayerRoleRec(&rec)
	if err != nil {
		t.Log.Warn("Failed creating testing account role record >%v<", err)
		return rec, err
	}
	return rec, nil
}
