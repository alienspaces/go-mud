package server

import (
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
)

// loggerWithContext provides a logger with function context
func loggerWithContext(l logger.Logger, functionName string) logger.Logger {
	if l == nil {
		return nil
	}
	return l.WithPackageContext("core/server").WithFunctionContext(functionName)
}
