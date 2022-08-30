# Copyright (C) 2020 CompuZest, Inc. - All Rights Reserved
#
# Unauthorized copying of this file, via any medium, is strictly prohibited
# Proprietary and confidential
#
# NOTICE: All information contained herein is, and remains the property of
# CompuZest, Inc. The intellectual and technical concepts contained herein are
# proprietary to CompuZest, Inc. and are protected by trade secret or copyright
# law. Dissemination of this information or reproduction of this material is
# strictly forbidden unless prior written permission is obtained from CompuZest, Inc.

FROM alpine:latest AS aws-iam-authenticator

RUN apk add --no-cache --virtual .build-deps curl

RUN curl -o /tmp/aws-iam-authenticator https://amazon-eks.s3-us-west-2.amazonaws.com/1.14.6/2019-08-22/bin/linux/amd64/aws-iam-authenticator && \
    chmod +x /tmp/aws-iam-authenticator && \
    install -b -c -v /tmp/aws-iam-authenticator /usr/local/bin/ && \
    aws-iam-authenticator version

RUN apk del .build-deps && rm /tmp/*

FROM hashicorp/terraform:1.0.9

ARG INFRACOST_API_KEY
ARG ZLI_TAG
ARG GITHUB_TOKEN

RUN apk add --update --no-cache curl jq bash openssl python3 py3-pip py3-setuptools libc6-compat unzip
RUN python3 -m pip --no-cache-dir install --upgrade awscli

RUN curl -sL https://github.com/infracost/infracost/releases/download/v0.10.10/infracost-linux-amd64.tar.gz | tar xz -C /tmp && \
    mv /tmp/infracost-linux-amd64 /usr/local/bin/infracost && \
    mkdir ~/.config && \
    mkdir ~/.config/infracost/

# Replace infracost_key
RUN echo $'https://pricing.api.infracost.io: \n\
  api_key: '$INFRACOST_API_KEY >> ~/.config/infracost/credentials.yml

RUN curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl" && \
    chmod +x ./kubectl && \
    mv ./kubectl /usr/local/bin/kubectl && \
    kubectl version --client

RUN VERSION=$(curl --silent "https://api.github.com/repos/argoproj/argo-cd/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/') && \
    curl -sSL -o /usr/local/bin/argocd https://github.com/argoproj/argo-cd/releases/download/$VERSION/argocd-linux-amd64 && \
    chmod +x /usr/local/bin/argocd

COPY ./scripts .
COPY ./helpers ./helpers

RUN find / -type f -iname "*.sh" -exec chmod +x {} \; && \
    find /argocd -type f -iname "*.sh" -exec chmod +x {} \; && \
    find /client -type f -iname "*.sh" -exec chmod +x {} \;

COPY --from=aws-iam-authenticator /usr/local/bin/aws-iam-authenticator /usr/local/bin/aws-iam-authenticator

RUN curl -fsSL -o /tmp/get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 && \
    chmod 700 /tmp/get_helm.sh && \
    /tmp/get_helm.sh && \
    helm plugin install https://github.com/hypnoglow/helm-s3.git && \
    rm /tmp/get_helm.sh

# get the zlifecycle-internal-cli
RUN ./helpers/gh_download.sh -v $ZLI_TAG -f zlifecycle-internal-cli_${ZLI_TAG#v}_linux_x86_64.zip -t $GITHUB_TOKEN -r zlifecycle/zlifecycle-internal-cli -o /tmp/cli.zip && \
    unzip -d /tmp /tmp/cli.zip && \
    install -b -c -v /tmp/zlifecycle-internal-cli /usr/local/bin/ && \
    zlifecycle-internal-cli -h && \
    rm /tmp/cli.zip && \
    rm /tmp/zlifecycle-internal-cli

# Adding gitlab.com to known_hosts to avoid "Host key verification failed." error
RUN mkdir ~/.ssh && \
    touch ~/.ssh/known_hosts && \
    ssh-keyscan -t rsa gitlab.com >> ~/.ssh/known_hosts && \
    ssh-keyscan -t ed25519 gitlab.com >> ~/.ssh/known_hosts && \
    ssh-keyscan -t ed25519 github.com >> ~/.ssh/known_hosts

ENTRYPOINT ["sh"]
