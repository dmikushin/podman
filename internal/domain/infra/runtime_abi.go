//go:build !remote

package infra

import (
	"context"
	"fmt"

	ientities "github.com/dmikushin/podman-shared/internal/domain/entities"
	"github.com/dmikushin/podman-shared/internal/domain/infra/tunnel"
	"github.com/dmikushin/podman-shared/pkg/bindings"
	"github.com/dmikushin/podman-shared/pkg/domain/entities"
)

// NewTestingEngine factory provides a libpod runtime for testing-specific operations
func NewTestingEngine(facts *entities.PodmanConfig) (ientities.TestingEngine, error) {
	switch facts.EngineMode {
	case entities.ABIMode:
		r, err := NewLibpodTestingRuntime(facts.FlagSet, facts)
		return r, err
	case entities.TunnelMode:
		ctx, err := bindings.NewConnectionWithOptions(context.Background(), bindings.Options{
			URI:         facts.URI,
			Identity:    facts.Identity,
			TLSCertFile: facts.TLSCertFile,
			TLSKeyFile:  facts.TLSKeyFile,
			TLSCAFile:   facts.TLSCAFile,
			Machine:     facts.MachineMode,
		})
		return &tunnel.TestingEngine{ClientCtx: ctx}, err
	}
	return nil, fmt.Errorf("runtime mode '%v' is not supported", facts.EngineMode)
}
