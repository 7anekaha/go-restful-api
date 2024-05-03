# Use the official golang image from Docker Hub
FROM golang:1.22

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download and install Go dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o myapp .

# Expose the port on which the Go application will run
EXPOSE 3000

# Command to run the executable
CMD ["./myapp"]
