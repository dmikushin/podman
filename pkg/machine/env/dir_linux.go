package env

import (
	"github.com/dmikushin/podman-shared/v5/pkg/rootless"
	"github.com/dmikushin/podman-shared/v5/pkg/util"
)

func getRuntimeDir() (string, error) {
	if !rootless.IsRootless() {
		return "/run", nil
	}
	return util.GetRootlessRuntimeDir()
}
