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
	err := http.ListenAndServe(address, nil)
	if err != nil {
		// TODO pick another port
		return err
	}
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header

	p, found := siteMap[path]
	if !found {
		w.WriteHeader(http.StatusNotFound)
		p, found = siteMap["404.html"]
	}
	if !found {
		fmt.Fprintf(w, "404 page not found: %s", path)
		return
	}

	p, err := readFile(p.Path, true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Error expanding %s\n%s", p.Path, err)
		fmt.Fprintf(w, "Error expanding %s\n%s", p.Path, err)
	}
	w.Write(p.Body)
}
