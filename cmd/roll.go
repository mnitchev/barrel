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
	rootfs, err := cmd.Flags().GetString("rootfs")
	if err != nil {
		fmt.Errorf("Failed to get rootfs path from flags: %s", err)
		panic(err)
	}
	container := runner.Container{
		Command:    args[0],
		Args:       args[1:],
		Stdin:      os.Stdin,
		Stdout:     os.Stdout,
		Stderr:     os.Stdout,
		RootfsPath: rootfs,
	}
	exitCode, err := runner.Run(container)
	if err != nil {
		fmt.Printf("Running command failed. Error: %s\n", err)
	}

	os.Exit(exitCode)
}
