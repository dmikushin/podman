package vmconfigs

import (
	"github.com/dmikushin/podman-shared/pkg/machine/define"
	"github.com/dmikushin/podman-shared/pkg/machine/env"
)

func getPipe(name string) *define.VMFile {
	pipeName := env.WithPodmanPrefix(name)
	return &define.VMFile{Path: `\\.\pipe\` + pipeName}
}
