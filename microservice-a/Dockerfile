FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY microservice-a/go.mod microservice-a/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY microservice-a/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o microservice-a ./cmd/main.go

# Create a minimal image
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/microservice-a .

# Expose ports
EXPOSE 8090

# Run the application
CMD ["./microservice-a"]