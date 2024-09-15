# Use Go base image
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first
COPY go.mod go.sum ./

# Download dependencies. Dependencies will be cached if the go.mod and go.sum files haven't changed.
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go app
RUN go build -o todo-api

# Expose the API port
EXPOSE 8080

# Command to run the API
CMD ["./todo-api"]
