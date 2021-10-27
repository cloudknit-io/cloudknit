# syntax=docker/dockerfile:1

# BUILDER IMAGE
FROM golang:1.16-alpine as builder

# source files
COPY app /svc/app
COPY main.go /svc/main.go
# mod files
COPY go.mod /svc/go.mod
COPY go.sum /svc/go.sum

WORKDIR /svc

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o zlifecycle-state-manager main.go

# APPLICATION IMAGE
FROM alpine

ARG TERRAFORM_VERSION=0.14.10

RUN apk add --update --no-cache git openssh

# add known hosts
RUN mkdir -p ~/.ssh \
    && touch ~/.ssh/known_hosts \
    && ssh-keyscan -t rsa github.com >> /root/.ssh/known_hosts \
    && ssh-keyscan -t rsa gitlab.com >> /root/.ssh/known_hosts \
    && ssh-keyscan -t rsa bitbucket.org >> /root/.ssh/known_hosts

# install terraform
RUN wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
    && unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
    && rm terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
    && mv terraform /usr/bin/terraform

WORKDIR /svc
COPY --from=builder /svc/zlifecycle-state-manager /svc/zlifecycle-state-manager

EXPOSE 8080

ENTRYPOINT [ "/svc/zlifecycle-state-manager" ]
