FROM golang:1.23-alpine AS builder

COPY . /github.com/ipv02/auth/source/
WORKDIR /github.com/ipv02/auth/source/

RUN go mod download
RUN go build -o ./bin/auth_crud_server cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/ipv02/auth/source/bin/auth_crud_server .

CMD ["./auth_crud_server"]