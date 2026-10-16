package cmd

import (
	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Execute assigned tasks on this machine",
	Long: `The worker runs on a participating machine. It executes tasks assigned
by the master.

Remote connectivity is planned for a later release.`,
}

var workerStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the worker on this machine",
	Long: `Start runs the worker on this machine.

The worker executes tasks assigned by the master.

If a worker daemon is already running, start reports it. If a master daemon
is running, start refuses; stop the master first. Otherwise it launches the
daemon as a worker.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return startNode(cmd.Context(), cmd.OutOrStdout(), "worker")
	},
}

func init() {
	workerCmd.AddCommand(workerStartCmd)
	rootCmd.AddCommand(workerCmd)
}
