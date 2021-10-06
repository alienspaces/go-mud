package model

import (
	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-boilerplate/server/core/model"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/configurer"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/preparer"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/repositor"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/storer"

	"gitlab.com/alienspaces/go-boilerplate/server/service/player/internal/repository/player"
	"gitlab.com/alienspaces/go-boilerplate/server/service/player/internal/repository/playerrole"
)

// Model -
type Model struct {
	model.Model
	// Allows auth token verification to be mocked for testing
	AuthVerifyTokenFunc func(provider, token string) (*VerifiedData, error)
}

// NewModel -
func NewModel(c configurer.Configurer, l logger.Logger, s storer.Storer) (*Model, error) {

	m := &Model{
		Model: model.Model{
			Config: c,
			Log:    l,
			Store:  s,
		},
		AuthVerifyTokenFunc: nil,
	}

	m.AuthVerifyTokenFunc = m.authVerifyToken
	m.RepositoriesFunc = m.NewRepositories

	return m, nil
}

// NewRepositories - Custom repositories for this model
func (m *Model) NewRepositories(p preparer.Preparer, tx *sqlx.Tx) ([]repositor.Repositor, error) {

	repositoryList := []repositor.Repositor{}

	playerRepo, err := player.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new player repository >%v<", err)
		return nil, err
	}

	repositoryList = append(repositoryList, playerRepo)

	playerRoleRepo, err := playerrole.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new player role repository >%v<", err)
		return nil, err
	}

	repositoryList = append(repositoryList, playerRoleRepo)

	return repositoryList, nil
}

// PlayerRepository -
func (m *Model) PlayerRepository() *player.Repository {

	r := m.Repositories[player.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", player.TableName)
		return nil
	}

	return r.(*player.Repository)
}

// PlayerRoleRepository -
func (m *Model) PlayerRoleRepository() *playerrole.Repository {

	r := m.Repositories[playerrole.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", playerrole.TableName)
		return nil
	}

	return r.(*playerrole.Repository)
}
