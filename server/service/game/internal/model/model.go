package model

import (
	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-mud/server/core/model"
	"gitlab.com/alienspaces/go-mud/server/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/preparer"
	"gitlab.com/alienspaces/go-mud/server/core/type/repositor"
	"gitlab.com/alienspaces/go-mud/server/core/type/storer"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/action"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/actioncharacter"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/actioncharacterobject"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/actionmonster"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/actionmonsterobject"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/actionobject"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/character"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/characterinstance"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/characterinstanceview"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/dungeon"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/dungeoninstance"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/location"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/locationinstance"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/locationinstanceview"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/monster"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/monsterinstance"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/monsterinstanceview"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/object"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/objectinstance"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/repository/objectinstanceview"
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

	dungeonInstanceRepo, err := dungeoninstance.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new dungeon instance repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, dungeonInstanceRepo)

	locationRepo, err := location.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new location repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, locationRepo)

	locationInstanceRepo, err := locationinstance.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new location instance repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, locationInstanceRepo)

	locationInstanceViewRepo, err := locationinstanceview.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new location instance view repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, locationInstanceViewRepo)

	characterRepo, err := character.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new character repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, characterRepo)

	characterInstanceRepo, err := characterinstance.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new character instance repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, characterInstanceRepo)

	characterInstanceViewRepo, err := characterinstanceview.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new character instance view repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, characterInstanceViewRepo)

	monsterRepo, err := monster.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new monster repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, monsterRepo)

	monsterInstanceRepo, err := monsterinstance.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new monster instance repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, monsterInstanceRepo)

	monsterInstanceViewRepo, err := monsterinstanceview.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new monster instance view repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, monsterInstanceViewRepo)

	objectRepo, err := object.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new object repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, objectRepo)

	objectInstanceRepo, err := objectinstance.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new object instance repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, objectInstanceRepo)

	objectInstanceViewRepo, err := objectinstanceview.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new object instance view repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, objectInstanceViewRepo)

	dungeonActionRepo, err := action.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new dungeon action repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, dungeonActionRepo)

	dungeonActionCharacterRepo, err := actioncharacter.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new dungeon action character repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, dungeonActionCharacterRepo)

	dungeonActionCharacterObjectRepo, err := actioncharacterobject.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new dungeon action character object repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, dungeonActionCharacterObjectRepo)

	dungeonActionMonsterRepo, err := actionmonster.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new dungeon action monster repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, dungeonActionMonsterRepo)

	dungeonActionMonsterObjectRepo, err := actionmonsterobject.NewRepository(m.Log, p, tx)
	if err != nil {
		m.Log.Warn("Failed new dungeon action monster object repository >%v<", err)
		return nil, err
	}
	repositoryList = append(repositoryList, dungeonActionMonsterObjectRepo)

	dungeonActionObjectRepo, err := actionobject.NewRepository(m.Log, p, tx)
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

// DungeonInstanceRepository -
func (m *Model) DungeonInstanceRepository() *dungeoninstance.Repository {

	r := m.Repositories[dungeoninstance.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", dungeoninstance.TableName)
		return nil
	}

	return r.(*dungeoninstance.Repository)
}

// LocationRepository -
func (m *Model) LocationRepository() *location.Repository {

	r := m.Repositories[location.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", location.TableName)
		return nil
	}

	return r.(*location.Repository)
}

// LocationInstanceRepository -
func (m *Model) LocationInstanceRepository() *locationinstance.Repository {

	r := m.Repositories[locationinstance.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", locationinstance.TableName)
		return nil
	}

	return r.(*locationinstance.Repository)
}

// LocationInstanceViewRepository -
func (m *Model) LocationInstanceViewRepository() *locationinstanceview.Repository {

	r := m.Repositories[locationinstanceview.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", locationinstanceview.TableName)
		return nil
	}

	return r.(*locationinstanceview.Repository)
}

