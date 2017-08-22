package utils

import (
	"fmt"
	"os"
)

// A PathError is an error with a source path.
//
// An os.PathError is unfortunately not a PathError, but this is still
// useful for deciding whether to wrap other errors.
type PathError interface {
	error
	Path() string
}

type pathError struct {
	cause error
	path  string
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
