# Build the manager binary
FROM golang:1.13 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

COPY controller/main.go controller/main.go
COPY controller/api/ controller/api/
COPY controller/controllers/ controller/controllers/

COPY cni/ cni/
COPY pkg/ pkg/

# Build ipam
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o bifrost-ipam cni/ipam/main.go cni/ipam/dns.go

# Build controller
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager controller/main.go
