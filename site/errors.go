package site

import (
	"github.com/osteele/gojekyll/utils"
)

// combineErrors is a wrapper for utils.CombineErrors for backward compatibility
func combineErrors(errs []error) error {
	return utils.CombineErrors(errs)
}
