package common

import (
	"github.com/dmikushin/podman-shared/v5/cmd/podman/registry"
)

var (
	// Pull in configured json library
	json = registry.JSONLibrary()
)
