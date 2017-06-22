package gojekyll

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jaschaephraim/lrserver"
	"github.com/pkg/browser"
)

// Server serves the site on HTTP.
type Server struct {
	Site *Site
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

	page, found := site.PageForURL(urlpath)
	if !found {
		rw.WriteHeader(http.StatusNotFound)
		page, found = site.Paths["404.html"]
	}
	if !found {
		fmt.Fprintf(rw, "404 page not found: %s", urlpath)
		return
	}

	mimeType := mime.TypeByExtension(page.OutputExt())
	if mimeType != "" {
		rw.Header().Set("Content-Type", mimeType)
	}
	var w io.Writer = rw
	if strings.HasPrefix(mimeType, "text/html;") {
		w = scriptTagInjector{w}
	}
	err := page.Write(site, w)
	if err != nil {
		fmt.Printf("Error rendering %s: %s", urlpath, err)
	}
}

type scriptTagInjector struct {
	w io.Writer
}

var liveReloadScriptTag = []byte(`<script src="http://localhost:35729/livereload.js"></script>`)
var liveReloadSearchBytes = []byte(`</head>`)
var liveReloadReplacementBytes = append(liveReloadScriptTag, liveReloadSearchBytes...)

// Write injects a livereload script tag at the end of the HTML head, if present,
// else at the beginning of the document.
// It depends on the fact that dynamic page rendering makes a single Write call,
// so that it's guaranteed to find the marker within a single invocation argument.
// It doesn't parse HTML, so it could be spoofed but probably only intentionally.
func (i scriptTagInjector) Write(content []byte) (n int, err error) {
	if !bytes.Contains(content, liveReloadScriptTag) && bytes.Contains(content, liveReloadSearchBytes) {
		content = bytes.Replace(content, liveReloadSearchBytes, liveReloadReplacementBytes, 1)
	}
	if !bytes.Contains(content, liveReloadScriptTag) {
		content = append(liveReloadScriptTag, content...)
	}
	return i.w.Write(content)
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
