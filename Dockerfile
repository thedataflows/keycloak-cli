# Build stage
FROM goreleaser/goreleaser:v2.16.0 AS builder

# Set working directory
WORKDIR /app

# Copy source code and goreleaser config
COPY . .

# Build the application using goreleaser
# Use --snapshot for local builds without git tags
# Filter for linux/amd64 only to speed up Docker build
ARG VERSION=v0.0.0-snapshot
ENV GORELEASER_CURRENT_TAG=$VERSION
RUN goreleaser build --snapshot --clean --single-target

# Final stage
FROM alpine:3.22

# Install ca-certificates for HTTPS requests
RUN apk add --no-cache ca-certificates

# Create non-root user
RUN addgroup -g 1000 restricted && \
    adduser -D -s /bin/sh -u 1000 -G restricted restricted

# Copy binary from builder stage
COPY --from=builder /app/dist/default_linux_amd64_v1/keycloak-cli /usr/local/bin/

USER restricted

WORKDIR /tmp

ENTRYPOINT ["keycloak-cli"]
