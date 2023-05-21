package error

import (
	"fmt"
	"strings"
)

type ValidationErrorType int

const (
	ValidationErrorUnsupported ValidationErrorType = iota
	ValidationErrorInvalid
	ValidationErrorInvalidAction
)

func (t ValidationErrorType) String() string {
	return [...]string{"unsupported", "invalid", "invalid_action"}[t]
}

type LinkedFields struct {
	LinkedField string
	Fields      []string
}

func CreateRegistry(et ValidationErrorType, fields ...string) Registry {
	errorCollection := Registry{}

	if et == ValidationErrorInvalidAction {
		e, _ := ToError(NewInvalidActionError(""))
		errorCollection[e.ErrorCode] = e
	}

	for _, f := range fields {
		errCode := CreateErrorCode(et, f)
		message := fmt.Sprintf("The property '%s' is %s.", f, et)

		var e Error
		switch et {
		case ValidationErrorInvalid:
			e, _ = ToError(NewInvalidError(f, message))
		case ValidationErrorUnsupported:
			e, _ = ToError(NewUnsupportedError(f, message))
		}
		errorCollection[errCode] = e
	}

	return errorCollection
}

func CreateLinkedRegistry(et ValidationErrorType, linkedFields []LinkedFields) Registry {
	errorCollection := Registry{}

	for _, f := range linkedFields {
		errCode := CreateErrorCode(et, f.LinkedField)
		combinationMsg := strings.Join(f.Fields, " & ")
		message := fmt.Sprintf("The combination of %s is %s.", combinationMsg, et)

		var e Error
		if et == ValidationErrorInvalid {
			e, _ = ToError(NewInvalidError(f.LinkedField, message))
		} else {
			e, _ = ToError(NewUnsupportedError(f.LinkedField, message))
		}
		errorCollection[errCode] = e
	}

	return errorCollection
}
