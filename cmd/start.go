package cmd

import (
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the master on this machine",
	Long: `Start runs the master on this machine.

The master coordinates connected machines, discovers workers, manages
connections, and distributes workloads.

If a master daemon is already running, start reports it. If a worker daemon
is running, start refuses; stop the worker first. Otherwise it launches the
daemon as a master.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return startNode(cmd.Context(), cmd.OutOrStdout(), "master")
	},
}

func init() {
	masterCmd.AddCommand(startCmd)
}
