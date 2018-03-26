// File is used to host the HTTP Server abstraction, start and stop.

package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	ghandlers "github.com/gorilla/handlers"
	nomad "github.com/hashicorp/nomad/api"
	"github.com/jackc/pgx"
	"github.com/joyent/triton-service-groups/config"
	"github.com/joyent/triton-service-groups/server/handlers"
	"github.com/joyent/triton-service-groups/server/router"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	ErrConfig = fmt.Errorf("server: valid config required")
)

type HTTPServer struct {
	Addr string
	Bind string
	Port uint16

	logger zerolog.Logger
	pool   *pgx.ConnPool
	nomad  *nomad.Client

	http.Server
}

func New(cfg config.HTTPServer, pool *pgx.ConnPool, nomad *nomad.Client) *HTTPServer {
	log.Debug().Msg("http: creating new HTTP server")
	addr := fmt.Sprintf("%s:%d", cfg.Bind, cfg.Port)

	return &HTTPServer{
		Addr:   addr,
		Bind:   cfg.Bind,
		Port:   cfg.Port,
		logger: cfg.Logger,
		pool:   pool,
		nomad:  nomad,
	}
}

func (srv *HTTPServer) Start() {
	log.Debug().Msg("http: starting up HTTP server")

	srv.setup()
}

func (srv *HTTPServer) setup() {
	log.Debug().Msg("http: mounting routes as endpoints")

	router := router.WithRoutes(RoutingTable)
	authHandler := handlers.AuthHandler(srv.pool, router)
	contextHandler := handlers.ContextHandler(srv.pool, srv.nomad, authHandler)
	srv.Handler = ghandlers.LoggingHandler(srv.logger, contextHandler)

	ln := srv.listenWithRetry()

	go func() {
		log.Info().Msgf("http: started serving at %q", srv.Addr)
		err := srv.Serve(ln)
		if err != nil {
			log.Warn().Err(err)
		}
	}()
}

// listenWithRetry attempts to listen on our socket, failing after 10 seconds.
func (srv *HTTPServer) listenWithRetry() net.Listener {
	var (
		err error
		ln  net.Listener
	)

	for i := 0; i < 10; i++ {
		ln, err = net.Listen("tcp", srv.Addr)
		if err == nil {
			log.Debug().
				Str("http-bind", srv.Bind).
				Int("http-port", int(srv.Port)).
				Msgf("http: server listening at %q", srv.Addr)
			return ln
		}

		time.Sleep(time.Second)
	}
	return nil
}

// Stop handles gracefully shutting down the server, finally forcing shutdown
// after 3 seconds.
func (srv *HTTPServer) Stop(ctx context.Context) error {
	log.Debug().Msg("http: gracefully shutting down HTTP server")

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
