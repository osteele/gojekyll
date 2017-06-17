package main

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"path"

	"github.com/fsnotify/fsnotify"
)

// Server serves the site on HTTP.
type Server struct{ Site *Site }

// Run runs the server.
func (s *Server) Run() error {
	address := "localhost:4000"
	if err := s.watchFiles(); err != nil {
		return err
	}
	printSetting("Server address:", "http://"+address+"/")
	printSetting("Server running...", "press ctrl-c to stop.")
	http.HandleFunc("/", s.handler)
	return http.ListenAndServe(address, nil)
}

func (s *Server) handler(rw http.ResponseWriter, r *http.Request) {
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

func (s *Server) watchFiles() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					// TODO rebuild the site
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	return watcher.Add(s.Site.Source)
}
