package networks

// These structs are copied by hand from localcloud-daemon/internal/shared/ipc/ipc.go.
// Keep them in sync: the JSON tags are the wire format.

type REQUEST_TYPE string

const (
	GET    REQUEST_TYPE = "GET"
	POST   REQUEST_TYPE = "POST"
	PUT    REQUEST_TYPE = "PUT"
	DELETE REQUEST_TYPE = "DELETE"
)

type IPCRequest[T any] struct {
	RequestID string       `json:"id"`
	Resource  string       `json:"resource"`
	Type      REQUEST_TYPE `json:"type"`
	Payload   T            `json:"payload"`
}

type IPCResponse[T any] struct {
	RequestID string `json:"requestID"`
	Resource  string `json:"resource"`
	Error     string `json:"error,omitempty"`
	Response  T      `json:"response,omitempty"`
}

// ClientHello is the payload the daemon sends when a client connects
// (the anonymous struct in the daemon's SocketServer.newConnection).
type ClientHello struct {
	ID string
}
