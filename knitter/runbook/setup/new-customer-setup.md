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

- Create a mailing group for `<client>@compuzest.com` on G Suite
- Create new Github Service Account and register it under `<client>@compuzest.com`, username should follow the format `<client>-zlifecycle`
- Run the `bin/init_customer.sh` script to create an IL repo and assign permissions to the Zlifecycle service account
    ```shell script
    bin/create_github_repo.sh -t <token> -c <customer>
    ```
- Generate Personal Access Token for the Zlifecycle Service Account to be used by the Operator (Check LastPass secret note: "zLifecycle - k8s secrets")
  - In the scope select all options for `repo` and `workflow`
- Generate SSH key for the Github Service Account to be used by the Operator (Check LastPass secret note: "zLifecycle - k8s secrets")
    ```shell script
    ssh-keygen -b 2048 -t rsa -f <folder/to/generate/key> -q -N "" -C "<client>@compuzest.com"
    ```
- Add zLifecycle as an OAuth application by going `Repository Settings -> Developer Settings -> OAuth Apps -> New OAuth App`
    * Application name: `zLifecycle-<client>`
    * Homepage URL: `https://<client>-admin.zlifecycle.com`
    * Application description (OPTIONAL): `zLifecycle instance for <client>`
    * Authorization callback URL: `https://<client>-admin.zlifecycle.com/api/dex/callback`
- Generate a new client secret from the Application OAuth page
    
#### Setup Customer Terraform Shared State
Each Customer has a separate S3 bucket and Dynamo DB for their shared terraform state files, so they are isolated from other customers state file.

Run `tfstate` terraform to provision S3 bucket and Dynamo DB table that will be used for Terraform Shared State for the Customer.

- Create a folder with companies name under `s3-bucket/customers`
- Copy files from `zmart` directory and change names & s3 & Dynamo DB location for state (IMPORTANT: remove preexisting `tfstate` file)
- Run following:
```bash
cd s3-bucket/customers/<client>
terraform init
terraform apply
```
