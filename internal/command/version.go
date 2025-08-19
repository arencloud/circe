package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

type VersionCmd struct {
	command *cobra.Command
}

func NewVersionCmd() *VersionCmd {
	c := &VersionCmd{}
	c.command = &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			versionStr := Version
			if Commit != "none" || Date != "unknown" {
				versionStr = fmt.Sprintf("%s (commit %s, date %s)", Version, Commit, Date)
			}
			cmd.Printf("circe version %s\n", versionStr)
		},
	}
	return c
}
