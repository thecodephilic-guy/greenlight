# Stage 1: Build the binary
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git and ca-certificates (needed for fetching dependencies and HTTPS calls)
RUN apk add --no-cache git ca-certificates

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build a statically linked binary for Linux amd64 with git hash/tag and build time
RUN VERSION=$(git describe --always --dirty --tags --long 2>/dev/null || echo "unknown") && \
    BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ) && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}" \
    -o /app/bin/api ./cmd/api

# Stage 2: Final lightweight runtime container
FROM alpine:3.20

WORKDIR /app

# Add ca-certificates so your app can make outbound HTTPS requests (e.g. to Resend or Neon)
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user for security
RUN adduser -D -g '' appuser
USER appuser

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/bin/api /app/api

# The port Greenlight listens on
EXPOSE 4000

# Run the app in production mode
ENTRYPOINT ["/app/api"]
CMD ["-port=4000", "-env=production"]
