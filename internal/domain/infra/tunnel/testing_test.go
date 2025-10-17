//go:build !remote

package tunnel

import "github.com/dmikushin/podman-shared/v5/internal/domain/entities"

var _ entities.TestingEngine = &TestingEngine{}
