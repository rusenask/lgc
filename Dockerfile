# golang image where workspace (GOPATH) configured at /go.
FROM golang:1.5-onbuild

# Copy the local package files to the container’s workspace.
ADD . /go/src/github.com/rusenask/lgc

# Active vendor experiment
RUN export GO15VENDOREXPERIMENT=1
# Build the LGC command inside the container.
RUN go install github.com/rusenask/lgc

# Run the lgc command when the container starts.
ENTRYPOINT /go/bin/lgc

# http server listens on port 3000.
EXPOSE 3000
