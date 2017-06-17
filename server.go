package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fsnotify/fsnotify"
)

type Server struct{ Site *Site }

func server(site *Site) error {
	server := Server{site}
	address := "localhost:4000"
	if err := server.watchFiles(); err != nil {
		return err
	}
	printSetting("Server address:", "http://"+address+"/")
	printSetting("Server running...", "press ctrl-c to stop.")
	http.HandleFunc("/", server.handler)
	return http.ListenAndServe(address, nil)
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	site := s.Site
	urlpath := r.URL.Path

	// w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	p, found := site.Paths[urlpath]
	if !found {
		w.WriteHeader(http.StatusNotFound)
		p, found = site.Paths["404.html"]
	}
	if !found {
		fmt.Fprintf(w, "404 page not found: %s", urlpath)
		return
	}

	err := p.Write(w)
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
