package error

import "fmt"

type ValidationErrorType int

const (
	ValidationErrorUnsupported ValidationErrorType = iota
	ValidationErrorInvalid
)

func (t ValidationErrorType) String() string {
	return [...]string{"unsupported", "invalid"}[t]
}

type LinkedFields struct {
	LinkedField string
	FieldA      string
	FieldB      string
}

func CreateRegistry(et ValidationErrorType, fields []string) Registry {
	errorCollection := Registry{}

	for _, f := range fields {
		errCode := CreateErrorCode(et, f)
		message := fmt.Sprintf("The property '%s' is %s.", f, et)

		var e Error
		if et == ValidationErrorInvalid {
			e, _ = ToError(NewInvalidError(f, message))
		} else {
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
		message := fmt.Sprintf("The combination of %s and %s is %s.", f.FieldA, f.FieldB, et)

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
