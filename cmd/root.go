package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "barrel",
	Short: "Barrel is a desperate attempt to think of an original name for a container-rutnime-thingy",
}

func init() {
	initRollCommand()
	initPinCommand()
	initLimitMemoryCommand()
	rootCmd.AddCommand(barrelCmd)
	rootCmd.AddCommand(pinCmd)
	rootCmd.AddCommand(limitMemoryCmd)
}

func Execute() {
	if err := barrelCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
