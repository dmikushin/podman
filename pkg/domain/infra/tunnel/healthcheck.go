package tunnel

import (
	"context"

	"github.com/dmikushin/podman-shared/libpod/define"
	"github.com/dmikushin/podman-shared/pkg/bindings/containers"
	"github.com/dmikushin/podman-shared/pkg/domain/entities"
)

func (ic *ContainerEngine) HealthCheckRun(_ context.Context, nameOrID string, _ entities.HealthCheckOptions) (*define.HealthCheckResults, error) {
	return containers.RunHealthCheck(ic.ClientCtx, nameOrID, nil)
}
