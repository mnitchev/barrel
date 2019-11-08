package cmd

import (
	"os"

	"github.com/mnitchev/barrel/cgroups"
	"github.com/spf13/cobra"
)

var pinCmd = &cobra.Command{
	Use:   "pin-cpu",
	Short: "pins the specified cgroup cpuset to a cpu. use in conjunction with the --cgroup flag of the roll command",
	Run:   pinCommand,
}

func initPinCommand() {
	pinCmd.Flags().StringP("cgroup", "c", "", "name of the cgroup")
	pinCmd.MarkFlagRequired("cgroup")
	pinCmd.Flags().StringP("cpu-indexes", "i", "", "list of cpu indexes. Example: 1-4,6,7-9")
	pinCmd.MarkFlagRequired("cpu-indexes")
}

func pinCommand(cmd *cobra.Command, args []string) {
	cgroupName := getStringFlag(cmd, "cgroup")
	cpuIndexes := getStringFlag(cmd, "cpu-indexes")
	if err := cgroups.PinCPU(cgroupName, cpuIndexes); err != nil {
		os.Exit(1)
	}
}
