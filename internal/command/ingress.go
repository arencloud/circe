package command

import (
	"circe/pkg/netpol"
	"circe/pkg/unmarshalcsv"
	"github.com/spf13/cobra"
)

type IngressGenerateCommand struct {
	command     *cobra.Command
	input       string
	output      string
	headerStart int
}

func NewIngressCommand() *IngressGenerateCommand {
	c := &IngressGenerateCommand{
		command: &cobra.Command{
			Use:   "ingress",
			Short: "generates network policies based on inputs from CSV or XLSX",
		},
	}
	c.command.Flags().StringVarP(&c.input, "input", "i", "", "input file (CSV or XLSX)")
	c.command.Flags().StringVarP(&c.output, "output", "o", ".", "output directory to save egress policies, default is current directory")
	c.command.Flags().IntVarP(&c.headerStart, "header", "", 0, "header starting index in the input (CSV/XLSX), indicating which row to treat as header; default is 0")
	c.command.Run = c.Run
	return c
}

func (c *IngressGenerateCommand) Run(command *cobra.Command, args []string) {
	var unmarshalled []unmarshalcsv.UnmarshalledData
	if err := unmarshalcsv.Unmarshal(&unmarshalled, c.input, c.headerStart); err != nil {
		panic(err)
	}
	// Use generic policies filtered to Ingress only
	n := netpol.NewGenericPoliciesForDirection(unmarshalled, c.output, "Ingress")
	if err := n.RenderGeneric(); err != nil {
		panic(err)
	}
}
