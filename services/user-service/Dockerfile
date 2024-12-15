FROM golang:latest

WORKDIR /app

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Download the Go module dependencies
RUN go mod download

COPY . .

RUN go build -o /myapp ./cmd

EXPOSE 8082
CMD ["/myapp"]