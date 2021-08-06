# Create SSM secret

## Overview

For customer secrets, we currently use Parameter Store on AWS Systems Manager

## When to use this runbook

When you want to create a customer secret

## Prerequisites

1. [aws-cli](https://github.com/aws/aws-cli) - AWS CLI

## Initial Steps Overview

- [Create a secret](#create-a-secret)

## Detailed Steps

### Create a secret
Secret scopes:
* organization secret: `/<customer>/<secret_name>`
* team secret: `/<customer>/<team>/<secret_name`
* environment secret: `/<customer>/<team>/<environment>/<secret_name>`
* environment component secret: `/<customer>/<team>/<environment>/<environment_component>/<secret_name>`

1. Run the following command to create a SecureString secret in Parameter Store on AWS Systems Manager
```shell script
region="us-east-1"
# if using default profile, remove the --profile line, otherwise uncomment the line below
#profile=<customer_aws_profile>
value=<secret>
# key should be in a path format (depends on scope)
key=<key>
aws ssm put-parameter \
  --region "$region" \
  --profile "$profile" \
  --type SecureString \
  --name $key \
  --value "$secret" > /dev/null
```
