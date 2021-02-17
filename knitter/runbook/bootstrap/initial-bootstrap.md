# Initial Bootstrap for zLifecycle in a totally new setup

## Overview

This is used for initial bootstrap of zLifecycle in a totally new setup. Should only be needed once.

## When to use this runbook
This is used for initial bootstrap of zLifecycle in a totally new setup. Should only be needed once.

## Initial Steps Overview

1. [Create AWS Service Account](#create-aws-service-account)
1. [Setup zLifecycle Terraform Shared State](#setup-zLifecycle-terraform-shared-state)
1. [Create ECR Repositories](#create-ecr-repositories)

## Detailed Steps

#### Create AWS Service Account
Create an AWS IAM account that can be used as service account by Terraform to provision resources

For now give it admin access but in future it should have RBAC to limit the surface area

Add following to ~/.aws/credentials file locally:

```
[zlifecycle-shared]
aws_access_key_id = [access key if here]
aws_secret_access_key = [secret access key here]
```

#### Setup zLifecycle Terraform Shared State
zLifecycle environments (e.g. demo, dev) are managed by terraform workspaces. These terraform workspaces share a parent state directory maintained in terraform (`zlifecycle-tfstate`) that needs to be initialized before environments can be created. This bootstrap script is for this use case, where no zlifecycle environments exist yet.

Run `tfstate` terraform to provision S3 bucket and Dynamo DB table that will be used for Terraform Shared State.

```bash
cd tfstate
terraform init
terraform apply
```

Commit & push the tfstate file to github since there is no shared state bucket yet.

#### Create ECR Repositories

Manually create ECR repositories in AWS:
`zlifecycle-il-operator`
`zlifecycle-terraform`

