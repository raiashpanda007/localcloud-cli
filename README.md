# localcloud-cli

`localcloud` is the command-line client for Local Cloud. It doesn't run the
master or worker itself; it controls the **daemon** (`localcloud-daemon`) over
one unix socket per role.

## Development setup

The CLI spawns the daemon from a **hardcoded path** (development only):

```
/Users/ashwinrai/Projects/local-cloud/localcloud-daemon/bin/local-cloud-daemon
```

> TODO: this path is hardcoded in `internals/daemon/spawn.go` (`DAEMON_BIN_PATH`).
> It will be replaced when the CLI is built and distributed. Later, the daemon
> will run as a systemd service instead of being spawned by the CLI.

Build the daemon to that path, then run the CLI:

```sh
cd localcloud-daemon && go build -o bin/local-cloud-daemon .
cd localcloud-cli && go run . master start
```

## `localcloud master start` / `localcloud worker start`

`start` first connects to the role's socket and reads the greeting the daemon
sends on connect:

| Daemon state | Result |
| --- | --- |
| Running as the requested role (its socket greets) | Prints `This node is already running as <role>.`, exit 0 |
| Running as the other role | Refuses and tells you to stop it first, exit 1 |
| Not running (no socket, or a stale one that refuses) | Spawns the daemon detached, waits up to 5 s for the greeting, prints `Started <role> daemon (pid N)` |

Sockets (must match the daemon):

| Role | Path |
| --- | --- |
| master | `/tmp/daemon-master.sock` |
| worker | `/tmp/daemon-worker.sock` (the daemon doesn't open it yet, so worker checks only report "not running" until it does) |

The greeting is `{"requestID":"","resource":"CLIENT","response":{"ID":"<uuid>"}}`
with no trailing newline. `internals/networks/protocol.go` holds copies of the
daemon's `IPCRequest`/`IPCResponse` structs, which must be kept in sync by hand.

- The spawned daemon's logs are appended to `/tmp/localcloud-daemon.log`
  (`tail -f` it). Environment variables are passed through, e.g.
  `LOG_LEVEL=debug LOG_FORMAT=text localcloud master start`.
- The daemon runs in its own session, so it keeps running after the CLI exits
  and Ctrl+C in the terminal doesn't stop it. Stop it with `kill <pid>`.
- If a daemon crashed and left a stale socket file, the connect is refused and
  `start` spawns a new daemon, which then fails with "address already in use"
  until the stale file is removed (`rm -f <socket path>`).
