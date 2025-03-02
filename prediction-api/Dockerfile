# Start with the official Golang image to compile the application
FROM golang:1.23 AS build

# Set the working directory in the container
WORKDIR /api

# Copy all files
COPY ./prediction-api/ ./

# Download all dependencies
RUN go mod download

# Build the application as a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o api

# Expose port 8080 to traffic
EXPOSE 8080

# Define the default command to run the service
CMD ["./api"]
