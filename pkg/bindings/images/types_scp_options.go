package images

import (
	"net/url"

	"github.com/dmikushin/podman-shared/pkg/bindings/internal/util"
)

// ToParams formats struct fields to be passed to API service
func (o *ScpOptions) ToParams() (url.Values, error) {
	return util.ToParams(o)
}
