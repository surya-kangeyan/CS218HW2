# Stage 1: Build the Go application
FROM golang:1.17-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to leverage Docker cache
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o sys-monitor .

# Stage 2: Create a minimal image for running the application
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/sys-monitor .

# Expose the port the application listens on
EXPOSE 8080

# Command to run the executable
CMD ["./sys-monitor"]
