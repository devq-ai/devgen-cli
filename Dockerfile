# DevGen CLI - Multi-stage Docker build for optimized production image
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o devgen .

# Production stage
FROM scratch

# Copy ca-certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /app/devgen /usr/local/bin/devgen

# Copy example files
COPY --from=builder /app/example-playbook.yaml /workspace/example-playbook.yaml
COPY --from=builder /app/README.md /workspace/README.md

# Set working directory
WORKDIR /workspace

# Create non-root user
USER 1000:1000

# Set environment variables
ENV DEVGEN_CONFIG_DIR=/workspace/.devgen
ENV DEVGEN_LOG_LEVEL=info
ENV TERM=xterm-256color

# Expose port for server mode
EXPOSE 8080

# Default command
ENTRYPOINT ["/usr/local/bin/devgen"]
CMD ["--help"]

# Metadata
LABEL maintainer="DevQ.ai Team <dion@devq.ai>"
LABEL description="DevGen CLI - Development Generation Tool with Charm UI"
LABEL version="1.0.0"
LABEL org.opencontainers.image.source="https://github.com/devq-ai/devgen-cli"
LABEL org.opencontainers.image.documentation="https://github.com/devq-ai/devgen-cli/blob/main/README.md"
LABEL org.opencontainers.image.licenses="MIT"
