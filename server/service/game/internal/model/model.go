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
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/dungeoncharacter"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/dungeonlocation"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/dungeonmonster"
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

	dungeonRepo, err := dungeon.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new dungeon repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, dungeonRepo)

	dungeonLocationRepo, err := dungeonlocation.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new dungeon location repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, dungeonLocationRepo)

	dungeonCharacterRepo, err := dungeoncharacter.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new dungeon character repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, dungeonCharacterRepo)

	dungeonMonsterRepo, err := dungeonmonster.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new dungeon monster repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, dungeonMonsterRepo)

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

// DungeonCharacterRepository -
func (m *Model) DungeonCharacterRepository() *dungeoncharacter.Repository {

	r := m.Repositories[dungeoncharacter.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", dungeoncharacter.TableName)
		return nil
	}

	return r.(*dungeoncharacter.Repository)
}

// DungeonMonsterRepository -
func (m *Model) DungeonMonsterRepository() *dungeonmonster.Repository {

	r := m.Repositories[dungeonmonster.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", dungeonmonster.TableName)
		return nil
	}

	return r.(*dungeonmonster.Repository)
}
