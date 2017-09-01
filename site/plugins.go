package site

import (
	"github.com/osteele/gojekyll/plugins"
	"github.com/osteele/gojekyll/utils"
)

func (s *Site) installPlugins() error {
	s.plugins = s.cfg.Plugins
	installed := utils.StringSet{}
	// Install plugins and call their ModifyPluginList methods.
	// Repeat until no plugins have been added.
	for len(s.plugins) > len(installed) {
		// Collect plugins into a list instead of map, in order to preserve order
		pending := utils.StringList(s.plugins).Reject(installed.Contains)
		if err := plugins.Install(pending, s); err != nil {
			return err
		}
		for _, name := range pending {
			p, ok := plugins.Lookup(name)
			if ok {
				s.plugins = p.ModifyPluginList(s.plugins)
			}
		}
		installed.AddStrings(pending)
	}
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
