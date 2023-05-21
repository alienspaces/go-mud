package server

import (
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
)

func Logger(l logger.Logger, functionName string) logger.Logger {
	return l.WithPackageContext("core/server").WithFunctionContext(functionName)
}
