package secrets

import (
	"github.com/dmikushin/podman-shared/v5/cmd/podman/registry"
	"github.com/dmikushin/podman-shared/v5/cmd/podman/validate"
	"github.com/spf13/cobra"
)

var (
	// Command: podman _secret_
	secretCmd = &cobra.Command{
		Use:   "secret",
		Short: "Manage secrets",
		Long:  "Manage secrets",
		RunE:  validate.SubCommandExists,
	}
)

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Command: secretCmd,
	})
}
