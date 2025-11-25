FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY *.go ./


RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]
