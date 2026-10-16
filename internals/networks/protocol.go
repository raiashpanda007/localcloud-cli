package networks

// Wire protocol between the CLI and the daemon over the unix socket.
// Every message is a single-line JSON object (json.Marshal) followed by "\n".
// One request gets exactly one response. Keep these structs in sync with the daemon.

const (
	REQUEST_STATUS = "status"
)

type Request struct {
	Type string `json:"type"`
}

type Response struct {
	OK    bool   `json:"ok"`
	Role  string `json:"role,omitempty"` // "master" | "worker"
	PID   int    `json:"pid,omitempty"`
	Error string `json:"error,omitempty"`
}
