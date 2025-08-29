package filters

// This file contains a very inefficient implementation, that scans the string multiple times.
// Replace it by a transducer if it shows up in hot spots.

import (
	"fmt"
	"regexp"
	"strings"
)

var smartifyTransforms = []struct {
	match *regexp.Regexp
	repl  string
}{
	{regexp.MustCompile("(^|[^[:alnum:]])``(.+?)''"), "$1“$2”"},
	{regexp.MustCompile(`(^|[^[:alnum:]])'`), "$1‘"},
	{regexp.MustCompile(`'`), "’"},
	{regexp.MustCompile(`(^|[^[:alnum:]])"`), "$1“"},
	{regexp.MustCompile(`"($|[^[:alnum:]])`), "”$1"},
	{regexp.MustCompile(`(^|\s)--($|\s)`), "$1–$2"},
	{regexp.MustCompile(`(^|\s)---($|\s)`), "$1—$2"},
}

// replace these wherever they appear
var smartifyReplaceSpans = map[string]string{
	"...":  "…",
	"(c)":  "©",
	"(r)":  "®",
	"(tm)": "™",
}

// replace these only if bounded by space or word boundaries
var smartifyReplaceWords = map[string]string{
	// "---": "–",
	// "--":  "—",
}

var smartifyReplacements map[string]string
var smartifyReplacementPattern *regexp.Regexp

func init() {
	smartifyReplacements = map[string]string{}
	var disjuncts []string
	regexQuoter := regexp.MustCompile(`[\(\)\.]`)
	escape := func(s string) string {
		return regexQuoter.ReplaceAllString(s, `\$0`)
	}
	for k, v := range smartifyReplaceSpans {
		disjuncts = append(disjuncts, escape(k))
		smartifyReplacements[k] = v
	}
	for k, v := range smartifyReplaceWords {
		disjuncts = append(disjuncts, fmt.Sprintf(`(\b|\s|^)%s(\b|\s|$)`, escape(k)))
		smartifyReplacements[k] = fmt.Sprintf("$1%s$2", v)
	}
	p := fmt.Sprintf(`(%s)`, strings.Join(disjuncts, `|`))
	smartifyReplacementPattern = regexp.MustCompile(p)
}

func smartifyFilter(s string) string {
	for _, rule := range smartifyTransforms {
		s = rule.match.ReplaceAllString(s, rule.repl)
	}
	s = smartifyReplacementPattern.ReplaceAllStringFunc(s, func(w string) string {
		return smartifyReplacements[w]
	})
	return s
}
