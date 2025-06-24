# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:latest
WORKDIR /app

# Install bash
RUN apk add --no-cache bash

COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY start.sh /app/start.sh
COPY wait-for-it.sh /app/wait-for-it.sh
RUN chmod +x /app/start.sh /app/wait-for-it.sh

EXPOSE 8080
ENTRYPOINT ["/app/start.sh"]