package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"

	"raiashpanda007/localcloud-cli/internals/daemon"
	"raiashpanda007/localcloud-cli/internals/networks"
)

// startNode makes this machine run the daemon as role. It reuses a daemon
// that is already running as role, refuses if it is running as the other
// role, and spawns a new daemon if none is running.
func startNode(ctx context.Context, out io.Writer, role string) error {
	resp, err := networks.Status(ctx)

	switch {
	case err == nil && resp.Role == role:
		fmt.Fprintf(out, "This node is already running as %s (pid %d).\n", role, resp.PID)
		return nil

	case err == nil:
		return fmt.Errorf("this node is running as %s (pid %d); stop it first (kill %d), then run `localcloud %s start` again", resp.Role, resp.PID, resp.PID, role)

	case errors.Is(err, networks.ErrNoDaemon):
		process, err := daemon.Spawn(role)
		if err != nil {
			return err
		}

		resp, err := daemon.WaitUntilReady(ctx, process, role)
		if err != nil {
			return err
		}

		fmt.Fprintf(out, "Started %s daemon (pid %d). Logs: %s\n", role, resp.PID, daemon.DAEMON_LOG_PATH)
		return nil

	default:
		return err
	}
}
