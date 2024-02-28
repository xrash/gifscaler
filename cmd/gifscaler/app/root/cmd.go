package root

import (
	"os"

	"github.com/spf13/cobra"
)

type cliopts struct {
	keepWorkdir bool
	usePalette  bool
	scale       uint
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

	cmd.Flags().BoolVarP(
		&c.options.keepWorkdir,
		"keep-workdir",
		"",
		false,
		"Keep the temporary directory in which all the work is done.",
	)

	cmd.Flags().BoolVarP(
		&c.options.keepWorkdir,
		"use-palette",
		"",
		true,
		"Generate a palette from all frames and use it when producing the output. This will maintain the transparency of the frames.",
	)

	cmd.Flags().UintVarP(
		&c.options.scale,
		"scale",
		"",
		2,
		"What scale to use when rescaling the frames with scale2x. This will be passed as the \"-k\" parameter to scale2x.",
	)

	return cmd
}