// CharacterRepository -
func (m *Model) CharacterRepository() *character.Repository {

	r := m.Repositories[character.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", character.TableName)
		return nil
	}

	return r.(*character.Repository)
}

// CharacterInstanceRepository -
func (m *Model) CharacterInstanceRepository() *characterinstance.Repository {

	r := m.Repositories[characterinstance.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", characterinstance.TableName)
		return nil
	}

	return r.(*characterinstance.Repository)
}

// CharacterInstanceViewRepository -
func (m *Model) CharacterInstanceViewRepository() *characterinstanceview.Repository {

	r := m.Repositories[characterinstanceview.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", characterinstanceview.TableName)
		return nil
	}

	return r.(*characterinstanceview.Repository)
}

// MonsterRepository -
func (m *Model) MonsterRepository() *monster.Repository {

	r := m.Repositories[monster.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", monster.TableName)
		return nil
	}

	return r.(*monster.Repository)
}

// MonsterInstanceRepository -
func (m *Model) MonsterInstanceRepository() *monsterinstance.Repository {

	r := m.Repositories[monsterinstance.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", monsterinstance.TableName)
		return nil
	}

	return r.(*monsterinstance.Repository)
}

// MonsterInstanceViewRepository -
func (m *Model) MonsterInstanceViewRepository() *monsterinstanceview.Repository {

	r := m.Repositories[monsterinstanceview.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", monsterinstanceview.TableName)
		return nil
	}

	return r.(*monsterinstanceview.Repository)
}

// ObjectRepository -
func (m *Model) ObjectRepository() *object.Repository {

	r := m.Repositories[object.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", object.TableName)
		return nil
	}

	return r.(*object.Repository)
}

// ObjectInstanceRepository -
func (m *Model) ObjectInstanceRepository() *objectinstance.Repository {

	r := m.Repositories[objectinstance.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", objectinstance.TableName)
		return nil
	}

	return r.(*objectinstance.Repository)
}

// ObjectInstanceViewRepository -
func (m *Model) ObjectInstanceViewRepository() *objectinstanceview.Repository {

	r := m.Repositories[objectinstanceview.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", objectinstanceview.TableName)
		return nil
	}

	return r.(*objectinstanceview.Repository)
}

// ActionRepository -
func (m *Model) ActionRepository() *action.Repository {

	r := m.Repositories[action.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", action.TableName)
		return nil
	}

	return r.(*action.Repository)
}

// ActionCharacterRepository -
func (m *Model) ActionCharacterRepository() *actioncharacter.Repository {

	r := m.Repositories[actioncharacter.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", actioncharacter.TableName)
		return nil
	}

	return r.(*actioncharacter.Repository)
}

// ActionCharacterObjectRepository -
func (m *Model) ActionCharacterObjectRepository() *actioncharacterobject.Repository {

	r := m.Repositories[actioncharacterobject.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", actioncharacterobject.TableName)
		return nil
	}

	return r.(*actioncharacterobject.Repository)
}

// ActionMonsterRepository -
func (m *Model) ActionMonsterRepository() *actionmonster.Repository {

	r := m.Repositories[actionmonster.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", actionmonster.TableName)
		return nil
	}

	return r.(*actionmonster.Repository)
}

// ActionMonsterObjectRepository -
func (m *Model) ActionMonsterObjectRepository() *actionmonsterobject.Repository {

	r := m.Repositories[actionmonsterobject.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", actionmonsterobject.TableName)
		return nil
	}

	return r.(*actionmonsterobject.Repository)
}

// ActionObjectRepository -
func (m *Model) ActionObjectRepository() *actionobject.Repository {

	r := m.Repositories[actionobject.TableName]
	if r == nil {
		m.Log.Warn("Repository >%s< is nil", actionobject.TableName)
		return nil
	}

	return r.(*actionobject.Repository)
}

func (m *Model) Logger(functionName string) logger.Logger {
	return m.Log.WithPackageContext("model").WithFunctionContext(functionName)
}
