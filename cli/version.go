package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/matoous/linkfix/internal/version"
)

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of linkfix",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("linkfix %s (%s)\n", version.GitTag, version.GitCommit)
		},
	}
}
