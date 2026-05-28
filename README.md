# Auth Service

gRPC-based authentication service with PostgreSQL.

## Prerequisites

- Go 1.26+
- PostgreSQL 18+
- protobuf 34+

## Quick Start

```sh
# Install dependencies
go mod download
```

## Testing

Integration tests use [testcontainers-go](https://github.com/testcontainers/testcontainers-go) with PostgreSQL.

```sh
cd testing
go test ./...
```

## Project Structure

```
.
├── cmd/          # CLI entry points
├── gen/          # Generated protobuf code (ignored)
├── migrations/   # SQL migrations
├── proto/        # Protobuf service definitions
└── testing/      # Integration tests
```
