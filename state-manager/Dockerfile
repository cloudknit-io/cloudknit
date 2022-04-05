# syntax=docker/dockerfile:1
FROM alpine

ARG TERRAFORM_VERSION=1.0.9

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

COPY --from=builder zlifecycle-state-manager /

EXPOSE 8080

ENTRYPOINT [ "/zlifecycle-state-manager" ]
