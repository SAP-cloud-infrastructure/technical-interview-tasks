# Build stage
FROM golang:1.22-bullseye AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o simple-web-app .

# Runtime stage
FROM debian:bullseye-slim

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary from builder
COPY --from=builder /app/simple-web-app .

# Expose port
EXPOSE 8080

# Run the application
ENTRYPOINT ["./simple-web-app"]
