package runner

import (
	"fmt"

	coreconfig "gitlab.com/alienspaces/go-mud/backend/core/config"
	"gitlab.com/alienspaces/go-mud/backend/core/server"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
)

func buildHandlerConfig(r Runner) (map[string]server.HandlerConfig, error) {

	type handlerConfigFunc func(l logger.Logger, config map[string]server.HandlerConfig) (map[string]server.HandlerConfig, error)

	handlerConfigFuncs := []handlerConfigFunc{
		r.templateHandlerConfig,
	}

	var err error
	handlerConfig := map[string]server.HandlerConfig{}

	for _, handlerConfigFunc := range handlerConfigFuncs {
		handlerConfig, err = handlerConfigFunc(r.Log, handlerConfig)
		if err != nil {
			err := fmt.Errorf("failed adding handler configuration >%v<", err)
			r.Log.Warn(err.Error())
			return nil, err
		}
	}

	handlerConfig = resolveHandlerNames(handlerConfig)

	appHome := r.Config.Get(coreconfig.AppServerHome)
	handlerConfig, err = server.ResolveHandlerSchemaLocationRoot(handlerConfig, appHome)
	if err != nil {
		err := fmt.Errorf("failed to resolve template API handler location root >%v<", err)
		r.Log.Warn(err.Error())
		return nil, err
	}

	handlerConfig = server.ResolveHandlerSchemaLocation(handlerConfig, "schema/template")

	handlerConfig = server.ResolveDocumentationSummary(handlerConfig)

	return handlerConfig, nil
}

func resolveHandlerNames(handlerConfig map[string]server.HandlerConfig) map[string]server.HandlerConfig {
	for name := range handlerConfig {
		hc := handlerConfig[name]
		hc.Name = name
		handlerConfig[name] = hc
	}
	return handlerConfig
}
