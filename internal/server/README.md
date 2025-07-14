# Zerodupe Server

Backed server for deduplication file storage system that splits files into chunks and only stores unique chunks, saving storage space.

## Requirements

- Go Version (Go 1.20+)
- Docker

## Server Components

The server components follow a clean architecture with clear separation of concerns:

## Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│             │     │             │     │             │
│     API     │────▶│   Handler   │────▶│   Storage   │
│             │     │             │     │             │
└─────────────┘     └─────────────┘     └─────────────┘
```

## Running the Server

From the project root:

```sh
go run ./cmd/zerodupe-server/main.go
```

Or build and run:

```sh
go build -o zerodupe-server ./cmd/zerodupe-server
./zerodupe-server
```

## **Running Tests**

```sh
go test ./...
```