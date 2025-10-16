//go:build (amd64 || arm64) && !windows

package main

import (
	"errors"

	"github.com/dmikushin/podman-shared/pkg/machine/define"
)

func extractConnectionString(podmanSocket *define.VMFile, _ *define.VMFile) (string, error) {
	if podmanSocket == nil {
		return "", errors.New("socket of machine is not set")
	}
	return "unix://" + podmanSocket.Path, nil
}
