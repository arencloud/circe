package command

func InitialiseRootCmd() *RootCommand {
	rootCommand := NewRootCommand()
	networkPolicyCommand := NewNetworkPolicyCmd()
	versionCmd := NewVersionCmd()
	rootCommand.Command.AddCommand(
		networkPolicyCommand.commnad,
		versionCmd.command,
	)
	return rootCommand
}
