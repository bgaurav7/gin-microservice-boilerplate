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

# Runtime stage
FROM alpine:3.19

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from the builder stage
COPY --from=builder /go/bin/server /server

# Copy configuration files
COPY --from=builder /go/src/app/config /config
COPY --from=builder /go/src/app/migrations /migrations

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["/server"]
