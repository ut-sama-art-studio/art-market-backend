# Use the official Golang image as the base for building the app
FROM golang:1.20 AS builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code directories into the container
COPY ./database ./database
COPY ./server ./server
COPY ./graph ./graph
COPY ./services ./services
COPY ./utils ./utils
COPY ./middlewares ./middlewares
COPY ./tests ./tests
COPY ./migrations ./migrations

# Ensure dependencies from internal packages and external modules are downloaded
RUN go get -d ./...

# Build the application
RUN go build -o /app/server-binary ./server/*.go

# Runtime image
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/server-binary .

# Expose the API port
EXPOSE 8080

# Start the API server
CMD ["./server-binary"]


# building container
#   docker build -t art-market-api .

# running with .env mounted at runtime
#   docker run --network art-market-network -v -p 8080:8080 --name art-market-api art-market-api:latest  