package cmd

// Version is the build version.
//
// Make initializes this to the git commit hash.
var Version string

func init() {
	if Version == "" {
		Version = "develop"
	}
}
