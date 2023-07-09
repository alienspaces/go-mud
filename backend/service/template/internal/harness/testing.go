package harness

import (
	coreConfig "gitlab.com/alienspaces/go-mud/backend/core/config"
	"gitlab.com/alienspaces/go-mud/backend/core/harness"
	"gitlab.com/alienspaces/go-mud/backend/core/log"
	"gitlab.com/alienspaces/go-mud/backend/core/store"
	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/backend/core/type/storer"
	templateConfig "gitlab.com/alienspaces/go-mud/backend/service/template/internal/config"

	"gitlab.com/alienspaces/go-mud/backend/service/template/internal/model"
	"gitlab.com/alienspaces/go-mud/backend/service/template/internal/record"
)

// Testing -
type Testing struct {
	harness.Testing
	Data         Data
	DataConfig   DataConfig
	teardownData teardownData
}

// teardownData -
type teardownData struct {
	TemplateRecs []record.Template
}

// NewTesting -
func NewTesting(c configurer.Configurer, l logger.Logger, s storer.Storer, config DataConfig) (t *Testing, err error) {

	// harness
	t = &Testing{}
	t.Config = c
	t.Log = l
	t.Store = s

	// modeller
	t.ModellerFunc = t.Modeller

	// data
	t.CreateDataFunc = t.CreateData
	t.RemoveDataFunc = t.RemoveData

	t.DataConfig = config
	t.Data = Data{}
	t.teardownData = teardownData{}

	err = t.Init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

// NewDefaultDependencies provides a set of default dependencies.
func NewDefaultDependencies(items ...coreConfig.Item) (configurer.Configurer, logger.Logger, storer.Storer, error) {

	// configurer
	c, err := templateConfig.NewConfig(items, false)
	if err != nil {
		return nil, nil, nil, err
	}

	// logger
	l, err := log.NewLogger(c)
	if err != nil {
		return nil, nil, nil, err
	}

	// storer
	s, err := store.NewStore(c, l)
	if err != nil {
		return nil, nil, nil, err
	}

	return c, l, s, nil
}

// Modeller -
func (t *Testing) Modeller() (modeller.Modeller, error) {

	m, err := model.NewModel(t.Config, t.Log, t.Store)
	if err != nil {
		t.Log.Warn("failed new model >%v<", err)
		return nil, err
	}

	return m, nil
}

// CreateData - Custom data
func (t *Testing) CreateData() error {
	l := t.Logger("CreateData")

	for _, templateConfig := range t.DataConfig.TemplateConfig {
		l.Debug("processing template configuration")

		_, err := t.processTemplateConfig(templateConfig)
		if err != nil {
			l.Warn("failed processing template configuration >%v<", err)
			return err
		}
	}

	return nil
}

// RemoveData -
func (t *Testing) RemoveData() error {
	l := t.Logger("RemoveData")

	for _, rec := range t.teardownData.TemplateRecs {
		err := t.removeTemplateRec(&rec)
		if err != nil {
			l.Warn("failed removing template record >%v<", err)
			return err
		}
	}

	t.Data = Data{}
	t.teardownData = teardownData{}

	return nil
}

// Logger - Returns a logger with package context and provided function context
func (t *Testing) Logger(fCtx string) logger.Logger {
	return t.Log.WithPackageContext("harness").WithFunctionContext(fCtx)
}
