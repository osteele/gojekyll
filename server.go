package gojekyll

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jaschaephraim/lrserver"
)

// Server serves the site on HTTP.
type Server struct {
	Site *Site
	mu   sync.Mutex
	lr   *lrserver.Server
}

// Run runs the server.
func (s *Server) Run(logger func(label, value string)) error {
	address := "localhost:4000"
	if err := s.watchFiles(); err != nil {
		return err
	}
	s.lr = lrserver.New(lrserver.DefaultName, lrserver.DefaultPort)
	go s.lr.ListenAndServe()
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
	var (
		site      = s.Site
		events    = make(chan string)
		debounced = debounce(time.Second, events)
	)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				events <- event.Name
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	go func() {
		for {
			names := <-debounced
			// Resolve to URLS *before* reloading the site, in case the latter
			// remaps permalinks.
			urls := map[string]bool{}
			for _, name := range names {
				relpath, err := filepath.Rel(site.Source, name)
				if err != nil {
					log.Println("error:", err)
					continue
				}
				url, found := site.GetFileURL(relpath)
				if !found {
					log.Println("error:", name, "does not match a site URL")
				}
				urls[url] = true
			}
			s.syncReloadSite()
			for url := range urls {
				s.lr.Reload(url)
			}
		}
	}()

	return watcher.Add(site.Source)
}

// debounce relays values from input to output, merging successive values within interval
// TODO consider https://github.com/ReactiveX/RxGo
func debounce(interval time.Duration, input chan string) (output chan []string) {
	output = make(chan []string)
	var (
		pending = []string{}
		ticker  = time.Tick(interval) // nolint: staticcheck
	)
	go func() {
		for {
			select {
			case value := <-input:
				pending = append(pending, value)
			case <-ticker:
				if len(pending) > 0 {
					output <- pending
					pending = []string{}
				}
			}
		}
	}()
	return
}
