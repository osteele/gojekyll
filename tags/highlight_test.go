package tags

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/osteele/gojekyll/cache"
	"github.com/osteele/gojekyll/config"
	"github.com/osteele/liquid"
	"github.com/stretchr/testify/require"
)

var highlightTagTests = []struct{ in, out string }{
	{`{% highlight ruby %}
	def foo
	  puts 'foo'
	end
	{% endhighlight %}`, "chroma"},
	{`{% highlight ruby linenos %}
	def foo
	  puts 'foo'
	end
	{% endhighlight %}`, "lntable"},
}

func TestHighlightTag(t *testing.T) {
	cache.Disable()
	engine := liquid.NewEngine()
	cfg := config.Default()
	AddJekyllTags(engine, &cfg, []string{}, func(string) (string, bool) { return "", false })

	for i, test := range highlightTagTests {
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			s, err := engine.ParseAndRenderString(test.in, liquid.Bindings{})
			require.NoError(t, err)
			re := regexp.MustCompile(fmt.Sprintf(`class="%s"`, test.out))
			require.True(t, re.MatchString(s))
		})
	}
}
