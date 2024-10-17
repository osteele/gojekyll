package site

import (
	"errors"
	"strings"
)

func combineErrors(errs []error) error {
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	default:
		messages := make([]string, len(errs))
		for i, e := range errs {
			messages[i] = e.Error()
		}
		return errors.New(strings.Join(messages, "\n"))
	}
}
