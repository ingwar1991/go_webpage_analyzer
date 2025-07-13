# Stage 1: Build the Go binary
FROM golang:1.23.9-alpine AS builder

# Set environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go binary
RUN go build -o server ./cmd/web

# Stage 2: Create a minimal runtime image
FROM alpine:latest

# Set working directory in runtime container
WORKDIR /app

# Copy binary and templates to the runtime container
COPY --from=builder /app/server .
COPY --from=builder /app/web/templates ./web/templates

# Set the entrypoint
ENTRYPOINT ["./server"]

# Expose the port your server listens on (optional, but helpful)
EXPOSE 8080
