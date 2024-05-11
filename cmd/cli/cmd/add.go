package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

func (c *Cli) AddCmd() *cobra.Command {
	// addCmd represents the add command
	var addCmd = &cobra.Command{
		Use:   "add [data]",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// fmt.Println(args)
			data := args[0]
			err := c.bc.AddBlock(data)

			if err != nil {
				fmt.Fprintf(os.Stderr, "%s", err)
				runtime.Goexit()
			}

			fmt.Printf("Added block: %x\n", c.bc.LastHash)
		},
	}

	return addCmd
}
