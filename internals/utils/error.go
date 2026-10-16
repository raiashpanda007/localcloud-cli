package utils

import (
	"log/slog"
	"os"
	"runtime/debug"
)

/*
Mirrors the daemon's error system (localcloud-daemon/internal/utils/error.go).
1. Invalid data sort off.
2. Service not responded (Not a big problem to kill the cli)
3. Not solvable have to kill the cli
*/
type CliError struct {
	Error   error
	Trace   string
	Message string
	Code    int
}

type SystemError struct {
	rootErrorChan chan CliError
}

func Init() SystemError {
	errChan := make(chan CliError)

	return SystemError{
		rootErrorChan: errChan,
	}
}

func (r *SystemError) NewCliError(err error, message string, code int) {
	var errRecieved = CliError{
		Error:   err,
		Message: message,
		Trace:   string(debug.Stack()),
		Code:    code,
	}

	slog.Error("ERROR RECIEVED ::", "message", message, "error", err, "trace", errRecieved.Trace)

	r.rootErrorChan <- errRecieved
}

func (r *SystemError) Handler() {
	for {
		cliError := <-r.rootErrorChan

		if cliError.Code == 1 {
			slog.Warn("Invalid data received; continuing cli", "message", cliError.Message, "error", cliError.Error)
			continue
		}

		if cliError.Code == 2 {
			slog.Warn("Service did not respond; continuing cli", "message", cliError.Message, "error", cliError.Error)
			continue
		}

		if cliError.Code == 3 {
			slog.Error("Can't resolve it killing the cli", "message", cliError.Message, "error", cliError.Error, "trace", cliError.Trace)
			os.Exit(1)
		}
	}
}
