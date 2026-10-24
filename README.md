# Yet Another Fibonacci Calculator

A small Go HTTP API that returns the first N Fibonacci numbers, starting at F(0) = 0.

## Run

Requires Go 1.26.4 or later. The project uses only the standard library.

```sh
go run ./cmd/yafibserver
```

The server listens on port 8080.

```sh
curl "http://localhost:8080/api/yafib?count=5"
```

```json
{"numbers":[0,1,1,2,3]}
```

To build an executable:

```sh
go build -o yafibserver ./cmd/yafibserver
```

On Windows, use an output name of yafibserver.exe.

## API

GET /api/yafib accepts a required integer query parameter, count.

- Counts 0 through 94 return HTTP 200 with a numbers array.
- A count of zero returns {"numbers":[]}.
- Missing, empty, malformed, negative, or oversized counts return HTTP 400.
- Unsupported methods return HTTP 405 with Allow: GET.
- Unexpected calculation failures return HTTP 500 without exposing internal error details.

Errors use the JSON shape {"error":"message"}. Responses have Content-Type: application/json.

## Design

The engine builds the complete sequence iteratively in O(N) time and O(N) space
for the returned slice. This matches an API that returns every term: computing
one term faster does not remove the work needed to produce the full sequence.
Caching adds unnecessary state for this small, bounded calculation.

Values use uint64. The largest representable Fibonacci number is
F(93) = 12200160415121876738, so a sequence starting at F(0) can contain at most
94 terms. The engine rejects larger counts before allocating the result.
This bound also makes an overflow check inside the loop unnecessary.

The HTTP package consumes a small SequenceGenerator interface, allowing request
handling to be tested independently of calculation. Engine validation errors are
identified using errors.Is and mapped to HTTP 400; unexpected failures are logged
and mapped to HTTP 500.

The executable wires the engine and handler together with a dedicated HTTP mux.
A five-second header timeout bounds how long clients can take to send headers.
The listen address is fixed at :8080 for this exercise.

## Logging

The executable creates a standard-library log/slog logger with JSON output to
standard error and passes it into the HTTP server. Each record includes a time,
level, message, and structured fields.

- INFO: server startup, including the listen address.
- ERROR: unexpected calculation failures, with method, path, count, and error.
- ERROR: response-writing failures, with the intended HTTP status and error.
- ERROR: server failures, with the listen address and error. Fatal server
  failures exit with status 1.

The HTTP server's internal error logger uses the same structured handler.
Expected input validation failures are returned to the client without error
logs. Successful requests are not logged individually; access logging is
omitted to keep the example focused. The calculation engine does not log.

For example, a startup record has this shape:

```json
{"time":"2026-10-24T12:00:00Z","level":"INFO","msg":"starting server","addr":":8080"}
```

## Layout

```text
cmd/yafibserver/   Executable entry point
internal/engine/  Sequence calculation and boundary tests
internal/server/  HTTP handler, contract tests, and engine integration tests
```

## Verify

```sh
go test ./... -cover
go vet ./...
gofmt -l cmd internal
```

Tests cover sequence values, both count boundaries, HTTP validation, supported
methods, response headers, error mapping, and real engine responses through the
handler. The formatting command prints nothing when all files are formatted.

## Numeric clients

JSON contains integer literals. Clients that decode all numbers as IEEE 754
double-precision values, including JavaScript's default JSON parser, can lose
precision for the larger terms. Exact consumers should use a JSON decoder that
preserves integer precision.
