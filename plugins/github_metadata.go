package plugins

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/google/go-github/github"
	"github.com/osteele/gojekyll/config"
	"golang.org/x/oauth2"
)

func init() {
	register("jekyll-github-metadata", jekyllGithubMetadataPlugin{})
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
