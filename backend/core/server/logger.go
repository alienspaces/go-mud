package server

import (
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
)

// loggerWithContext provides a contextual logger for usage in HTTP handler methods
func loggerWithContext(l logger.Logger, functionName string) logger.Logger {
	if l == nil {
		return nil
	}
	return l.WithPackageContext("core/server/http").WithFunctionContext(functionName)
}
