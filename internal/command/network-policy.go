package command

import "github.com/spf13/cobra"

type NetworkPolicyCmd struct {
	commnad *cobra.Command
}

func NewNetworkPolicyCmd() *NetworkPolicyCmd {
	c := &NetworkPolicyCmd{
		commnad: &cobra.Command{
			Use:   "network-policy",
			Short: "network policy menu",
		},
	}
	egressGen := NewEgressCommand()
	ingressGen := NewIngressCommand()
	c.commnad.AddCommand(
		egressGen.command,
		ingressGen.command,
	)
	return c
}
