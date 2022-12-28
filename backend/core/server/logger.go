package server

import (
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
)

// Logger provides a contextual logger
func (rnr *Runner) Logger(functionName string) logger.Logger {
	return rnr.Log.WithPackageContext("core/server").WithFunctionContext(functionName)
}

// HTTPLogger provides a contextual logger for usage in HTTP handler methods
func HTTPLogger(l logger.Logger, functionName string) logger.Logger {
	if l == nil {
		return nil
	}
	return l.WithPackageContext("core/server/http").WithFunctionContext(functionName)
}
