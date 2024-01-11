package main

import (
	"fmt"
	"os"

	"github.com/eduardooliveira/stLib/cli/commands"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mmp",
	Short: "Interact with MMP through commands!",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	importCommand := commands.InitImport()

	rootCmd.AddCommand(importCommand)
}

func main() {
	Execute()
}
