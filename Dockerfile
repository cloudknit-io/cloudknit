# Build the manager binary
FROM golang:1.17 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/
COPY templates/ templates/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager main.go

# add known hosts
RUN mkdir -p /.ssh \
    && touch /.ssh/known_hosts \
    && ssh-keyscan -t rsa github.com >> /.ssh/known_hosts \
    && ssh-keyscan -t rsa gitlab.com >> /.ssh/known_hosts \
    && ssh-keyscan -t rsa bitbucket.org >> /.ssh/known_hosts

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/base:debug
# USER nonroot:nonroot

WORKDIR /
COPY --from=builder /workspace/manager .
COPY --from=builder /workspace/templates ./templates

COPY --from=builder /.ssh /root/.ssh

ENTRYPOINT ["/manager"]
