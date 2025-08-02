# -- Build the binary
FROM golang:1.24.3-alpine3.21 AS builder

# Install build dependencies for CGO and SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Enable CGO and build the server binary
ENV CGO_ENABLED=1
RUN go build -o zerodupe-server ./cmd/zerodupe-server/main.go

# -- Create minimal image
FROM alpine:3.21

WORKDIR /server

COPY --from=builder /app/zerodupe-server .
COPY --from=builder /app/internal/server/storage /server/internal/server/storage

# Install runtime SQLite library
RUN apk add --no-cache sqlite-libs

CMD ["./zerodupe-server", "--port", "8080", "--secret", "supersecret"]
