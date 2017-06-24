package server

import (
	"fmt"
	"io"
	"log"
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
	err := s.Site.InitializeRenderingPipeline()
	if err != nil {
		return err
	}
	address := "localhost:4000"
	logger("Server address:", "http://"+address+"/")
	if err := s.StartLiveReloader(); err != nil {
		return err
	}
	if err := s.watchFiles(); err != nil {
		return err
	}
	http.HandleFunc("/", s.handler)
	c := make(chan error)
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

	p, found := site.URLPage(urlpath)
	if !found {
		rw.WriteHeader(http.StatusNotFound)
		log.Println("Not found:", urlpath)
		p, found = site.Paths["404.html"]
	}
	if !found {
		fmt.Fprintf(rw, "404 page not found: %s", urlpath) // nolint: gas
		return
	}

	mimeType := mime.TypeByExtension(p.OutputExt())
	if mimeType != "" {
		rw.Header().Set("Content-Type", mimeType)
	}
	var w io.Writer = rw
	if strings.HasPrefix(mimeType, "text/html;") {
		w = NewLiveReloadInjector(w)
	}
	err := p.Write(site, w)
	if err != nil {
		fmt.Printf("Error rendering %s: %s", urlpath, err)
	}
}

func (s *Server) reloadSite() {
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
