package site

import (
	"github.com/osteele/gojekyll/plugins"
	"github.com/osteele/gojekyll/utils"
)

func (s *Site) installPlugins() error {
	names := s.cfg.Plugins
	installed := map[string]bool{}
	// install plugins and call their ModifyPluginList lists
	for {
		pending := []string{}
		for _, name := range names {
			if !installed[name] {
				pending = append(pending, name)
			}
		}
		if len(pending) == 0 {
			break
		}
		if err := plugins.Install(pending, s); err != nil {
			return err
		}
		for _, name := range names {
			if !installed[name] {
				p, ok := plugins.Lookup(name)
				if ok {
					names = p.ModifyPluginList(names)
				}
			}
		}
		for _, name := range pending {
			installed[name] = true
		}
	}
	s.plugins = names
	return nil
}

func (s *Site) runHooks(h func(plugins.Plugin) error) error {
	for _, name := range s.plugins {
		p, ok := plugins.Lookup(name)
		if ok {
			if err := h(p); err != nil {
				return utils.WrapError(err, "running plugin")
			}
		}
	}
	return nil
}
