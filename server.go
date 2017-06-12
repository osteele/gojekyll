package main

import (
	"fmt"
	"net/http"
)

func server() error {
	address := "localhost:4000"
	printSetting("Server address:", "http://"+address+"/")
	printSetting("Server running...", "press ctrl-c to stop.")
	http.HandleFunc("/", handler)
	return http.ListenAndServe(address, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	p, found := siteMap[path]
	if !found {
		w.WriteHeader(http.StatusNotFound)
		p, found = siteMap["404.html"]
	}
	if !found {
		fmt.Fprintf(w, "404 page not found: %s", path)
		return
	}

	p, err := readFile(p.Path, siteData, true)
	if err != nil {
		fmt.Printf("Error rendering %s: %s", path, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if _, err := w.Write(p.Body); err != nil {
		fmt.Printf("Error writing %s: %s", path, err)
	}
}
