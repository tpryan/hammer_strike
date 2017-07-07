# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

RUN apt-get update && apt-get -y install apache2-utils

# Copy the local package files to the container's workspace.
ADD . /go/src/abrunner
RUN mkdir /go/src/abrunner/logs

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go install abrunner

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/abrunner

# Document that the service listens on port 8080.
EXPOSE 80