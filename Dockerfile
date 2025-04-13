FROM alpine:3.21

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary built by GoReleaser
COPY buildkite-mcp-server /app/

# Set the entrypoint to run the server in stdio mode
ENTRYPOINT ["/app/buildkite-mcp-server", "stdio"]
