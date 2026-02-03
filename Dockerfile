# Use the official Golang image as the base image
FROM golang:1.25.5 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the source code
COPY . .

# Download dependencies
RUN cd src && go mod download

# Build the Go application
RUN cd src && go build -o main ./cmd/prod/main.go

FROM ubuntu:22.04
WORKDIR /root/

# Install ca-certificates
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the built binary from the builder stage
COPY --from=builder ./app/src/main .

EXPOSE 8080

# Run the binary
CMD ["./main"]