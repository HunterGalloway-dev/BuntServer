# syntax=docker/dockerfile:1

FROM golang:1.23.0

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY . .

# Build
RUN go build -o bin/docker-buntserver ./cmd/api/main.go

RUN apt-get update
RUN apt-get install nodejs -y

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080

# Run
CMD ["./bin/docker-buntserver"]