package server

func (s *Server) watchAndReload() error {
	site := s.Site
	return site.WatchFiles(func(filenames []string) {
		// This resolves filenames to URLS *before* reloading the site, in case the latter
		// remaps permalinks.
		urls := map[string]bool{}
		for _, relpath := range filenames {
			url, ok := site.FilenameURLPath(relpath)
			if ok {
				urls[url] = true
			}
		}
		s.reload(len(filenames))
		for url := range urls {
			s.lr.Reload(url)
		}
	})
}
