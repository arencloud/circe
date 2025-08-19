package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

type RootCommand struct {
	Command *cobra.Command
}

// These variables can be overridden at build time using -ldflags, e.g.:
//
//	go build -ldflags "-X 'circe/internal/command.Version=v0.2.0' -X 'circe/internal/command.Commit=abcdef1' -X 'circe/internal/command.Date=2025-08-19'" ./cmd/main
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func NewRootCommand() *RootCommand {
	versionStr := Version
	// Enrich version string with commit/date if provided
	if Commit != "none" || Date != "unknown" {
		versionStr = fmt.Sprintf("%s (commit %s, date %s)", Version, Commit, Date)
	}

	c := &RootCommand{
		Command: &cobra.Command{
			Use:     "circe",
			Short:   "Circe: CLI tool for conversion",
			Long:    ``,
			Version: versionStr,
		},
	}
	c.Command.SetVersionTemplate("circe version {{.Version}}\n")
	cobra.OnInitialize(c.initConfig)
	return c
}

func (c *RootCommand) initConfig() {

}
