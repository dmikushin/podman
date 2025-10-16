//go:build !remote

package tunnel

import "github.com/dmikushin/podman-shared/internal/domain/entities"

var _ entities.TestingEngine = &TestingEngine{}
