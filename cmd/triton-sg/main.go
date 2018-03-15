package main

import (
	"os"

	"github.com/joyent/triton-service-groups/buildtime"
	"github.com/joyent/triton-service-groups/cli"
	"github.com/rs/zerolog/log"
	"github.com/sean-/conswriter"
)

var (
	// variables populated by govvv(1)
	Version    = "dev"
	BuildDate  string
	DocsDate   string
	GitCommit  string
	GitBranch  string
	GitState   string
	GitSummary string
)

func main() {
	exportBuildtimeConsts()

	defer func() {
		p := conswriter.GetTerminal()
		err := p.Wait()
		if err != nil {
			log.Warn().Err(err)
		}
	}()

	if err := cli.Execute(); err != nil {
		log.Error().Err(err).Msg("unable to run")
		os.Exit(1)
	}
}

func exportBuildtimeConsts() {
	buildtime.GitCommit = GitCommit
	buildtime.GitBranch = GitBranch
	buildtime.GitState = GitState
	buildtime.GitSummary = GitSummary
	buildtime.BuildDate = BuildDate
	if DocsDate != "" {
		buildtime.DocsDate = DocsDate
	} else {
		buildtime.DocsDate = BuildDate
	}
	buildtime.Version = Version
}
