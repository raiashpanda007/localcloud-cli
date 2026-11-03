package networks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"syscall"
	"time"
)

const (
	// Must match the daemon's socket paths (master.go and worker.go).
	MASTER_SOCKET_PATH = "/tmp/daemon-master.sock"
	WORKER_SOCKET_PATH = "/tmp/daemon-worker.sock"

	REPLY_TIMEOUT = 2 * time.Second
)

// ErrNoDaemon means nothing is listening on the role's socket: the socket file
// is missing, or it is stale and refuses connections.
var ErrNoDaemon = errors.New("no daemon is listening")

// socketPath returns the socket for role, or an error for an unknown role.
func socketPath(role string) (string, error) {
	switch role {
	case "master":
		return MASTER_SOCKET_PATH, nil
	case "worker":
		return WORKER_SOCKET_PATH, nil
	default:
		return "", fmt.Errorf("unknown role %q", role)
	}
}

// Connect dials the role's socket and reads the greeting the daemon sends on connect.
// It returns ErrNoDaemon when no daemon is listening.
func Connect(ctx context.Context, role string) (ClientHello, error) {
	path, err := socketPath(role)
	if err != nil {
		return ClientHello{}, err
	}
	return connectPath(ctx, path)
}

func connectPath(ctx context.Context, path string) (ClientHello, error) {
	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "unix", path)
	if err != nil {
		if errors.Is(err, syscall.ENOENT) || errors.Is(err, syscall.ECONNREFUSED) {
			return ClientHello{}, ErrNoDaemon
		}
		return ClientHello{}, err
	}
	defer conn.Close()

	if err := conn.SetReadDeadline(time.Now().Add(REPLY_TIMEOUT)); err != nil {
		return ClientHello{}, err
	}

	// The greeting has no trailing newline, so decode one JSON value instead of reading a line.
	var greeting IPCResponse[ClientHello]
	if err := json.NewDecoder(conn).Decode(&greeting); err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return ClientHello{}, fmt.Errorf("daemon did not greet within %s (is it an old build?)", REPLY_TIMEOUT)
		}
		return ClientHello{}, fmt.Errorf("unable to read daemon greeting: %w", err)
	}

	if greeting.Error != "" {
		return ClientHello{}, fmt.Errorf("daemon refused the connection: %s", greeting.Error)
	}
	if greeting.Resource != "CLIENT" {
		return ClientHello{}, fmt.Errorf("unexpected greeting resource %q", greeting.Resource)
	}

	return greeting.Response, nil
}
