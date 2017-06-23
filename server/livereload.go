package server

import (
	"io"

	"github.com/jaschaephraim/lrserver"
)

// liveReloadScriptTag is inserted into the HTML page.
var liveReloadScriptTag = []byte(`<script src="http://localhost:35729/livereload.js"></script>`)

// StartLiveReloader starts the Live Reload server as a go routine, and returns immediately
func (s *Server) StartLiveReloader() error {
	s.lr = lrserver.New(lrserver.DefaultName, lrserver.DefaultPort)
	s.lr.SetStatusLog(nil)
	go s.lr.ListenAndServe() // nolint: errcheck
	return nil
}

// NewLiveReloadInjector returns a writer that injects the Live Reload JavaScript
// into its wrapped content.
func NewLiveReloadInjector(w io.Writer) io.Writer {
	return TagInjector{w, liveReloadScriptTag}
}
