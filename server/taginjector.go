package server

import (
	"bytes"
	"io"
)

var closeHeadTag = []byte(`</head>`) // TODO also look for </HEAD>

// TagInjector wraps a writer and adds a script tag to its content.
// It depends on the fact that dynamic page rendering makes a single Write call,
// so that it's guaranteed to find the marker within a single invocation argument.
// It doesn't parse HTML, so it could be spoofed but probably only intentionally.
type TagInjector struct {
	w         io.Writer
	insertion []byte
}

// Write injects a livereload script tag at the end of the HTML head, if present,
// else at the beginning of the document.
func (i TagInjector) Write(b []byte) (n int, err error) {
	n = len(b)
	if !bytes.Contains(b, i.insertion) && bytes.Contains(b, closeHeadTag) {
		r := append(i.insertion, closeHeadTag...)
		b = bytes.Replace(b, closeHeadTag, r, 1)
	}
	if !bytes.Contains(b, i.insertion) {
		b = append(i.insertion, b...)
	}
	_, err = i.w.Write(b)
	return
}
