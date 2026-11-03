package daemon

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"raiashpanda007/localcloud-cli/internals/networks"
)

const (
	// TODO: hardcoded for development only. Replace when the cli is built and
	// distributed (the daemon will be managed by systemd instead of spawned here).
	// Build it with: cd localcloud-daemon && go build -o bin/local-cloud-daemon .
	DAEMON_BIN_PATH string = "/Users/ashwinrai/Projects/local-cloud/localcloud-daemon/bin/local-cloud-daemon"

	// The spawned daemon's stdout/stderr are appended here.
	DAEMON_LOG_PATH string = "/tmp/localcloud-daemon.log"

	READY_TIMEOUT       = 5 * time.Second
	READY_POLL_INTERVAL = 100 * time.Millisecond
)

// Spawn starts the daemon as a detached process running the given role.
func Spawn(role string) (*exec.Cmd, error) {
	if _, err := os.Stat(DAEMON_BIN_PATH); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("daemon binary not found at %s; build it with `cd localcloud-daemon && go build -o bin/local-cloud-daemon .`", DAEMON_BIN_PATH)
		}
		return nil, fmt.Errorf("unable to check daemon binary %s: %w", DAEMON_BIN_PATH, err)
	}

	logFile, err := os.OpenFile(DAEMON_LOG_PATH, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("unable to open daemon log file %s: %w", DAEMON_LOG_PATH, err)
	}
	// The child gets its own copy of the descriptor, so the cli can close its copy.
	defer logFile.Close()

	cmd := exec.Command(DAEMON_BIN_PATH, role)
	cmd.Stdin = nil
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	// New session: Ctrl+C in the cli's terminal doesn't reach the daemon,
	// and the daemon keeps running after the cli exits.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("unable to start daemon %s: %w", DAEMON_BIN_PATH, err)
	}

	return cmd, nil
}

// WaitUntilReady polls the role's socket until the daemon greets, the spawned
// process exits, the timeout passes, or ctx is cancelled.
func WaitUntilReady(ctx context.Context, cmd *exec.Cmd, role string) error {
	exited := make(chan error, 1)
	go func() {
		exited <- cmd.Wait()
	}()

	timeout := time.NewTimer(READY_TIMEOUT)
	defer timeout.Stop()
	ticker := time.NewTicker(READY_POLL_INTERVAL)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-exited:
			return fmt.Errorf("daemon exited before it was ready (%v); see %s", err, DAEMON_LOG_PATH)
		case <-timeout.C:
			return fmt.Errorf("daemon was not ready within %s; see %s", READY_TIMEOUT, DAEMON_LOG_PATH)
		case <-ticker.C:
			_, err := networks.Connect(ctx, role)
			if errors.Is(err, networks.ErrNoDaemon) {
				continue // socket not up yet
			}
			if err != nil {
				return err
			}
			return nil
		}
	}
}
