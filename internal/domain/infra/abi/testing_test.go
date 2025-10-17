//go:build !remote

package abi

import "github.com/dmikushin/podman-shared/v5/internal/domain/entities"

var _ entities.TestingEngine = &TestingEngine{}
