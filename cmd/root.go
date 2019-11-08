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
	barrelCmd.Flags().StringP("rootfs", "r", "", "path to rootfs for the contained process")
	barrelCmd.MarkFlagRequired("rootfs")
	barrelCmd.Flags().StringP("cgroup-name", "c", "", "name of the cgroup to add the taskt to. will create a new one if missing")
	barrelCmd.MarkFlagRequired("cgroup-name")
	rootCmd.AddCommand(barrelCmd)
}

func Execute() {
	if err := barrelCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
