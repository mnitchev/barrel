package cmd

import (
	"os"

	"github.com/mnitchev/barrel/runner"
	"github.com/spf13/cobra"
)

var barrelCmd = &cobra.Command{
	Use:   "roll",
	Short: "starts a the specified command in a container",
	Run:   rollCommand,
}

func rollCommand(cmd *cobra.Command, args []string) {
	container := runner.Container{
		Command: args[0],
		Args:    args[1:],
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stdout,
	}
	if err := runner.Run(container); err != nil {
		panic(err)
	}
}
