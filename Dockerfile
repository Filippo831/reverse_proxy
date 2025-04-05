FROM golang

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /reverse_proxy

EXPOSE 8081 8082

CMD ["/reverse_proxy"]
