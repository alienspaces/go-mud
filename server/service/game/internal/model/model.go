package model

import (
	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-mud/server/core/model"
	"gitlab.com/alienspaces/go-mud/server/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/preparer"
	"gitlab.com/alienspaces/go-mud/server/core/type/repositor"
	"gitlab.com/alienspaces/go-mud/server/core/type/storer"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/dungeon"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/dungeonlocation"
)

// Model -
type Model struct {
	model.Model
}

// NewModel -
func NewModel(c configurer.Configurer, l logger.Logger, s storer.Storer) (*Model, error) {

	m := &Model{
		model.Model{
			Config: c,
			Log:    l,
			Store:  s,
		},
	}

	m.RepositoriesFunc = m.NewRepositories

	return m, nil
}

// NewRepositories - Custom repositories for this model
func (m *Model) NewRepositories(p preparer.Preparer, tx *sqlx.Tx) ([]repositor.Repositor, error) {

	repositoryList := []repositor.Repositor{}

	tr, err := dungeon.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new game repository >%v<", err)
		return nil, err
	}

	repositoryList = append(repositoryList, tr)

	return repositoryList, nil
}

// DungeonRepository -
func (m *Model) DungeonRepository() *dungeon.Repository {

	r := m.Repositories[dungeon.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", dungeon.TableName)
		return nil
	}

	return r.(*dungeon.Repository)
}

// DungeonLocationRepository -
func (m *Model) DungeonLocationRepository() *dungeonlocation.Repository {

	r := m.Repositories[dungeonlocation.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", dungeonlocation.TableName)
		return nil
	}

	return r.(*dungeonlocation.Repository)
}
