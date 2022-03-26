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
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/dungeonaction"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/dungeonactioncharacter"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/dungeonactioncharacterobject"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/dungeonactionmonster"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/dungeonactionmonsterobject"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/dungeonactionobject"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/dungeoncharacter"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/dungeonlocation"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/dungeonmonster"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/dungeonobject"
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
func (m *Model) NewRepositories(p preparer.Repository, tx *sqlx.Tx) ([]repositor.Repositor, error) {

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

	dungeonObjectRepo, err := dungeonobject.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new dungeon object repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, dungeonObjectRepo)

	dungeonActionRepo, err := dungeonaction.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new dungeon action repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, dungeonActionRepo)

	dungeonActionCharacterRepo, err := dungeonactioncharacter.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new dungeon action character repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, dungeonActionCharacterRepo)

	dungeonActionCharacterObjectRepo, err := dungeonactioncharacterobject.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new dungeon action character object repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, dungeonActionCharacterObjectRepo)

	dungeonActionMonsterRepo, err := dungeonactionmonster.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new dungeon action monster repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, dungeonActionMonsterRepo)

	dungeonActionMonsterObjectRepo, err := dungeonactionmonsterobject.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new dungeon action monster object repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, dungeonActionMonsterObjectRepo)

	dungeonActionObjectRepo, err := dungeonactionobject.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new dungeon action object repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, dungeonActionObjectRepo)

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

// DungeonObjectRepository -
func (m *Model) DungeonObjectRepository() *dungeonobject.Repository {

	r := m.Repositories[dungeonobject.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", dungeonobject.TableName)
		return nil
	}

	return r.(*dungeonobject.Repository)
}

// DungeonActionRepository -
func (m *Model) DungeonActionRepository() *dungeonaction.Repository {

	r := m.Repositories[dungeonaction.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", dungeonaction.TableName)
		return nil
	}

	return r.(*dungeonaction.Repository)
}

// DungeonActionCharacterRepository -
func (m *Model) DungeonActionCharacterRepository() *dungeonactioncharacter.Repository {

	r := m.Repositories[dungeonactioncharacter.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", dungeonactioncharacter.TableName)
		return nil
	}

	return r.(*dungeonactioncharacter.Repository)
}

// DungeonActionCharacterObjectRepository -
func (m *Model) DungeonActionCharacterObjectRepository() *dungeonactioncharacterobject.Repository {

	r := m.Repositories[dungeonactioncharacterobject.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", dungeonactioncharacterobject.TableName)
		return nil
	}

	return r.(*dungeonactioncharacterobject.Repository)
}

// DungeonActionMonsterRepository -
func (m *Model) DungeonActionMonsterRepository() *dungeonactionmonster.Repository {

	r := m.Repositories[dungeonactionmonster.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", dungeonactionmonster.TableName)
		return nil
	}

	return r.(*dungeonactionmonster.Repository)
}

// DungeonActionMonsterObjectRepository -
func (m *Model) DungeonActionMonsterObjectRepository() *dungeonactionmonsterobject.Repository {

	r := m.Repositories[dungeonactionmonsterobject.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", dungeonactionmonsterobject.TableName)
		return nil
	}

	return r.(*dungeonactionmonsterobject.Repository)
}

// DungeonActionObjectRepository -
func (m *Model) DungeonActionObjectRepository() *dungeonactionobject.Repository {

	r := m.Repositories[dungeonactionobject.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", dungeonactionobject.TableName)
		return nil
	}

	return r.(*dungeonactionobject.Repository)
}
