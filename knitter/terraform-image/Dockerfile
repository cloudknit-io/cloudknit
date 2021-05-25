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

RUN curl -o /tmp/aws-iam-authenticator https://amazon-eks.s3-us-west-2.amazonaws.com/1.14.6/2019-08-22/bin/linux/amd64/aws-iam-authenticator
RUN chmod +x /tmp/aws-iam-authenticator
RUN install -b -c -v /tmp/aws-iam-authenticator /usr/local/bin/
RUN aws-iam-authenticator version

RUN apk del .build-deps
RUN rm /tmp/*

FROM hashicorp/terraform:0.13.2
RUN apk add --update curl

RUN apk update && apk upgrade && apk --no-cache add curl jq

RUN curl -sL https://github.com/infracost/infracost/releases/download/v0.8.6/infracost-linux-amd64.tar.gz | tar xz -C /tmp

RUN mv /tmp/infracost-linux-amd64 /usr/local/bin/infracost

RUN mkdir ~/.config && mkdir ~/.config/infracost/

# Replace infracost_key
RUN echo $'https://pricing.api.infracost.io: \n\
  api_key: ${infracost_key}' >> ~/.config/infracost/credentials.yml

RUN curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl"
RUN chmod +x ./kubectl

RUN mv ./kubectl /usr/local/bin/kubectl
RUN kubectl version --client

RUN VERSION=$(curl --silent "https://api.github.com/repos/argoproj/argo-cd/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/') && curl -sSL -o /usr/local/bin/argocd https://github.com/argoproj/argo-cd/releases/download/$VERSION/argocd-linux-amd64
RUN chmod +x /usr/local/bin/argocd

COPY ./scripts .

RUN find / -type f -iname "*.sh" -exec chmod +x {} \;
RUN find /argocd -type f -iname "*.sh" -exec chmod +x {} \;
RUN find /client -type f -iname "*.sh" -exec chmod +x {} \;
RUN find /terraform -type f -iname "*.sh" -exec chmod +x {} \;

COPY --from=aws-iam-authenticator /usr/local/bin/aws-iam-authenticator /usr/local/bin/aws-iam-authenticator

ENTRYPOINT ["sh"]
