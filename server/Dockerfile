# Stage 1: Build the Go binary
FROM golang:1.24 AS builder
WORKDIR /app

# Copy source files into the container
COPY . .

# Build the server binary for Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Stage 2: Create a minimal container image
FROM alpine:latest
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/server .

# Expose the port that the server listens on
EXPOSE 8080

# Run the server binary
CMD ["./server"]