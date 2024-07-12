FROM golang:1.22-alpine as builder

LABEL authors="zvilutis"

ENTRYPOINT ["top", "-b"]

# Create and set the working directory
WORKDIR /app

# Copy the project's Go module dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the project source code
COPY . .

# Build the Go application
RUN go build cmd/store/*.go

# Use a smaller Alpine image as the final stage
FROM alpine:latest

# Copy the built binary from the previous stage
COPY --from=builder /app/main /app/main

# Expose the application port (adjust as needed)
EXPOSE 8080

# Specify the command to run your application
CMD ["/app/main"]
