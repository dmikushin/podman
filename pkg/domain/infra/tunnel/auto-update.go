package tunnel

import (
	"context"
	"errors"

	"github.com/dmikushin/podman-shared/pkg/domain/entities"
)

func (ic *ContainerEngine) AutoUpdate(_ context.Context, _ entities.AutoUpdateOptions) ([]*entities.AutoUpdateReport, []error) {
	return nil, []error{errors.New("not implemented")}
}
