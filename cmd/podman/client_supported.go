//go:build amd64 || arm64

package main

import (
	"github.com/dmikushin/podman-shared/v5/pkg/machine/provider"
)

func getProvider() (string, error) {
	p, err := provider.Get()
	if err != nil {
		return "", err
	}
	return p.VMType().String(), nil
}
