FROM golang:1.24.3-alpine3.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o zerodupe-client cmd/zerodupe-client/main.go

ENTRYPOINT ["./zerodupe-client"]