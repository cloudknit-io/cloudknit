## Installation of AWS and K8s Locally

# Overview
Installation and setup of aws and kubernetes cli if not already done.

# When to use this runbook
If your machine doesnt have aws and k8s cli installed and configured

# Steps

- Install [aws-cli](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html), [kubernetes cli](https://kubernetes.io/docs/tasks/tools/)
- Setup AWS credentials, by adding following keys in the credentials file present at ~/.aws/credentials location

    `[default]`\
    `aws_access_key_id = {creds-provided-via-last-pass}`\
    `aws_secret_access_key = {creds-provided-via-last-pass}`\
    `[compuzest-shared]`\
    `aws_access_key_id = {creds-provided-via-last-pass}`\
    `aws_secret_access_key = {creds-provided-via-last-pass}`
    
- Add config information for kubernetes by running following command `aws eks --region {region} update-kubeconfig --name {environment-name}`
