package agent

import (
	"context"
	"os"

	nomad "github.com/hashicorp/nomad/api"
	"github.com/jackc/pgx"
	"github.com/joyent/triton-service-groups/buildtime"
	"github.com/joyent/triton-service-groups/config"
	"github.com/joyent/triton-service-groups/server"
	"github.com/rs/zerolog/log"
)

type Agent struct {
	signalCh    chan os.Signal
	shutdownCtx context.Context
	shutdown    func()
	config      *config.Config
	pool        *pgx.ConnPool
	nomad       *nomad.Client
}

func New(cfg *config.Config) *Agent {
	log.Debug().Msg("agent: initializing agent")

	return &Agent{
		config: cfg,
	}
}

func (a *Agent) Run(ctx context.Context) (err error) {
	log.Debug().Msgf("agent: running %s agent", buildtime.PROGNAME)

	a.shutdownCtx, a.shutdown = context.WithCancel(ctx)

	go a.handleSignals()

	if err = a.ensureDBPool(); err != nil {
		return err
	}

	if err = a.ensureNomadClient(); err != nil {
		return err
	}

	srv := server.New(a.config.HTTPServer, a.pool, a.nomad)
	srv.Start()

	for {
		<-a.shutdownCtx.Done()
		err := srv.Stop(a.shutdownCtx)
		if err != nil {
			log.Warn().Err(err)
		}
		return nil
	}
}

func (a *Agent) Stop() {
	log.Info().Msgf("agent: shutting down %s agent", buildtime.PROGNAME)

	a.stopSignalCh()
	a.pool.Close()
	a.shutdown()
}
