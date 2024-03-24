# Use a smaller base image for the server
FROM --platform=linux/amd64  golang:1.21-alpine3.19 AS builder

# Set the working directory
WORKDIR /build

# Copy the source files
COPY . ./

# downloading modules 
RUN go mod download

# Build the server app
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o server ./cmd/main.go


# Run from scratch
FROM --platform=linux/amd64 scratch 

WORKDIR /app 

COPY --from=builder /build/server /app/server

EXPOSE 10015

CMD ["/app/server"]
