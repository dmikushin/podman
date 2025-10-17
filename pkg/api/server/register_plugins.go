//go:build !remote

package server

import (
	"github.com/dmikushin/podman-shared/pkg/api/handlers/compat"
	"github.com/gorilla/mux"
)

func (s *APIServer) registerPluginsHandlers(r *mux.Router) error {
	r.Handle(VersionedPath("/plugins"), s.APIHandler(compat.UnsupportedHandler))
	// Added non version path to URI to support docker non versioned paths
	r.Handle("/plugins", s.APIHandler(compat.UnsupportedHandler))
	return nil
}
