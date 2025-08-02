# ZeroDupe

zerodupe is a content-addressable storage server designed for space efficiency through block-level deduplication. Clients interact with the server via a REST API to upload and download data. Files are split into fixed-size blocks, each identified by the hash of its content. The server stores each unique block only once. Metadata linking files to their constituent blocks is maintained separately.

## Architecture

The system is divided into three main layers:

1. **CMD**: The cmd layer is the entry point for building server and client binaries.
2. **Client**: The client layer is responsible for interacting with the user and sending requests to the server. It is implemented in the `pkg/client` directory.
3. **Server**: The server layer is responsible for handling requests from the client and storing/retrieving data. It is implemented in the `internal/server` directory.
4. **Hasher**: The hasher layer is responsible for calculating hashes for file and chunk data. It is implemented in the `pkg/hasher` directory and is used by both the client and server.

## 🐳 Running ZeroDupe with Docker Compose

### Prerequisites

- Go
- Docker

1. Build and start the server:

The server is a long-running service. Start it with

```bash
docker-compose up d zerodupe-server
```

- This will build the server image (if needed) and start it in the background.
- The server will listen on port 8080 by default.

2. Run Client Commands as Needed
   The client is a command-line tool for one-off actions (signup, upload, download, etc.).
   You do not keep the client running; instead, you run it only when you need it.

### Example: Sign up a new user

```bash
docker-compose run --rm zerodupe-client signup --server http://zerodupe-server:8080 --username <username> --password <password> --confirm-password <password>
```

Replace <username> and <password> with your desired credentials.

### Example: Upload a file

```bash
docker-compose run --rm \
  -v $(pwd)/path/to/file.txt:/app/file.txt \
  zerodupe-client upload --server http://zerodupe-server:8080 /app/file.txt
```

### Example: Download a file

```bash
docker-compose run --rm \
  -v $(pwd)/downloads:/app/downloads \
  zerodupe-client download --server http://zerodupe-server:8080 -o /app/downloads -n output.txt <FILE_HASH>
```

Replace <FILE_HASH> with the hash of the file you want to download.

### Stopping the Server

When you’re done, stop the server and clean up resources with:

```bash
docker-compose down
```

---

## Server Configuration

You can configure the server using environment variables or command-line flags:

| Flag / Env Var                                             | Description                   | Default      |
| ---------------------------------------------------------- | ----------------------------- | ------------ |
| `--port`, `PORT`                                           | Server port                   | 8080         |
| `--storage`, `STORAGE_DIR`                                 | Storage directory             | data/storage |
| `--secret`, `JWT_SECRET`                                   | JWT Secret (required)         |              |
| `--access-token-expiry-min`, `ACCESS_TOKEN_EXPIRY_MIN`     | Access token expiry (minutes) | 30           |
| `--refresh-token-expiry-hour`, `REFRESH_TOKEN_EXPIRY_HOUR` | Refresh token expiry (hours)  | 24           |

---

## Project Structure

- `cmd/zerodupe-server/` — Server entry point
- `cmd/zerodupe-client/` — Client entry point
- `internal/server/` — Server logic and API
- `pkg/client/` — Client logic and API
- `pkg/hasher/` — Hashing utilities

---

## Development

To build and run locally (requires Go):

```bash
go build -o zerodupe-server ./cmd/zerodupe-server/main.go
go build -o zerodupe-client ./cmd/zerodupe-client/main.go
```

## Running Tests

```bash
go test ./...
```

---

## Summary Table

| What you want to do | Command Example                                                                           |
| ------------------- | ----------------------------------------------------------------------------------------- |
| Start the server    | `docker-compose up -d zerodupe-server`                                                    |
| Sign up a user      | `docker-compose run --rm zerodupe-client signup --server http://zerodupe-server:8080 ...` |
| Upload a file       | `docker-compose run --rm -v $(pwd)/file.txt:/app/file.txt zerodupe-client upload ...`     |
| Download a file     | `docker-compose run --rm -v $(pwd)/downloads:/app/downloads zerodupe-client download ...` |
| Stop everything     | `docker-compose down`                                                                     |
