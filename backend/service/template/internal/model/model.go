package model

import (
	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-mud/backend/core/model"
	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/preparer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/querier"
	"gitlab.com/alienspaces/go-mud/backend/core/type/repositor"
	"gitlab.com/alienspaces/go-mud/backend/core/type/storer"

	"gitlab.com/alienspaces/go-mud/backend/service/template/internal/repository/template"
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
	m.QueriesFunc = m.NewQueriers

	return m, nil
}

// NewQueriers - Custom queriers for this model
func (m *Model) NewQueriers(p preparer.Query, tx *sqlx.Tx) ([]querier.Querier, error) {
	var querierList []querier.Querier

	return querierList, nil
}

// NewRepositories - Custom repositories for this model
func (m *Model) NewRepositories(p preparer.Repository, tx *sqlx.Tx) ([]repositor.Repositor, error) {
	l := m.Logger("NewRepositories")

	repositoryList := []repositor.Repositor{}

	tr, err := template.NewRepository(l, p, tx)
	if err != nil {
		l.Warn("Failed new template repository >%v<", err)
		return nil, err
	}

	repositoryList = append(repositoryList, tr)

	return repositoryList, nil
}

// TemplateRepository -
func (m *Model) TemplateRepository() *template.Repository {
	l := m.Logger("TemplateRepository")

	r := m.Repositories[template.TableName]
	if r == nil {
		l.Warn("Repository >%s< is nil", template.TableName)
		return nil
	}

	return r.(*template.Repository)
}

// Logger - Returns a logger with package context and provided function context
func (m *Model) Logger(functionName string) logger.Logger {
	return m.Log.WithPackageContext("model").WithFunctionContext(functionName)
}
