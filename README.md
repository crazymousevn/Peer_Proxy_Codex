# Peer_Proxy_Codex

Peer_Proxy_Codex contains the initial implementation for the PeerProxy project.

## Signaling Server

A basic signaling server is available under `cmd/signaling`. It exposes an HTTP
endpoint at `/signal` on port `8080`.

Run the server:

```bash
go run cmd/signaling/main.go
```

Send a POST request with body `hello` to receive `hello world` in response.
Requests must include an `Authorization` header with the value `Bearer secret`.

The endpoint also accepts a WebRTC SDP offer (JSON encoded) and responds with
an SDP answer to establish a peer connection.

Once the peer connection is active the server listens for WebRTC data channels
labelled `socks5` or `http`.

- A data channel labelled `socks5` expects a single SOCKS5 session. The server
  supports the `CONNECT` command for IPv4 addresses and will proxy TCP traffic
  to the requested destination.
- A data channel labelled `http` accepts HTTP proxy requests. Standard HTTP
  requests with absolute URLs are forwarded using the `net/http` client, and
  `CONNECT` requests establish a tunnel for HTTPS traffic.

## Development

Common tasks are available via the Makefile:

```bash
make build
make test
```

