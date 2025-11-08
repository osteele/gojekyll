package utils

import (
	"fmt"
	"os"
)

// A WrappedError decorates an error with a message
type WrappedError interface {
	error
	Cause() error
}

// WrapError returns an error decorated with a message.
// If the error is nil, it returns nil.
func WrapError(err error, m string) error {
	if err == nil {
		return nil
	}
	return &wrappedError{cause: err, message: m}
}

type wrappedError struct {
	cause   error
	message string
}

func (we *wrappedError) Cause() error {
	return we.cause
}

func (we *wrappedError) Error() string {
	return fmt.Sprintf("%s: %s", we.message, we.cause)
}

// A PathError is an error with a source path.
//
// An os.PathError is unfortunately not a PathError, but this is still
// useful for deciding whether to wrap other errors.
type PathError interface {
	WrappedError
	Path() string
}

type pathError struct {
	cause error
	path  string
}

func (pe *pathError) Cause() error {
	return pe.cause
}

func (pe *pathError) Error() string {
	return fmt.Sprintf("%s: %s", pe.path, pe.cause)
}

func (pe *pathError) Path() string {
	return pe.path
}

// WrapPathError returns an error that will print with a path.\
// It wraps its argument if it is not nil and does not already provide a path.
func WrapPathError(err error, path string) error {
	if err == nil {
		return nil
	}
	switch err := err.(type) {
	case PathError:
		return err
	case *os.PathError:
		return err
	default:
		return &pathError{path: path, cause: err}
	}
}

// CombineErrors combines multiple errors into a single error.
// Returns nil if there are no errors, the single error if there's only one,
// or a combined error with all error messages joined by newlines.
func CombineErrors(errs []error) error {
	// Filter out nil errors
	var nonNilErrs []error
	for _, err := range errs {
		if err != nil {
			nonNilErrs = append(nonNilErrs, err)
		}
	}

	switch len(nonNilErrs) {
	case 0:
		return nil
	case 1:
		return nonNilErrs[0]
	default:
		var msg string
		for i, e := range nonNilErrs {
			if i > 0 {
				msg += "\n"
			}
			msg += e.Error()
		}
		return fmt.Errorf("%s", msg)
	}
}
