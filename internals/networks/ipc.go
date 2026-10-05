package networks

import (
	"log/slog"
	"net"
)

const (
	WorkerSocketPath = "/tmp/worker.sock"
	MasterSocketPath = "/tmp/master.sock"
)

func ConnectUnixSocket(socketpath string) error {

	conn, err := net.Dial("unix", socketpath)

	if err != nil {
		return err
	}

	slog.Info("Connect to socket path :: ", socketpath)
	defer conn.Close()
	return nil
}
