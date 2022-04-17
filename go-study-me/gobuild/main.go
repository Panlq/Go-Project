package main

import (
	"fmt"
	"os"
)

var (
	buildtime = ""
	branch    = ""
	commit    = ""
	goversion = ""
)

func main() {
	args := os.Args
	if len(args) == 2 && (args[1] == "--version" || args[1] == "-v") {
		fmt.Printf("Build Time: %s\n", buildtime)
		fmt.Printf("GitCommit: %s:%s\n", branch, commit)
		fmt.Printf("GO Version: %s\n", goversion)
		// fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	}
}
