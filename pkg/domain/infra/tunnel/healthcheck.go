package tunnel

import (
	"context"

	"github.com/dmikushin/podman-shared/v5/libpod/define"
	"github.com/dmikushin/podman-shared/v5/pkg/bindings/containers"
	"github.com/dmikushin/podman-shared/v5/pkg/domain/entities"
)

func (ic *ContainerEngine) HealthCheckRun(_ context.Context, nameOrID string, _ entities.HealthCheckOptions) (*define.HealthCheckResults, error) {
	return containers.RunHealthCheck(ic.ClientCtx, nameOrID, nil)
}
