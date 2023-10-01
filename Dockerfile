# Use the official Golang base image
FROM golang:1.18-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code into the container
COPY . .

# Build the Go application
RUN go build ./cmd/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the Go application
CMD ["./main"]