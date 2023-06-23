package error

import (
	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
)

func NewDatabaseError(l logger.Logger, err error) error {
	l.Warn("database error >%v<", err)

	return coreerror.NewInternalError()
}

func NewInvalidIDError(field string, id string) error {
	return coreerror.NewInvalidDataError("ID >%s< is not a valid UUID", id)
}
