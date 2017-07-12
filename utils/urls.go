package utils

import "strings"

// URLJoin interpolates paths with "/", skipping empty paths and avoiding "//".
func URLJoin(paths ...string) string {
	url := ""
loop:
	for _, p := range paths {
		switch {
		case p == "":
			continue loop
		case url != "" && !strings.HasSuffix(url, "/") && !strings.HasPrefix(p, "/"):
			url += "/"
		case strings.HasSuffix(url, "/") && strings.HasPrefix(p, "/"):
			p = p[1:]
		}
		url += p
	}
	return url
}
