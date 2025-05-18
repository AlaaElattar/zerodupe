FROM golang:1.24.3-alpine3.21

WORKDIR /server

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o zerodupe-server cmd/zerodupe-server/main.go

CMD ["./zerodupe-server"]
