FROM golang:1.22.3

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.sum go.mod ./
COPY internal/ internal/
COPY cmd/ cmd/
COPY ../configs configs/
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o consumer_service ./cmd/consumer_service.go

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 9001
