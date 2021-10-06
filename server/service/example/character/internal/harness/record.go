package harness

import (
	"github.com/brianvoe/gofakeit/v5"
	"github.com/google/uuid"

	"gitlab.com/alienspaces/go-boilerplate/server/service/character/internal/model"
	"gitlab.com/alienspaces/go-boilerplate/server/service/character/internal/record"
)

func (t *Testing) createCharacterRec(chracterConfig CharacterConfig) (record.Character, error) {

	chracterRec := chracterConfig.Record

	if chracterRec.PlayerID == "" {
		playerID, _ := uuid.NewRandom()
		chracterRec.PlayerID = playerID.String()
	}

	if chracterRec.Name == "" {
		chracterRec.Name = gofakeit.Name()
	}

	t.Log.Info("Creating chracter testing record >%#v<", chracterRec)

	err := t.Model.(*model.Model).CreateCharacterRec(&chracterRec)
	if err != nil {
		t.Log.Warn("Failed creating testing chracter record >%v<", err)
		return chracterRec, err
	}

	return chracterRec, nil
}
