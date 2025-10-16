package network

import (
	"fmt"

	"github.com/dmikushin/podman-shared/cmd/podman/common"
	"github.com/dmikushin/podman-shared/cmd/podman/registry"
	"github.com/dmikushin/podman-shared/pkg/domain/entities"
	"github.com/spf13/cobra"
	"go.podman.io/common/pkg/completion"
)

var (
	networkUpdateDescription = `Update an existing podman network`
	networkUpdateCommand     = &cobra.Command{
		Use:               "update [options] NETWORK",
		Short:             "Update an existing podman network",
		Long:              networkUpdateDescription,
		RunE:              networkUpdate,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: common.AutocompleteNetworks,
		Example:           `podman network update podman1`,
	}
)

var (
	networkUpdateOptions entities.NetworkUpdateOptions
)

func networkUpdateFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	addDNSServerFlagName := "dns-add"
	flags.StringSliceVar(&networkUpdateOptions.AddDNSServers, addDNSServerFlagName, nil, "add network level nameservers")
	removeDNSServerFlagName := "dns-drop"
	flags.StringSliceVar(&networkUpdateOptions.RemoveDNSServers, removeDNSServerFlagName, nil, "remove network level nameservers")
	_ = cmd.RegisterFlagCompletionFunc(addDNSServerFlagName, completion.AutocompleteNone)
	_ = cmd.RegisterFlagCompletionFunc(removeDNSServerFlagName, completion.AutocompleteNone)
}
func init() {
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Command: networkUpdateCommand,
		Parent:  networkCmd,
	})
	networkUpdateFlags(networkUpdateCommand)
}

func networkUpdate(_ *cobra.Command, args []string) error {
	name := args[0]

	err := registry.ContainerEngine().NetworkUpdate(registry.Context(), name, networkUpdateOptions)
	if err != nil {
		return err
	}
	fmt.Println(name)
	return nil
}
