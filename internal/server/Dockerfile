# -- Build the binary
FROM golang:1.24.3-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o zerodupe-server cmd/zerodupe-server/main.go

# -- Create minimal image
FROM alpine:3.21

WORKDIR /server

COPY --from=builder /app/zerodupe-server .

CMD ["./zerodupe-server"]
