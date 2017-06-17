package gojekyll

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"path"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Server serves the site on HTTP.
type Server struct {
	Site *Site
	mu   sync.Mutex
}

type emptyType struct{}

var void emptyType

// Run runs the server.
func (s *Server) Run(logger func(label, value string)) error {
	address := "localhost:4000"
	if err := s.watchFiles(); err != nil {
		return err
	}
	logger("Server address:", "http://"+address+"/")
	logger("Server running...", "press ctrl-c to stop.")
	http.HandleFunc("/", s.handler)
	return http.ListenAndServe(address, nil)
}

func (s *Server) handler(rw http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	site := s.Site
	urlpath := r.URL.Path
	mimeType := mime.TypeByExtension(path.Ext(urlpath))
	if mimeType != "" {
		rw.Header().Set("Content-Type", mimeType)
	}

	p, found := site.PageForURL(urlpath)
	if !found {
		rw.WriteHeader(http.StatusNotFound)
		p, found = site.Paths["404.html"]
	}
	if !found {
		fmt.Fprintf(rw, "404 page not found: %s", urlpath)
		return
	}

	err := p.Write(rw)
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

func (s *Server) watchFiles() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	var (
		events    = make(chan emptyType)
		debounced = debounce(time.Second, events)
	)

	go func() {
		for {
			select {
			case <-watcher.Events:
				events <- void
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	go func() {
		for {
			<-debounced
			s.syncReloadSite()
		}
	}()

	return watcher.Add(s.Site.Source)
}

// debounce relays values from input to output, merging successive values within interval
// TODO consider https://github.com/ReactiveX/RxGo
func debounce(interval time.Duration, input chan emptyType) (output chan emptyType) {
	output = make(chan emptyType)
	var (
		pending = false
		ticker  = time.Tick(interval) // nolint: staticcheck
	)
	go func() {
		for {
			select {
			case <-input:
				pending = true
			case <-ticker:
				if pending {
					output <- void
					pending = false
				}
			}
		}
	}()
	return
}
