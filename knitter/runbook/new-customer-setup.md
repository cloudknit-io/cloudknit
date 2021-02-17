# Setup a new Customer

## Overview

Steps to setup a new Customer

## When to use this runbook
This is to be used when you are setting up brand new customer

## Initial Steps Overview

1. [Setup Github Service Account](#setup-github-service-account)
1. [Setup Customer Terraform Shared State](#setup-customer-terraform-shared-state)

## Detailed Steps

#### Setup Github Service Account
This account will be used by zLifecycle to Read and Write to various repos:

- Create new Github service account (example: zLifecycle with zLifecycle@compuzest.com email)
- Add new github service account to the customer github org and give perms to following repos
  - compuzest-zlifecycle-il - write access
  - helm-charts - read access
  - compuZest-zlifecycle-config - read access
- Generate Personal Token for the Github Service account to be used by secret created (Check LastPass secret note: "zLifecycle - k8s secrets")
  - In the scope select all options for repo and workflow
- Generate ssh key for the Github service account to be used by secret created (Check LastPass secret note: "zLifecycle - k8s secrets")

#### Setup Customer Terraform Shared State
Each Customer has a separate S3 bucket and Dynamo DB for their shared terraform state files so they are isolated from other customers state file.

Run `tfstate` terraform to provision S3 bucket and Dynamo DB table that will be used for Terraform Shared State for the Customer.

- Create a folder with companies name under `tfstate/customers`
- Copy files from compuzest directory and change names & s3 & Dynamo DB location for state
- Run following:
```bash
cd tfstate
terraform init
terraform apply
```
