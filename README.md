# localcloud-cli

`localcloud` is the command-line client for Local Cloud. It doesn't run the
master or worker itself; it controls the **daemon** (`localcloud-daemon`)
over a unix socket at `/tmp/daemon.sock`.

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

Only one daemon runs per machine. `start` asks it for its status over the
socket and then:

| Daemon state | Result |
| --- | --- |
| Running as the requested role | Prints `This node is already running as <role> (pid N).`, exit 0 |
| Running as the other role | Refuses and tells you to stop it first (`kill N`), exit 1 |
| Not running (no socket, or a stale one) | Spawns the daemon detached, waits up to 5 s for it to answer, prints `Started <role> daemon (pid N)` |

- The spawned daemon's logs are appended to `/tmp/localcloud-daemon.log`
  (`tail -f` it). Environment variables are passed through, e.g.
  `LOG_LEVEL=debug LOG_FORMAT=text localcloud master start`.
- The daemon runs in its own session, so it keeps running after the CLI exits
  and Ctrl+C in the terminal doesn't stop it. Stop it with `kill <pid>`.
- If a crashed daemon left a stale `/tmp/daemon.sock`, the daemon removes it
  on startup. Until it does, run `rm -f /tmp/daemon.sock`.
- If the daemon doesn't reply within 2 s, `start` fails instead of hanging.

## IPC protocol

Every message is one JSON object (`json.Marshal`) followed by `\n`, and each
request gets exactly one response. The structs live in
`internals/networks/protocol.go` and must match the daemon's.

```
→ {"type":"status"}
← {"ok":true,"role":"master","pid":1234}
← {"ok":false,"error":"unknown request type"}
```

To fake a daemon by hand:

```sh
nc -lU /tmp/daemon.sock      # then type a reply line after the request arrives
```
