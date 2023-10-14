FROM golang:1.21.2-alpine AS builder

COPY . /app
WORKDIR /app

RUN go mod download
RUN go build -o ./bin/grpc-server ./cmd/grpc-server/main.go

FROM alpine:latest

WORKDIR /root
COPY --from=builder /app/bin/grpc-server .

CMD ["./grpc-server"]