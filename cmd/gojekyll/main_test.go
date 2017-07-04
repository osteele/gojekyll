package main

import (
	"testing"
)

func TestBuild(t *testing.T) {
	parseAndRun([]string{"build", "-s", "../../testdata/example", "-q"})
}
