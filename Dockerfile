# Build stage
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates (needed for go mod download)
RUN apk add --no-cache git ca-certificates tzdata

# Create appuser for security
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o pulse \
    ./cmd/server/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create appuser for security
RUN adduser -D -g '' appuser

# Create data directory for SQLite
RUN mkdir -p /app/data && chown appuser:appuser /app/data

# Copy the binary
COPY --from=builder /build/pulse /app/pulse

# Set working directory
WORKDIR /app/data

# Use appuser
USER appuser

# Expose port
EXPOSE 8080

# Health check (simple TCP check since we don't have health endpoint yet)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD timeout 3s sh -c 'nc -z localhost 8080' || exit 1

# Run the binary
ENTRYPOINT ["/app/pulse"]
