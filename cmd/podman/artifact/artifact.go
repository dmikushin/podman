package artifact

import (
	"github.com/dmikushin/podman-shared/cmd/podman/registry"
	"github.com/dmikushin/podman-shared/cmd/podman/validate"
	"github.com/spf13/cobra"
)

// Command: podman _artifact_
var artifactCmd = &cobra.Command{
	Use:   "artifact",
	Short: "Manage OCI artifacts",
	Long:  "Manage OCI artifacts",
	RunE:  validate.SubCommandExists,
}

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Command: artifactCmd,
	})
}
