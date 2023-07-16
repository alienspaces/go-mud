package error

import (
	"errors"
	"fmt"
	"strings"
)

func IsError(e error) bool {
	var errorPtr Error
	return errors.As(e, &errorPtr)
}

func HasErrorCode(err error, c ErrorCode) bool {
	e, err := ToError(err)
	if err != nil {
		return false
	}

	return e.ErrorCode == c
}

func ToError(e error) (Error, error) {
	if e == nil {
		return Error{}, fmt.Errorf("err is nil when converting to coreerror.Error type")
	}

	var err Error
	if !errors.As(e, &err) {
		return Error{}, fmt.Errorf("failed to convert to coreerror.Error type >%v<", e)
	}

	if len(err.SchemaValidationErrors) == 0 {
		err.SchemaValidationErrors = nil
	}

	return err, nil
}

func ToErrors(errs ...error) ([]Error, error) {
	var results []Error

	for _, e := range errs {
		result, err := ToError(e)
		if err != nil {
			return nil, fmt.Errorf("failed to convert err to error result >%#v<", err)
		}

		results = append(results, result)
	}

	return results, nil
}

func ProcessParamError(err error) error {
	e, conversionErr := ToError(err)
	if conversionErr != nil {
		return err
	}

	if len(e.SchemaValidationErrors) == 0 {
		return NewParamError(e.Error())
	}

	errStr := strings.Builder{}
	errStr.WriteString("Invalid parameter(s): ")
	for i, sve := range e.SchemaValidationErrors {
		if sve.GetField() == "$" {
			errStr.WriteString(fmt.Sprintf("(%d) %s; ", i+1, sve.Message))
		} else {
			errStr.WriteString(fmt.Sprintf("(%d) %s: %s; ", i+1, sve.GetField(), sve.Message))
		}
	}

	formattedErrString := errStr.String()
	formattedErrString = formattedErrString[0 : len(formattedErrString)-2] // remove extra space and semicolon
	formattedErrString += "."
	return NewParamError(formattedErrString)
}
