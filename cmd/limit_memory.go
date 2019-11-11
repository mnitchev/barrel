package cmd

import (
	"os"

	"github.com/mnitchev/barrel/cgroups"
	"github.com/spf13/cobra"
)

var limitMemoryCmd = &cobra.Command{
	Use:   "limit-memory",
	Short: "pins the specified cgroup cpuset to a cpu. use in conjunction with the --cgroup flag of the roll command",
	Run:   limitMemoryCommand,
}

func initLimitMemoryCommand() {
	limitMemoryCmd.Flags().StringP("cgroup", "c", "", "name of the cgroup")
	limitMemoryCmd.MarkFlagRequired("cgroup")
	limitMemoryCmd.Flags().StringP("max", "m", "", "maximum amount of memory. Example: 1024, 50M, 30K")
	limitMemoryCmd.MarkFlagRequired("nax")
}

func limitMemoryCommand(cmd *cobra.Command, args []string) {
	cgroupName := getStringFlag(cmd, "cgroup")
	max := getStringFlag(cmd, "max")
	if err := cgroups.LimitMemory(cgroupName, max); err != nil {
		os.Exit(1)
	}
}
