# Build stage
FROM golang:1.23.4-alpine3.20

# Set working directory
WORKDIR /app

# Install air
RUN go install github.com/air-verse/air@latest

# Copy go mod and sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main main.go 


# Expose port
EXPOSE 8080

# Command to run the application
CMD ["air", "-c", ".air.toml"]
