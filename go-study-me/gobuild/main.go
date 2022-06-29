package main

import (
	"log"

	"github.com/panlq/gobuild/version"
)

func main() {
	if err := version.NewCmdVersion().Execute(); err != nil {
		log.Fatal(err)
	}
}
