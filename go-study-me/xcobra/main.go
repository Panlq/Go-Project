package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var diff bool

type option struct {
	diff bool
}

var defaultOption = option{}

func main() {
	cmd := &cobra.Command{
		Use: "no use",
		Run: func(c *cobra.Command, args []string) {
		},
	}

	pflag.BoolVar(&defaultOption.diff, "diff", false, "no working")
	// Ignore errors; CommandLine is set for ExitOnError.
	//nolint
	pflag.Parse()
	//nolint
	viper.BindPFlags(pflag.CommandLine)

	if err := cmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
