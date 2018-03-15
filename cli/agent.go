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
	Short: buildtime.PROGNAME + ` API server agent runner`,
	Long: fmt.Sprintf(`
%s - Triton Service Groups API

Runs the Triton Service Groups API server as an agent process.

The agent includes reading in the default and environment configuration then
serving the HTTP API. The HTTP server is configurable and can be bound to any
interface or port. Note that the server is configured to use keep alive by
default (Go standard library) in order to follow HTTP/1.1.

Agent will continue to run in the foreground until an interrupt signal has been
received. Sending SIGINT or SIGTERM will drain/shutdown all HTTP connections
gracefully, while performing a SIGKILL will not.

`, buildtime.PROGNAME),

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
