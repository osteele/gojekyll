package tags

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/osteele/liquid/render"
)

var highlightArgsRE = regexp.MustCompile(`^\s*(\S+)(\s+linenos)?\s*$`)

func highlightTag(rc render.Context) (string, error) {
	argStr, err := rc.ExpandTagArg()
	if err != nil {
		return "", err
	}
	args := highlightArgsRE.FindStringSubmatch(argStr)
	if args == nil {
		return "", fmt.Errorf("syntax error")
	}
	source, err := rc.InnerString()
	if err != nil {
		return "", err
	}

	// Determine lexer.
	l := lexers.Get(args[1])
	if l == nil {
		l = lexers.Analyse(source) // nolint: misspell // British spelling from chroma library
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	lineNum := args[2] != ""

	// Determine formatter.
	f := html.New(
		html.WithClasses(true),
		html.WithLineNumbers(lineNum),
		html.LineNumbersInTable(true),
	)

	// Determine style.
	s := styles.Get("")
	if s == nil {
		s = styles.Fallback
	}

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = f.Format(buf, s, it); err != nil {
		return "", err
	}
	return buf.String(), nil
}
