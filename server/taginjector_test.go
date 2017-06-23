package server

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

var tests = []struct{ in, out string }{
	{"pre</head>post", "pre:insertion:</head>post"},
	{"pre:insertion:</head>post", "pre:insertion:</head>post"},
	{"post", ":insertion:post"},
}

func TestTagInjector(t *testing.T) {
	for _, test := range tests {
		out := new(bytes.Buffer)
		// bytes.NewBufferString(c.in)
		w := TagInjector{out, []byte(":insertion:")}
		_, err := w.Write([]byte(test.in))
		require.NoError(t, err)
		// buf.String() // returns a string of what was written to it
		require.Equal(t, test.out, out.String())
	}
}
