FROM golang:1.24-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

COPY . .

# Build the Go app
RUN go build -o yugi ./cmd/yugi

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go build`
CMD ["./yugi"]