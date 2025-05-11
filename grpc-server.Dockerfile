FROM golang:1.24.1-alpine AS builder

COPY . /github.com/nogavadu/auth-service
WORKDIR /github.com/nogavadu/auth-service

RUN go mod download
RUN go build -o ./bin/auth-service cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/nogavadu/auth-service/bin/auth-service .

CMD ["./auth-service"]