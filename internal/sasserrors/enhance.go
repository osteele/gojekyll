package sasserrors

import (
	"fmt"
	"strings"
)

// Enhance adds helpful context to common sass.Start() errors
func Enhance(err error) error {
	if err == nil {
		return nil
	}

	errStr := err.Error()

	// Detect the common error when npm 'sass' package is used instead of 'sass-embedded'
	if strings.Contains(errStr, "unexpected EOF") || strings.Contains(errStr, "connection is shut down") {
		return fmt.Errorf("%w\n\n"+
			"This error typically occurs when the wrong Sass package is installed.\n"+
			"gojekyll requires the Dart Sass Embedded binary, not the pure JavaScript version.\n\n"+
			"Solution:\n"+
			"  - If using npm: install 'sass-embedded' instead of 'sass'\n"+
			"    Run: npm install -g sass-embedded\n"+
			"  - Or download Dart Sass from: https://github.com/sass/dart-sass/releases\n\n"+
			"The 'sass' package from npm is a pure JavaScript implementation that does not\n"+
			"support the embedded protocol required by gojekyll", err)
	}

	return err
}
