package agent

import (
	"github.com/jackc/pgx"
	"github.com/rs/zerolog/log"
)

func (a *Agent) ensureDBPool() error {
	log.Debug().Msg("agent: connecting to database")

	pool, err := pgx.NewConnPool(a.config.DBPool)
	if err != nil {
		return err
	}
	a.pool = pool

	return nil
}
