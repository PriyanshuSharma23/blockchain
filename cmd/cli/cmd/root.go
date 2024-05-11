package cmd

import (
	"runtime"

	"github.com/PriyanshuSharma23/custom_blockchain/internals/blockchain"
	"github.com/spf13/cobra"
)

type Cli struct {
	bc      *blockchain.Blockchain
	rootCmd *cobra.Command
}

func NewCli(bc *blockchain.Blockchain) *Cli {
	var rootCmd = &cobra.Command{
		Use:   "custom_blockchain",
		Short: "CLI to interact with blockchain",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	}

	cli := &Cli{
		bc:      bc,
		rootCmd: rootCmd,
	}

	cli.Init()

	return cli
}

func (c *Cli) Init() {
	c.rootCmd.AddCommand(c.PrintCmd())
	c.rootCmd.AddCommand(c.AddCmd())
}

func (c *Cli) Execute() {
	err := c.rootCmd.Execute()
	if err != nil {
		runtime.Goexit()
	}
}
