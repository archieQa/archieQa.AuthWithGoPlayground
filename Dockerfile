# Use the official Golang image
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Install necessary tools
RUN apk add --no-cache git

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Set environment variables for Go
ENV GO_ENV=development

# Command to run the application without building
CMD ["go", "run", "main.go"]