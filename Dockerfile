# syntax=docker/dockerfile:1
# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/app/main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /esp32_prometheus_exporter

# Runtime stage
FROM alpine:3.20
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /esp32_prometheus_exporter /esp32_prometheus_exporter
EXPOSE 2112
ENTRYPOINT ["/esp32_prometheus_exporter"]
