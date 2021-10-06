package harness

import (
	"gitlab.com/alienspaces/go-boilerplate/server/core/harness"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/modeller"
	"gitlab.com/alienspaces/go-boilerplate/server/service/character/internal/model"
	"gitlab.com/alienspaces/go-boilerplate/server/service/character/internal/record"
)

// Testing -
type Testing struct {
	harness.Testing
	Data       *Data
	DataConfig DataConfig
}

// DataConfig -
type DataConfig struct {
	CharacterConfig []CharacterConfig
}

// CharacterConfig -
type CharacterConfig struct {
	Record record.Character
}

// Data -
type Data struct {
	CharacterRecs []record.Character
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

	// Stand alone character records
	for _, characterConfig := range t.DataConfig.CharacterConfig {

		characterRec, err := t.createCharacterRec(characterConfig)
		if err != nil {
			t.Log.Warn("Failed creating character record >%v<", err)
			return err
		}
		t.Data.CharacterRecs = append(t.Data.CharacterRecs, characterRec)
	}

	return nil
}

// RemoveData -
func (t *Testing) RemoveData() error {

CHARACTER_RECS:
	for {
		if len(t.Data.CharacterRecs) == 0 {
			break CHARACTER_RECS
		}
		rec := record.Character{}
		rec, t.Data.CharacterRecs = t.Data.CharacterRecs[0], t.Data.CharacterRecs[1:]

		err := t.Model.(*model.Model).RemoveCharacterRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing character record >%v<", err)
			return err
		}
	}

	return nil
}
