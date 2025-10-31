
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY microservice-b/go.mod microservice-b/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY microservice-b/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o microservice-b ./cmd/main.go

# Create a minimal image
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/microservice-b .

# Expose ports
EXPOSE 8080 50051

# Run the application
CMD ["./microservice-b"]