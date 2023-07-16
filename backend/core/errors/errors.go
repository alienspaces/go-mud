// Package errors is responsible for capturing accidents or misfortunes.
package errors

// Error contains an array of errors
type Error struct {
	Errors []error
}

// Error returns all errors as a string and makes fault conform with the error interface
func (e *Error) Error() string {
	var err string
	for _, error := range e.Errors {
		err = err + error.Error() + "\n"
	}
	return err
}

// Pop will remove the first error that was added to the queue
func (e *Error) Pop() {
	e.Errors = e.Errors[1:]
}

// Add will add an error to the queue of registered errors
func (e *Error) Add(err error) {
	e.Errors = append(e.Errors, err)
}

// Count returns a count of errors that have been added
func (e *Error) Count() int {
	return len(e.Errors)
}
