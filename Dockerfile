FROM golang:1.23-alpine

# Set working directory
WORKDIR /app

# Copy go modules & install
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the app
RUN go build -o app ./cmd/main.go

# Expose port
EXPOSE 8080

# Run the app
CMD ["./app"]
