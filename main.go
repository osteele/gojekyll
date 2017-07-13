package main

import (
	"fmt"
	"os"

	"github.com/osteele/gojekyll/cmd"
)

func main() {
	err := cmd.ParseAndRun(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
