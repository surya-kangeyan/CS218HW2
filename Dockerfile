# Use the official Golang image to build app
FROM golang:1.17-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o sys-monitor .

# Use a minimal image for running the app
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/sys-monitor .

# Expose the application's port
EXPOSE 8080

# Command to run the executable
CMD ["./sys-monitor"]

