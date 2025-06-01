package plugins

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/google/go-github/github"
	"github.com/osteele/gojekyll/config"
	"github.com/osteele/liquid"
	"golang.org/x/oauth2"
)

func init() {
	register("jekyll-github-metadata", jekyllGithubMetadataPlugin{})
}

// jekyllGithubMetadataPlugin emulates the jekyll-github-metadata plugin.
type jekyllGithubMetadataPlugin struct{ plugin }

func (p jekyllGithubMetadataPlugin) ModifySiteDrop(s Site, d map[string]interface{}) error {
	var (
		cfg        = s.Config()
		isUserPage = false
		ref        = "master"
	)
	nwo, err := getCurrentRepo(cfg)
	if err != nil {
		return err
	}
	repo, err := getGitHubRepo(nwo)
	if err != nil {
		return err
	}
	if *repo.Name == fmt.Sprintf("%s.github.com", strings.ToLower(*repo.Owner.Login)) {
		isUserPage = true
		ref = "gh-pages"
	}
	gh := map[string]interface{}{
		"build_revision":     getBuildRevision(cfg.SourceDir()),
		"clone_url":          repo.CloneURL,
		"is_project_page":    !isUserPage,
		"is_user_page":       isUserPage,
		"issues_url":         repo.IssuesURL,
		"language":           repo.Language,
		"latest_release_url": repo.URL,
		"owner_gravatar_url": repo.Owner.GravatarID,
		"owner_name":         repo.Owner.Login,
		"owner_url":          repo.Owner.URL,
		"project_tagline":    repo.Description,
		"project_title":      repo.Name, // TODO is this right?
		"releases_url":       repo.ReleasesURL,
		"repo_clone_url":     repo.GitURL,
		"repository_name":    repo.Name,
		"repository_nwo":     nwo,
		"repository_url":     repo.URL,
		"show_downloads?":    repo.HasDownloads,
		"url":                repo.URL,
		"tar_url":            repoArchiveURL(repo, "tarball", ref),
		"zip_url":            repoArchiveURL(repo, "zipball", ref),

		// TODO
		// contributors public_repositories show_downloads releases versions
		// wiki_url

		// These may be replaced by environment variable values
		"api_url":        "https://api.github.com",
		"environment":    "development",
		"help_url":       "https://help.github.com",
		"hostname":       "https://github.com",
		"pages_hostname": "github.io",
	}
	for key, envName := range githubPagesEnvVars {
		if s := os.Getenv(envName); s != "" {
			gh[key] = s
		}
	}
	d["github"] = liquid.IterationKeyedMap(gh)
	return err
}

func getGitHubRepo(nwo string) (*github.Repository, error) {
	ctx := context.Background()
	var ts oauth2.TokenSource
	if tok := os.Getenv("JEKYLL_GITHUB_TOKEN"); tok != "" {
		ts = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: tok})
	} else if tok := os.Getenv("GITHUB_TOKEN"); tok != "" {
		ts = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: tok})
	} else if tok := os.Getenv("OCTOKIT_ACCESS_TOKEN"); tok != "" {
		ts = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: tok})
	}
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	nameAndOwner := strings.SplitN(nwo, "/", 2)
	repo, _, err := client.Repositories.Get(ctx, nameAndOwner[0], nameAndOwner[1])
	return repo, err
}

func repoArchiveURL(repo *github.Repository, format, ref string) string {
	u := strings.Replace(*repo.ArchiveURL, "{archive_format}", format, 1)
	u = strings.Replace(u, "{/ref}", "/"+ref, 1)
	return u
}

// A map of site.github key -> environment variable name
var githubPagesEnvVars = map[string]string{
	"api_url":        "PAGES_API_URL",
	"build_revision": "JEKYLL_BUILD_REVISION",
	"environment":    "PAGES_ENV",
	"help_url":       "PAGES_HELP_URL",
	"hostname":       "PAGES_GITHUB_HOSTNAME",
	"pages_hostname": "PAGES_PAGES_HOSTNAME",
}

var githubURLNWOMatcher = regexp.MustCompile(`origin\s+https://github.com/(.+?/.+)\.git\b`).FindSubmatch

func getCurrentRepo(c *config.Config) (string, error) {
	if nwo := os.Getenv("PAGES_REPO_NWO"); nwo != "" {
		return nwo, nil
	}
	if s, ok := c.String("repository"); ok {
		return s, nil
	}
	cmd := exec.Command("git", "remote", "-v")
	cmd.Dir = c.SourceDir()
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	if m := githubURLNWOMatcher(out); m != nil {
		return string(m[1]), nil
	}
	return "", fmt.Errorf("jekyll-github-metadata failed to find current repository")
}

func getBuildRevision(dir string) string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(bytes.TrimSpace(out))
}
