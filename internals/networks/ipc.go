package networks

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"syscall"
	"time"
)

const (
	// Must match SOCKET_SERVER_PATH in localcloud-daemon/main.go.
	SOCKET_SERVER_PATH string = "/tmp/daemon.sock"

	REPLY_TIMEOUT = 2 * time.Second
)

// ErrNoDaemon means nothing is listening on the socket: either the socket
// file doesn't exist, or it is stale (left behind by a crashed daemon).
var ErrNoDaemon = errors.New("no daemon is running on this machine")

// Status asks the running daemon which role it runs as.
func Status(ctx context.Context) (Response, error) {
	return sendRequest(ctx, SOCKET_SERVER_PATH, Request{Type: REQUEST_STATUS})
}

func sendRequest(ctx context.Context, socketpath string, req Request) (Response, error) {
	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "unix", socketpath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, syscall.ECONNREFUSED) {
			return Response{}, ErrNoDaemon
		}
		return Response{}, fmt.Errorf("unable to connect to daemon socket %s: %w", socketpath, err)
	}
	defer conn.Close()

	payload, err := json.Marshal(req)
	if err != nil {
		return Response{}, fmt.Errorf("unable to encode %q request: %w", req.Type, err)
	}

	// The deadline covers both the write and the reply, so a daemon that never
	// answers (e.g. an old build without the protocol) can't hang the cli.
	if err := conn.SetDeadline(time.Now().Add(REPLY_TIMEOUT)); err != nil {
		return Response{}, fmt.Errorf("unable to set socket deadline: %w", err)
	}

	if _, err := conn.Write(append(payload, '\n')); err != nil {
		return Response{}, fmt.Errorf("unable to send %q request to daemon: %w", req.Type, err)
	}

	line, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return Response{}, fmt.Errorf("daemon did not reply to %q within %s (is it an old build without the IPC protocol?)", req.Type, REPLY_TIMEOUT)
		}
		return Response{}, fmt.Errorf("unable to read %q reply from daemon: %w", req.Type, err)
	}

	var resp Response
	if err := json.Unmarshal(line, &resp); err != nil {
		return Response{}, fmt.Errorf("invalid reply from daemon %q: %w", string(line), err)
	}

	if !resp.OK {
		return resp, fmt.Errorf("daemon rejected %q request: %s", req.Type, resp.Error)
	}

	return resp, nil
}
