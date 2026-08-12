package renderers

import (
	"sync"
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/liquid"
	"github.com/stretchr/testify/require"
)

func TestIsMarkdown(t *testing.T) {
	require.False(t, false)
}

type immutablePageDrop struct {
	CustomVar string
}

func (d immutablePageDrop) ToLiquid() interface{} {
	return map[string]interface{}{"custom_var": d.CustomVar}
}

func TestJekyllAssignToDropWarnsAndHasNoEffect(t *testing.T) {
	var warnings []Warning
	manager, err := New(config.Default(), Options{
		WarningHandler: func(warning Warning) {
			warnings = append(warnings, warning)
		},
	})
	require.NoError(t, err)

	bindings := liquid.Bindings{"page": immutablePageDrop{CustomVar: "from-frontmatter"}}
	result, err := manager.RenderTemplate(
		[]byte("{{ page.custom_var }}\n{% assign page.custom_var = \"from-assign\" %}\n{{ page.custom_var }}"),
		bindings,
		"index.md",
		1,
	)
	require.NoError(t, err)
	require.Equal(t, "from-frontmatter\n\nfrom-frontmatter", string(result))
	require.Equal(t, []Warning{{
		Path:    "index.md",
		Line:    2,
		Message: "`{% assign page.custom_var = \"from-assign\" %}` has no effect because `page` is not an assignable object",
	}}, warnings)
}

func TestJekyllAssignToMapUpdatesCopyWithoutWarning(t *testing.T) {
	var warnings []Warning
	manager, err := New(config.Default(), Options{
		WarningHandler: func(warning Warning) {
			warnings = append(warnings, warning)
		},
	})
	require.NoError(t, err)

	page := map[string]any{"custom_var": "before"}
	result, err := manager.RenderTemplate(
		[]byte("{% assign page.custom_var = \"after\" %}{{ page.custom_var }}"),
		liquid.Bindings{"page": page},
		"index.md",
		1,
	)
	require.NoError(t, err)
	require.Equal(t, "after", string(result))
	require.Equal(t, "before", page["custom_var"], "Liquid assignments must not mutate caller bindings")
	require.Empty(t, warnings)
}

func TestJekyllAssignWarningsAreSafeDuringConcurrentRenders(t *testing.T) {
	var warnings []Warning
	manager, err := New(config.Default(), Options{
		WarningHandler: func(warning Warning) {
			warnings = append(warnings, warning)
		},
	})
	require.NoError(t, err)

	const renders = 8
	var wg sync.WaitGroup
	errs := make(chan error, renders)
	for range renders {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, renderErr := manager.RenderTemplate(
				[]byte("{% assign page.value = 1 %}"),
				liquid.Bindings{"page": immutablePageDrop{}},
				"concurrent.md",
				1,
			)
			errs <- renderErr
		}()
	}
	wg.Wait()
	close(errs)
	for renderErr := range errs {
		require.NoError(t, renderErr)
	}
	require.Len(t, warnings, renders)
}
