# syntax=docker/dockerfile:1

# BUILDER IMAGE
FROM golang:1.17-alpine as builder

# source files
COPY app /svc/app
COPY main.go /svc/main.go
# mod files
COPY go.mod /svc/go.mod
COPY go.sum /svc/go.sum

WORKDIR /svc

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o zlifecycle-state-manager main.go

# APPLICATION IMAGE
FROM alpine
ARG TERRAFORM_VERSION=1.0.9
RUN wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
    && unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
    && rm terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
    && mv terraform /usr/bin/terraform

WORKDIR /svc
COPY --from=builder /svc/zlifecycle-state-manager /svc/zlifecycle-state-manager

EXPOSE 8080

ENTRYPOINT [ "/svc/zlifecycle-state-manager" ]
