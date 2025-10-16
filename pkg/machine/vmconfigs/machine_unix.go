//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd

package vmconfigs

import (
	"github.com/dmikushin/podman-shared/pkg/machine/define"
)

func getPipe(_ string) *define.VMFile {
	return nil
}
