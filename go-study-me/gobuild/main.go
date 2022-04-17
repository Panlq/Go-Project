package main

import (
	"fmt"
	"os"

	"github.com/panlq/gobuild/version"
)

func main() {
	args := os.Args
	if len(args) == 2 && (args[1] == "--version" || args[1] == "-v") {
		fmt.Fprintf(os.Stdout, "client version %#v\n", version.Get())
	}
}
