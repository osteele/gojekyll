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
	"github.com/osteele/gojekyll/site"
	"github.com/pkg/browser"
)

// Server serves the site on HTTP.
type Server struct {
	Site *site.Site
	mu   sync.Mutex
	lr   *lrserver.Server
}

// Run runs the server.
func (s *Server) Run(open bool, logger func(label, value string)) error {
	cfg := s.Site.Config()
	s.Site.SetAbsoluteURL("")
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	logger("Server address:", "http://"+address+"/")
	if err := s.StartLiveReloader(); err != nil {
		return err
	}
	if err := s.watchAndReload(); err != nil {
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

	var (
		site    = s.Site
		urlpath = r.URL.Path
	)
	p, found := site.URLPage(urlpath)

	if !found {
		rw.WriteHeader(http.StatusNotFound)
		log.Println("Not found:", urlpath)
		p, found = site.Routes["404.html"]
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
	err := site.WriteDocument(w, p)
	if err != nil {
		fmt.Printf("Error rendering %s: %s", urlpath, err)
	}
}

func (s *Server) reloadSite() {
	s.mu.Lock()
	defer s.mu.Unlock()

	start := time.Now()
	fmt.Printf("%s Reloading site...", start.Format(time.Stamp))
	site, err := s.Site.Reload()
	if err != nil {
		fmt.Println()
		fmt.Println(err.Error())
	}
	s.Site = site
	s.Site.SetAbsoluteURL("")
	fmt.Printf("reloaded in %.2fs\n", time.Since(start).Seconds())
}
