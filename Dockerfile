# Build stage
FROM golang:1.23.6-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dapr-actor-gen ./cmd

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/dapr-actor-gen .

# Make the binary executable
RUN chmod +x ./dapr-actor-gen

# Create a non-root user
RUN adduser -D -s /bin/sh appuser
USER appuser

ENTRYPOINT ["./dapr-actor-gen"]