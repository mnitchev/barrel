package cmd

import (
	"fmt"
	"os"

	"github.com/mnitchev/barrel/runner"
	"github.com/spf13/cobra"
)

var barrelCmd = &cobra.Command{
	Use:   "roll",
	Short: "starts a the specified command in a container",
	Args:  cobra.MinimumNArgs(1),
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
	exitCode, err := runner.Run(container)
	if err != nil {
		fmt.Printf("Running command failed. Error: %s\n", err)
	}

	os.Exit(exitCode)
}
