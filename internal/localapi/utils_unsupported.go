//go:build !amd64 && !arm64

package localapi

import (
	"context"

	"github.com/dmikushin/podman-shared/pkg/domain/entities"
	"github.com/sirupsen/logrus"
)

func CheckPathOnRunningMachine(ctx context.Context, path string) (*LocalAPIMap, bool) {
	logrus.Debug("CheckPathOnRunningMachine is not supported")
	return nil, false
}

func CheckIfImageBuildPathsOnRunningMachine(ctx context.Context, containerFiles []string, options entities.BuildOptions) ([]string, entities.BuildOptions, bool) {
	logrus.Debug("CheckIfImageBuildPathsOnRunningMachine is not supported")
	return nil, options, false
}

func IsHyperVProvider(ctx context.Context) (bool, error) {
	logrus.Debug("IsHyperVProvider is not supported")
	return false, nil
}

func ValidatePathForLocalAPI(path string) error {
	logrus.Debug("ValidatePathForLocalAPI is not supported")
	return nil
}
