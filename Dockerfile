# Build stage
FROM golang:1.24-alpine AS builder

# Install dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /go/src/app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/server ./cmd/server

# Intermediate stage for preparing files with proper permissions
FROM alpine:3.19 as prepare

# Create app directories
RUN mkdir -p /app/config /app/migrations /app/internal/infrastructure/rbac

# Copy the binary and config files
COPY --from=builder /go/bin/server /app/server
COPY --from=builder /go/src/app/config /app/config
COPY --from=builder /go/src/app/migrations /app/migrations
COPY --from=builder /go/src/app/internal/infrastructure/rbac /app/internal/infrastructure/rbac

# Set permissions for nonroot user (uid 65532)
RUN chmod +x /app/server && \
    chown -R 65532:65532 /app

# Runtime stage
FROM gcr.io/distroless/static:nonroot

# Copy prepared files with correct permissions
COPY --from=prepare /app /app

# Set working directory
WORKDIR /app

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["/app/server"]
