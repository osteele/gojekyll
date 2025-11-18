package plugins

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/osteele/gojekyll/tags"
	"github.com/osteele/liquid"
	"github.com/osteele/liquid/render"
	liquidtags "github.com/osteele/liquid/tags"
)

func init() {
	register("jekyll-gist", jekyllGistPlugin{})
}

type jekyllGistPlugin struct{ plugin }

func (p jekyllGistPlugin) ConfigureTemplateEngine(e *liquid.Engine) error {
	e.RegisterTag("gist", gistTag)
	return nil
}

func gistTag(ctx render.Context) (string, error) {
	argsline, err := ctx.ExpandTagArg()
	if err != nil {
		return "", err
	}
	args, err := tags.ParseArgs(argsline)
	if err != nil {
		return "", err
	}
	if len(args.Args) < 1 {
		return "", fmt.Errorf("gist tag: missing argument")
	}

	gistID := args.Args[0]
	var filename string
	if len(args.Args) >= 2 {
		filename = args.Args[1]
	}

	// Generate script tag
	scriptURL := fmt.Sprintf("https://gist.github.com/%s.js", gistID)
	if filename != "" {
		scriptURL += fmt.Sprintf("?file=%s", filename)
	}
	output := fmt.Sprintf(`<script src="%s"> </script>`, scriptURL)

	// Check if noscript is enabled in config
	noscriptEnabled := false
	if site := ctx.Get("site"); site != nil {
		if siteDrop, ok := liquid.FromDrop(site).(liquidtags.IterationKeyedMap); ok {
			if gistConfig, ok := siteDrop["gist"].(map[string]interface{}); ok {
				if noscript, ok := gistConfig["noscript"].(bool); ok {
					noscriptEnabled = noscript
				}
			}
		}
	}

	// If noscript is enabled, fetch and include the raw gist content
	if noscriptEnabled {
		code, err := fetchGistContent(gistID, filename)
		if err == nil && code != "" {
			escapedCode := html.EscapeString(code)
			output += fmt.Sprintf("<noscript><pre>%s</pre></noscript>", escapedCode)
		}
		// Silently ignore fetch errors to match jekyll-gist behavior
	}

	return output, nil
}

// fetchGistContent retrieves the raw content of a gist from GitHub
func fetchGistContent(gistID, filename string) (string, error) {
	// Build the raw content URL
	// Format: https://gist.githubusercontent.com/{user}/{id}/raw/{file}
	// If no filename is specified, GitHub returns the first file
	url := fmt.Sprintf("https://gist.githubusercontent.com/%s/raw", gistID)
	if filename != "" {
		url += "/" + filename
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch gist: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(body)), nil
}
