package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"raiashpanda007/localcloud-cli/internals/daemon"
	"raiashpanda007/localcloud-cli/internals/networks"
)

// startNode makes this machine run the daemon as role. It reuses a daemon
// that is already running as role, refuses if it is running as the other
// role, and spawns a new daemon if none is running.
func startNode(ctx context.Context, out io.Writer, role string) error {
	other := "worker"
	if role == "worker" {
		other = "master"
	}

	_, err := networks.Connect(ctx, role)
	switch {
	case err == nil:
		fmt.Fprintf(out, "This node is already running as %s.\n", role)
		return nil
	case !errors.Is(err, networks.ErrNoDaemon):
		return err
	}

	// Dev only: LOCALCLOUD_DEV=1 lets one machine run both roles.
	if os.Getenv("LOCALCLOUD_DEV") != "1" {
		_, err = networks.Connect(ctx, other)
		switch {
		case err == nil:
			return fmt.Errorf("this node is running as %s; stop it first, then run `localcloud %s start` again", other, role)
		case !errors.Is(err, networks.ErrNoDaemon):
			return err
		}
	}

	process, err := daemon.Spawn(role)
	if err != nil {
		return err
	}

	if err := daemon.WaitUntilReady(ctx, process, role); err != nil {
		return err
	}

	fmt.Fprintf(out, "Started %s daemon (pid %d). Logs: %s\n", role, process.Process.Pid, daemon.DAEMON_LOG_PATH)
	return nil
}
