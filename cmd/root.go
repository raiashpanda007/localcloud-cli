package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"raiashpanda007/localcloud-cli/internals/utils"
)

// sysError is the cli's central error bus (same pattern as the daemon).
// Set up in Execute before any command runs.
var sysError utils.SystemError

var rootCmd = &cobra.Command{
	Use:   "localcloud",
	Short: "Pool computing resources across your machines",
	Long: `Local Cloud connects computers you control and combines their computing
resources into a single pool.

A master coordinates connected machines, manages workers, and distributes
workloads. Workers run on participating machines and execute assigned tasks.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sysError = utils.Init()
	go sysError.Handler()

	if err := rootCmd.ExecuteContext(rootCtx); err != nil {
		stop()
		os.Exit(1)
	}
}
