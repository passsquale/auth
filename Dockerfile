FROM golang:1.22.3-alpine AS builder

COPY . /github.com/passsquale/auth/source/
WORKDIR /github.com/passsquale/auth/source/

RUN go mod download
RUN go build -o ./bin/auth cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/passsquale/auth/source/bin/auth .

CMD ["./auth"]