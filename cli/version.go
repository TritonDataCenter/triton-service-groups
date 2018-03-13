package cli

import (
	"fmt"

	"github.com/joyent/triton-service-groups/buildtime"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: buildtime.PROGNAME + ` version information`,
	Long:  fmt.Sprintf(`Display %s version information`, buildtime.PROGNAME),

	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("%s:\n", buildtime.PROGNAME)
		fmt.Printf("\tversion: %s\n", buildtime.Version)
		fmt.Printf("\tdate: %s\n", buildtime.BuildDate)
		fmt.Printf("\tcommit: %s\n", buildtime.GitCommit)
		fmt.Printf("\tbranch: %s\n", buildtime.GitBranch)
		fmt.Printf("\tstate: %s\n", buildtime.GitState)
		fmt.Printf("\tsummary: %s\n", buildtime.GitSummary)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
