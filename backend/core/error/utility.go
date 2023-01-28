package error

import (
	"errors"
	"fmt"
)

func IsError(err error) bool {
	var errorPtr Error
	return errors.As(err, &errorPtr)
}

func ToError(err error) (Error, error) {
	if err == nil {
		return Error{}, fmt.Errorf("err is nil when converting to coreerror.Error type")
	}

	var errorPtr Error
	if !errors.As(err, &errorPtr) {
		return Error{}, fmt.Errorf("failed to convert to coreerror.Error type >%v<", err)
	}

	if len(errorPtr.SchemaValidationErrors) == 0 {
		errorPtr.SchemaValidationErrors = nil
	}

	return errorPtr, nil
}
