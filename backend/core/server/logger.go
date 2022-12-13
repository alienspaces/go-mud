package server

import (
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
)

func Logger(l logger.Logger, functionName string) logger.Logger {
	if l == nil {
		return nil
	}
	return l.WithPackageContext("core/server").WithFunctionContext(functionName)
}
