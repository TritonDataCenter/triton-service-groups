package cli

import (
	"fmt"

	"github.com/joyent/triton-service-groups/buildtime"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: buildtime.PROGNAME + ` version information`,
	Long: fmt.Sprintf(`
%s - Triton Service Groups API

Displays version, SCM, and build information.

`, buildtime.PROGNAME),

	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("%s\n\n", buildtime.PROGNAME)
		fmt.Printf("Version: %s\n", buildtime.Version)
		fmt.Printf("Date: %s\n", buildtime.BuildDate)
		fmt.Printf("Commit: %s\n", buildtime.GitCommit)
		fmt.Printf("Branch: %s\n", buildtime.GitBranch)
		fmt.Printf("State: %s\n", buildtime.GitState)
		fmt.Printf("Summary: %s\n", buildtime.GitSummary)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
