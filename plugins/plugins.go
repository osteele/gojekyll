// Package plugins holds emulated Jekyll plugins.
//
// Unlike Jekyll, these are baked into the executable -- both because as of 2017.07 package "plugin' currently
// works only on Linux, but also because the gojekyll implementation is immature and any possible interfaces
// are far from baked.
package plugins

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/google/go-github/github"
	"github.com/kyokomi/emoji"
	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/utils"
	"github.com/osteele/liquid"
)

// Site is the site interface that is available to a plugin.
type Site interface {
	AddDocument(pages.Document, bool)
	Config() *config.Config
	TemplateEngine() *liquid.Engine
	Pages() []pages.Page
}

// Plugin describes the hooks that a plugin can override.
type Plugin interface {
	Initialize(Site) error
	ConfigureTemplateEngine(*liquid.Engine) error
	ModifySiteDrop(Site, map[string]interface{}) error
	PostRead(Site) error
	PostRender([]byte) []byte
}

type plugin struct{}

func (p plugin) Initialize(Site) error                             { return nil }
func (p plugin) ConfigureTemplateEngine(*liquid.Engine) error      { return nil }
func (p plugin) ModifySiteDrop(Site, map[string]interface{}) error { return nil }
func (p plugin) PostRead(Site) error                               { return nil }
func (p plugin) PostRender(b []byte) []byte                        { return b }

// Lookup returns a plugin if it has been registered.
func Lookup(name string) (Plugin, bool) {
	p, found := directory[name]
	return p, found
}

// Install installs a plugin from the plugin directory.
func Install(names []string, site Site) {
	for _, name := range names {
		p, found := directory[name]
		if found {
			if err := p.Initialize(site); err != nil {
				panic(err)
			}
		} else {
			fmt.Printf("warning: gojekyll does not emulate the %s plugin.\n", name)
		}
	}
}

var directory = map[string]Plugin{}

// register installs a plugin in the plugin directory.
func register(name string, p Plugin) {
	directory[name] = p
}

func init() {
	register("jemoji", jemojiPlugin{})
	register("jekyll-github-metadata", jekyllGithubMetadataPlugin{})
	register("jekyll-mentions", jekyllMentionsPlugin{})
	register("jekyll-optional-front-matter", jekyllOptionalFrontMatterPlugin{})

	// the following plugins are always active
	// no warning but effect; the server runs in this mode anyway
	register("jekyll-live-reload", plugin{})
	register("jekyll-sass-converter", plugin{})
}

// Some small plugins are below. More involved plugins are in separate files.

// jekyll-jemoji

type jemojiPlugin struct{ plugin }

func (p jemojiPlugin) PostRender(b []byte) []byte {
	return utils.ApplyToHTMLText(b, func(s string) string {
		s = emoji.Sprint(s)
		return s
	})
}

// jekyll-github-metadata

type jekyllGithubMetadataPlugin struct{ plugin }

func (p jekyllGithubMetadataPlugin) ModifySiteDrop(s Site, d map[string]interface{}) error {
	cmd := exec.Command("git", "remote", "-v") // nolint: gas
	cmd.Dir = s.Config().SourceDir()
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	m := regexp.MustCompile(`origin\s+https://github.com/(.+?)/(.+)\.git\b`).FindStringSubmatch(string(out))
	owner, name := m[1], m[2]
	nwo := fmt.Sprintf("%s/%s", owner, name)
	client := github.NewClient(nil)
	var ctx = context.Background()
	repo, _, err := client.Repositories.Get(ctx, owner, name)
	userPage := *repo.Name == fmt.Sprintf("%s.github.com", strings.ToLower(*repo.Owner.Login))
	gh := map[string]interface{}{
		"clone_url":          repo.CloneURL,
		"language":           repo.Language,
		"owner_gravatar_url": repo.Owner.GravatarID,
		"owner_name":         repo.Owner.Login,
		"owner_url":          repo.Owner.URL,
		"project_tagline":    repo.Description,
		"project_title":      repo.Name, // TODO is this right?
		"repository_name":    repo.Name,
		"repository_nwo":     nwo,
		"repository_url":     repo.URL,
		"releases_url":       repo.ReleasesURL,
		"latest_release_url": repo.URL,
		"issues_url":         repo.IssuesURL,
		"show_downloads?":    repo.HasDownloads,
		"repo_clone_url":     repo.GitURL,
		"is_project_page":    !userPage,
		"is_user_page":       userPage,
		// TODO
		// contributors environment public_repositories releases releases_url
		// show_downloads tar_url zip_url url versions wiki_url
	}
	d["github"] = gh
	return err
}

// jekyll-mentions

type jekyllMentionsPlugin struct{ plugin }

var mentionPattern = regexp.MustCompile(`@(\w+)`)

func (p jekyllMentionsPlugin) PostRender(b []byte) []byte {
	return utils.ApplyToHTMLText(b, func(s string) string {
		return mentionPattern.ReplaceAllString(s, `<a href="https://github.com/$1" class="user-mention">@$1</a>`)
	})
}

// jekyll-optional-front-matter

type jekyllOptionalFrontMatterPlugin struct{ plugin }

var requireFrontMatterExclude = []string{
	"README",
	"LICENSE",
	"LICENCE",
	"COPYING",
	"CODE_OF_CONDUCT",
	"CONTRIBUTING",
	"ISSUE_TEMPLATE",
	"PULL_REQUEST_TEMPLATE",
}

func (p jekyllOptionalFrontMatterPlugin) Initialize(s Site) error {
	m := map[string]bool{}
	for _, k := range requireFrontMatterExclude {
		m[k] = true
	}
	s.Config().RequireFrontMatter = false
	s.Config().RequireFrontMatterExclude = m
	return nil
}

// helpers

// func (p plugin) stubbed(name string) {
// 	fmt.Printf("warning: gojekyll does not emulate the %s plugin. Some tags have been stubbed to prevent errors.\n", name)
// }

// func (p plugin) makeUnimplementedTag(pluginName string) liquid.Renderer {
// 	warned := false
// 	return func(ctx render.Context) (string, error) {
// 		if !warned {
// 			fmt.Printf("The %q tag in the %q plugin has not been implemented.\n", ctx.TagName(), pluginName)
// 			warned = true
// 		}
// 		return fmt.Sprintf(`<!-- unimplemented tag: %q -->`, ctx.TagName()), nil
// 	}
// }
