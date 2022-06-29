package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmdVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the application version information",
		Long:  "Print the application version information for the current context.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%#v\n", Get())
		},
	}
}
