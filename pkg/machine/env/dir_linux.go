package env

import (
	"github.com/dmikushin/podman-shared/pkg/rootless"
	"github.com/dmikushin/podman-shared/pkg/util"
)

func getRuntimeDir() (string, error) {
	if !rootless.IsRootless() {
		return "/run", nil
	}
	return util.GetRootlessRuntimeDir()
}
