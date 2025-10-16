//go:build !remote

package abi

import "github.com/dmikushin/podman-shared/internal/domain/entities"

var _ entities.TestingEngine = &TestingEngine{}
