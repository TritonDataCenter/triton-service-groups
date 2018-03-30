package testutils

import (
	nomad "github.com/hashicorp/nomad/api"
)

// NewNomadClient creates and returns a new Nomad client used for testing.
//
// TODO(justinwr): configure an HTTP client we can use to stub out requests to
// Nomad for unit testing
func NewNomadClient() (*nomad.Client, error) {
	nomadCfg := nomad.DefaultConfig()
	return nomad.NewClient(nomadCfg)
}
