package cli

import (
	"context"
	"fmt"

	"github.com/joyent/triton-service-groups/agent"
	"github.com/joyent/triton-service-groups/buildtime"
	"github.com/joyent/triton-service-groups/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: buildtime.PROGNAME + ` agent process`,
	Long:  fmt.Sprintf(`Starts the %s agent and runs until an interrupt is received.`, buildtime.PROGNAME),

	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info().Msgf("agent: starting %s agent", buildtime.PROGNAME)

		cfg, err := config.NewDefault()
		if err != nil {
			return err
		}

		a := agent.New(cfg)
		if err = a.Run(context.Background()); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(agentCmd)
}
