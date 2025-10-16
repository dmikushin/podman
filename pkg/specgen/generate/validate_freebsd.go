//go:build !remote

package generate

import "github.com/dmikushin/podman-shared/pkg/specgen"

// verifyContainerResources does nothing on freebsd as it has no cgroups
func verifyContainerResources(_ *specgen.SpecGenerator) ([]string, error) {
	return nil, nil
}
