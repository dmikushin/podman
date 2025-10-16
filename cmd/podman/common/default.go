package common

import (
	"github.com/dmikushin/podman-shared/cmd/podman/registry"
)

var (
	// Pull in configured json library
	json = registry.JSONLibrary()
)
