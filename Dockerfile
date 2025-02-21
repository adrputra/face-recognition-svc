# Stage 1: Build the Go application
FROM golang:1.21-alpine AS builder

# Install git (required for go mod)
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the Go application source code into the container
COPY . .

# Build the Go application inside the container
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/myapp .

# Stage 2: Create a small final image
FROM alpine:latest

WORKDIR /app

# Copy the compiled Go application from the builder stage
COPY --from=builder /app/myapp /app/myapp

# Copy the configuration file from the builder stage (if needed)
COPY --from=builder /app/config.yaml /app/config.yaml

# Expose the port that the application listens on
EXPOSE 8001

# Command to run the application
CMD ["./myapp"]
