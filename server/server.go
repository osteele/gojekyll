package server

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/jaschaephraim/lrserver"
	"github.com/osteele/gojekyll/sites"
	"github.com/pkg/browser"
)

// Server serves the site on HTTP.
type Server struct {
	Site *sites.Site
	mu   sync.Mutex
	lr   *lrserver.Server
}

// Run runs the server.
func (s *Server) Run(open bool, logger func(label, value string)) error {
	address := "localhost:4000"
	if err := s.watchFiles(); err != nil {
		return err
	}
	s.lr = lrserver.New(lrserver.DefaultName, lrserver.DefaultPort)
	s.lr.SetStatusLog(nil)
	logger("Server address:", "http://"+address+"/")
	http.HandleFunc("/", s.handler)
	c := make(chan error)
	go s.lr.ListenAndServe() // nolint: errcheck
	go func() {
		c <- http.ListenAndServe(address, nil)
	}()
	logger("Server running...", "press ctrl-c to stop.")
	if open {
		if err := browser.OpenURL("http://" + address); err != nil {
			fmt.Println("Error opening page:", err)
		}
	}
	return <-c
}

func (s *Server) handler(rw http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	site := s.Site
	urlpath := r.URL.Path

	page, found := site.URLPage(urlpath)
	if !found {
		rw.WriteHeader(http.StatusNotFound)
		page, found = site.Paths["404.html"]
	}
	if !found {
		fmt.Fprintf(rw, "404 page not found: %s", urlpath) // nolint: gas
		return
	}

	mimeType := mime.TypeByExtension(page.OutputExt())
	if mimeType != "" {
		rw.Header().Set("Content-Type", mimeType)
	}
	var w io.Writer = rw
	if strings.HasPrefix(mimeType, "text/html;") {
		w = NewLiveReloadInjector(w)
	}
	err := page.Write(site, w)
	if err != nil {
		fmt.Printf("Error rendering %s: %s", urlpath, err)
	}
}

func (s *Server) syncReloadSite() {
	s.mu.Lock()
	defer s.mu.Unlock()

	start := time.Now()
	fmt.Printf("%s Reloading site...", start.Format(time.Stamp))
	if err := s.Site.Reload(); err != nil {
		fmt.Println()
		fmt.Println(err.Error())
	}
	fmt.Printf("reloaded in %.2fs\n", time.Since(start).Seconds())
}
