# lesiwlabs edge compute

This is a Go WASM project that manages redirects for @lesiw’s domains. It uses
the [`workers`][workers] library to glue together Go’s `net/http` handlers and
CloudFlare’s worker API.

By building this project in Go, it trades some slight overhead for portability
between CloudFlare and other edge compute services that are capable of running
WASM apps. New glue code libraries may need to be written on the occasion of
moving from one service to another, but the core logic will remain the same.

This project was generated from the [worker-tinygo template][template].

## Requirements

- [wrangler](https://developers.cloudflare.com/workers/wrangler/)
- tinygo 0.29.0 or later

## Development

### Commands

```
make dev     # run dev server
make build   # build Go Wasm binary
make deploy  # deploy worker
```

### Testing dev server

```
$ curl http://localhost:8787/hello
Hello!
```

```
$ curl -X POST -d "test message" http://localhost:8787/echo
test message
```

[workers]: https://github.com/syumai/workers
[template]: https://github.com/syumai/workers/tree/main/_templates/cloudflare/worker-tinygo
