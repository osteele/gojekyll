package server

import (
	"bytes"
	"io"
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
		w := TagInjector{out, []byte(":insertion:")}
		_, err := io.WriteString(w, test.in)
		require.NoError(t, err)
		require.Equal(t, test.out, out.String())
	}
}
