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

COPY --from=aws-iam-authenticator /usr/local/bin/aws-iam-authenticator /usr/local/bin/aws-iam-authenticator

ENTRYPOINT ["sh"]
