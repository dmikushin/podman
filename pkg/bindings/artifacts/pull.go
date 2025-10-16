package artifacts

import (
	"context"
	"net/http"

	"github.com/dmikushin/podman-shared/pkg/auth"
	"github.com/dmikushin/podman-shared/pkg/bindings"
	"github.com/dmikushin/podman-shared/pkg/domain/entities"
	imageTypes "go.podman.io/image/v5/types"
)

func Pull(ctx context.Context, name string, options *PullOptions) (*entities.ArtifactPullReport, error) {
	if options == nil {
		options = new(PullOptions)
	}

	conn, err := bindings.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	params, err := options.ToParams()
	if err != nil {
		return nil, err
	}
	params.Set("name", name)

	header, err := auth.MakeXRegistryAuthHeader(
		&imageTypes.SystemContext{
			AuthFilePath: options.GetAuthfile(),
		},
		options.GetUsername(),
		options.GetPassword(),
	)
	if err != nil {
		return nil, err
	}

	response, err := conn.DoRequest(ctx, nil, http.MethodPost, "/artifacts/pull", params, header)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var report entities.ArtifactPullReport
	if err := response.Process(&report); err != nil {
		return nil, err
	}

	return &report, nil
}
