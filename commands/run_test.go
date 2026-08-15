package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseAndRun(t *testing.T) {
	err := ParseAndRun([]string{"build", "-s", "testdata/site", "-q"})
	require.NoError(t, err)
}

// dispatchBySwitch is a stand-in for the original switch-based dispatcher,
// used only to compare dispatch performance against the map-based table.
func dispatchBySwitch(cmd string) (needsSite bool, ok bool) {
	switch cmd {
	case "benchmark", "plugins", "version":
		return false, true
	case "build", "clean", "render", "routes", "serve", "variables":
		return true, true
	}
	return false, false
}

var benchCommands = []string{
	benchmark.FullCommand(),
	pluginsApp.FullCommand(),
	versionCmd.FullCommand(),
	build.FullCommand(),
	clean.FullCommand(),
	render.FullCommand(),
	routes.FullCommand(),
	serve.FullCommand(),
	variables.FullCommand(),
}

func BenchmarkCommandDispatchMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, cmd := range benchCommands {
			entry, ok := commandTable[cmd]
			if !ok {
				b.Fatalf("missing %s", cmd)
			}
			_ = entry.needsSite
		}
	}
}

func BenchmarkCommandDispatchSwitch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, cmd := range benchCommands {
			needsSite, ok := dispatchBySwitch(cmd)
			if !ok {
				b.Fatalf("missing %s", cmd)
			}
			_ = needsSite
		}
	}
}
