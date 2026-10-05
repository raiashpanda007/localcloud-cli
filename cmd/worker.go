package cmd

import (
	"errors"

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

The worker executes tasks assigned by the master.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("worker does not run inside the cli")
	},
}

func init() {
	workerCmd.AddCommand(workerStartCmd)
	rootCmd.AddCommand(workerCmd)
}

func StartWorker() error {
	// TODO: Connect to worker IPC and start the worker process
}
