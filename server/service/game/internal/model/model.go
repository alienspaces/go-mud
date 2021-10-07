package model

import (
	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-mud/server/core/model"
	"gitlab.com/alienspaces/go-mud/server/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/preparer"
	"gitlab.com/alienspaces/go-mud/server/core/type/repositor"
	"gitlab.com/alienspaces/go-mud/server/core/type/storer"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/game"
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

	tr, err := game.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new game repository >%v<", err)
		return nil, err
	}

	repositoryList = append(repositoryList, tr)

	return repositoryList, nil
}

// GameRepository -
func (m *Model) GameRepository() *game.Repository {

	r := m.Repositories[game.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", game.TableName)
		return nil
	}

	return r.(*game.Repository)
}
