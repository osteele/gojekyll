package server

import (
	"io"
	"os"
	"regexp"

	"github.com/jaschaephraim/lrserver"
)

// liveReloadScriptTag is inserted into the HTML page.
var liveReloadScriptTag = []byte(`<script src="http://localhost:35729/livereload.js"></script>`)

// startLiveReloader starts the Live Reload server as a go routine, and returns immediately
func (s *Server) startLiveReloader() error {
	lr := lrserver.New(lrserver.DefaultName, lrserver.DefaultPort)
	s.lr = lr
	lr.SetStatusLog(nil)
	lr.ErrorLog().SetOutput(outputFilter{os.Stdout})
	go lr.ListenAndServe() // nolint: errcheck
	return nil
}

// NewLiveReloadInjector returns a writer that injects the Live Reload JavaScript
// into its wrapped content.
func NewLiveReloadInjector(w io.Writer) io.Writer {
	return TagInjector{w, liveReloadScriptTag}
}

// Remove the lines that match the exclusion pattern.
// TODO submit an upstream PR to make this unnecessary
type outputFilter struct{ w *os.File }

var excludeRE = regexp.MustCompile(`websocket: close 1006 \(abnormal closure\): unexpected EOF`)

func (w outputFilter) Write(b []byte) (int, error) {
	if excludeRE.Match(b) {
		return len(b), nil
	}
	return w.w.Write(b)
}
