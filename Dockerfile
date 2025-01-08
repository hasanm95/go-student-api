# Build stage
FROM golang:1.23.4-alpine3.20

# Set working directory
WORKDIR /app

# Install required dependencies including SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Install air
RUN go install github.com/air-verse/air@latest

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Create data directory for SQLite
RUN mkdir -p /app/data

# Build the application
RUN go build -o main main.go 

# Expose port
EXPOSE 8080

# Command to run the application
CMD ["air", "-c", ".air.toml"]