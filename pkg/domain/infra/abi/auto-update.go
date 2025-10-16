//go:build !remote

package abi

import (
	"context"

	"github.com/dmikushin/podman-shared/pkg/autoupdate"
	"github.com/dmikushin/podman-shared/pkg/domain/entities"
)

func (ic *ContainerEngine) AutoUpdate(ctx context.Context, options entities.AutoUpdateOptions) ([]*entities.AutoUpdateReport, []error) {
	return autoupdate.AutoUpdate(ctx, ic.Libpod, options)
}
