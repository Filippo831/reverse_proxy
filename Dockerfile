FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN go build -o /reverse_proxy ./cmd/reverse_proxy

FROM alpine:latest

COPY --from=builder /reverse_proxy /reverse_proxy

EXPOSE 8081
EXPOSE 8082

ENTRYPOINT ["/reverse_proxy"]

