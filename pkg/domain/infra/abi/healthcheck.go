//go:build !remote

package abi

import (
	"context"

	"github.com/dmikushin/podman-shared/libpod/define"
	"github.com/dmikushin/podman-shared/pkg/domain/entities"
)

func (ic *ContainerEngine) HealthCheckRun(ctx context.Context, nameOrID string, _ entities.HealthCheckOptions) (*define.HealthCheckResults, error) {
	status, err := ic.Libpod.HealthCheck(ctx, nameOrID)
	if err != nil {
		return nil, err
	}
	report := define.HealthCheckResults{
		Status: status.String(),
	}
	return &report, nil
}
