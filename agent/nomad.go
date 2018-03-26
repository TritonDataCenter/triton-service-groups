package agent

import (
	"fmt"

	nomad "github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog/log"
)

func (a *Agent) ensureNomadClient() error {
	log.Debug().Msg("agent: connecting to job scheduler")

	nomadCfg := nomad.DefaultConfig()
	scheme := "http"

	if a.config.Nomad.TLSConfig != nil {
		nomadCfg.TLSConfig = a.config.Nomad.TLSConfig
		scheme = "https"
	}

	nomadCfg.Address = fmt.Sprintf("%s://%s:%d",
		scheme, a.config.Nomad.Addr, a.config.Nomad.Port)

	c, err := nomad.NewClient(nomadCfg)
	if err != nil {
		return err
	}
	a.nomad = c

	return nil
}
