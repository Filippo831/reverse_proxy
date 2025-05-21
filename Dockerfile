FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/reverse_proxy/main.go ./main.go
COPY internal/ ./internal/

RUN go build -o /reverse_proxy 

FROM alpine:latest

COPY --from=builder /reverse_proxy /reverse_proxy

EXPOSE 8081
EXPOSE 8082

ENTRYPOINT ["/reverse_proxy"]

