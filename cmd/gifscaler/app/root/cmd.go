package root

import (
	"os"

	"github.com/spf13/cobra"
)

type cliopts struct {
}

type rootCommand struct {
	options cliopts
}

func (c *rootCommand) run(cmd *cobra.Command, args []string) {
	os.Exit(run(c.options, args))
}

func CreateCmd() *cobra.Command {
	c := &rootCommand{}

	cmd := &cobra.Command{
		Use:   "root",
		Short: "Resize gifs with scale2x",
		Long:  `Resize gifs with scale2x`,
		Run:   c.run,
	}

	return cmd
}
