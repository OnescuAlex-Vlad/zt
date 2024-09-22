# Use official Golang image as base
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the project files into the container
COPY . .

# Build the Go application
RUN go build -o /corebit cmd/corebit/main.go

# Set the command to run the application
CMD ["/corebit"]
