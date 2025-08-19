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
			Short: "generates network policies based on inputs from csv",
		},
	}
	c.command.Flags().StringVarP(&c.input, "input", "i", "", "csv file")
	c.command.Flags().StringVarP(&c.output, "output", "o", ".", "output directory to save egress policies, default is current directory")
	c.command.Flags().IntVarP(&c.headerStart, "header", "", 0, "header starting index in csv, which will indicate which row to consider as header, default is 0")
	c.command.Run = c.Run
	return c
}

func (c *IngressGenerateCommand) Run(command *cobra.Command, args []string) {
	var unmarshalled []unmarshalcsv.UnmarshalledData
	unmarshal, err := unmarshalcsv.NewUnmarshalCsv(c.input, c.headerStart)
	if err != nil {
		panic(err)
	}
	err = unmarshal.UnmarshalCsv(&unmarshalled)
	if err != nil {
		panic(err)
	}
	// Use generic policies filtered to Ingress only
	n := netpol.NewGenericPoliciesForDirection(unmarshalled, c.output, "Ingress")
	err = n.RenderGeneric()
	if err != nil {
		panic(err)
	}
}
