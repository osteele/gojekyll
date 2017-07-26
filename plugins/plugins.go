// Package plugins holds emulated Jekyll plugins.
//
// Unlike Jekyll, these are baked into the executable -- both because package "plugin'
// works only on Linux (as of 2017.07); and because the gojekyll implementation is immature and any possible interfaces
// are far from baked.
package plugins

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
	"github.com/kyokomi/emoji"
	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/utils"
	"github.com/osteele/liquid"
)

// Plugin describes the hooks that a plugin can override.
type Plugin interface {
	Initialize(Site) error
	ConfigureTemplateEngine(*liquid.Engine) error
	ModifySiteDrop(Site, map[string]interface{}) error
	PostRead(Site) error
	PostRender([]byte) []byte
}

// Site is the site interface that is available to a plugin.
type Site interface {
	AddDocument(pages.Document, bool)
	Config() *config.Config
	TemplateEngine() *liquid.Engine
	Pages() []pages.Page
}

// Lookup returns a plugin if it has been registered.
func Lookup(name string) (Plugin, bool) {
	p, found := directory[name]
	return p, found
}

// Install installs a registered plugin.
func Install(names []string, site Site) {
	for _, name := range names {
		if p, found := directory[name]; found {
			if err := p.Initialize(site); err != nil {
				panic(err)
			}
		} else {
			fmt.Printf("warning: gojekyll does not emulate the %s plugin.\n", name)
		}
	}
}

// Embed plugin to implement defaults implementations of the Plugin interface.
//
// This is internal until better baked.
type plugin struct{}

func (p plugin) Initialize(Site) error                             { return nil }
func (p plugin) ConfigureTemplateEngine(*liquid.Engine) error      { return nil }
func (p plugin) ModifySiteDrop(Site, map[string]interface{}) error { return nil }
func (p plugin) PostRead(Site) error                               { return nil }
func (p plugin) PostRender(b []byte) []byte                        { return b }

var directory = map[string]Plugin{}

// register installs a plugin in the plugin directory.
//
// This is internal until better baked.
func register(name string, p Plugin) {
	directory[name] = p
}

// Add the built-in plugins defined in this file.
// More extensive plugins are defined and registered in own files.
func init() {
	register("jemoji", jemojiPlugin{})
	register("jekyll-github-metadata", jekyllGithubMetadataPlugin{})
	register("jekyll-mentions", jekyllMentionsPlugin{})
	register("jekyll-optional-front-matter", jekyllOptionalFrontMatterPlugin{})

	// Gojekyll behaves as though the following plugins are always loaded.
	// Define them here so we don't see warnings that they aren't defined.
	register("jekyll-live-reload", plugin{})
	register("jekyll-sass-converter", plugin{})
}

// Some small plugins are below. More involved plugins are in separate files.

// jemojiPlugin emulates the jekyll-jemoji plugin.
type jemojiPlugin struct{ plugin }

func (p jemojiPlugin) PostRender(b []byte) []byte {
	return utils.ApplyToHTMLText(b, func(s string) string {
		s = emoji.Sprint(s)
		return s
	})
}

// jekyllGithubMetadataPlugin emulates the jekyll-github-metadata plugin.
type jekyllGithubMetadataPlugin struct{ plugin }

func (p jekyllGithubMetadataPlugin) ModifySiteDrop(s Site, d map[string]interface{}) error {
	nwo, err := getCurrentRepo(s.Config())
	if err != nil {
		return err
	}
	nameOwner := strings.SplitN(nwo, "/", 2)

	ctx := context.Background()
	var ts oauth2.TokenSource
	if tok := os.Getenv("JEKYLL_GITHUB_TOKEN"); tok != "" {
		ts = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: tok})
	} else if tok := os.Getenv("OCTOKIT_ACCESS_TOKEN"); tok != "" {
		ts = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: tok})
	}
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	repo, _, err := client.Repositories.Get(ctx, nameOwner[0], nameOwner[1])
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

		// TODO: build_revision: `git rev-parse HEAD`
		"api_url":        "https://api.github.com",
		"environment":    "development",
		"help_url":       "https://help.github.com",
		"hostname":       "https://github.com",
		"pages_hostname": "github.io",
		// TODO
		// contributors public_repositories show_downloads releases versions
		// url tar_url zip_url wiki_url
	}
	for key, envName := range githubPagesEnvVars {
		if s := os.Getenv(envName); s != "" {
			gh[key] = s
		}
	}
	d["github"] = gh
	return err
}

var githubPagesEnvVars = map[string]string{
	"build_revision": "JEKYLL_BUILD_REVISION",
	"environment":    "PAGES_ENV",
	"api_url":        "PAGES_API_URL",
	"help_url":       "PAGES_HELP_URL",
	"hostname":       "PAGES_GITHUB_HOSTNAME",
	"pages_hostname": "PAGES_PAGES_HOSTNAME",
}

func getCurrentRepo(c *config.Config) (nwo string, err error) {
	nwo = os.Getenv("PAGES_REPO_NWO")
	if nwo != "" {
		return
	}
	if s, ok := c.Variables["repository"]; ok {
		if s, ok := s.(string); ok {
			return s, nil
		}
	}
	cmd := exec.Command("git", "remote", "-v") // nolint: gas
	cmd.Dir = c.SourceDir()
	out, err := cmd.Output()
	if err != nil {
		return
	}
	m := regexp.MustCompile(`origin\s+https://github.com/(.+?)/(.+)\.git\b`).FindStringSubmatch(string(out))
	owner, name := m[1], m[2]
	nwo = fmt.Sprintf("%s/%s", owner, name)
	return
}

// jekyllMentionsPlugin emulates the jekyll-mentions plugin.
type jekyllMentionsPlugin struct{ plugin }

var mentionPattern = regexp.MustCompile(`@(\w+)`)

func (p jekyllMentionsPlugin) PostRender(b []byte) []byte {
	return utils.ApplyToHTMLText(b, func(s string) string {
		return mentionPattern.ReplaceAllString(s, `<a href="https://github.com/$1" class="user-mention">@$1</a>`)
	})
}

// jekyllOptionalFrontMatterPlugin emulates the jekyll-optional-front-matter plugin.
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
