FROM golang:1.19-alpine

# Set destination for COPY
WORKDIR /app

RUN apk add --no-cache git

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY internal ./internal
COPY cmd ./cmd

WORKDIR /app/cmd
RUN go clean -cache

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080
# Build
RUN go run .
