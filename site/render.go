package site

// Render renders the site's pages.
func (s *Site) Render() error {
	for _, c := range s.Collections {
		if err := c.Render(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Site) ensureRendered() (err error) {
	s.renderOnce.Do(func() {
		err = s.initializeRenderingPipeline()
		if err != nil {
			return
		}
		err = s.Render()
		if err != nil {
			return
		}
	})
	return
}
