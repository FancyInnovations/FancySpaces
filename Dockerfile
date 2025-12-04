# Use the official Golang image as the base image
FROM golang:1.25.5 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the source code
COPY . .

# Download dependencies
RUN go mod download

# Build the Go application
RUN go build -C fancyspaces -o main ./src/cmd/local/main.go

FROM ubuntu:22.04
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder ./app/fancyspaces/main .

EXPOSE 8080

# Run the binary
CMD ["./main"]