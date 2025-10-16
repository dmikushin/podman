//go:build !remote

package libpod

import (
	"net/http"

	"github.com/dmikushin/podman-shared/libpod"
	"github.com/dmikushin/podman-shared/pkg/api/handlers/utils"
	api "github.com/dmikushin/podman-shared/pkg/api/types"
	"github.com/dmikushin/podman-shared/pkg/domain/infra/abi"
)

func GetInfo(w http.ResponseWriter, r *http.Request) {
	runtime := r.Context().Value(api.RuntimeKey).(*libpod.Runtime)
	containerEngine := abi.ContainerEngine{Libpod: runtime}
	info, err := containerEngine.Info(r.Context())
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}
	utils.WriteResponse(w, http.StatusOK, info)
}
