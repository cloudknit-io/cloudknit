# Index
- [Overview](#overview)
- [Setup New Customer](#setup-new-customer)
- [Initial Bootstrap](#initial-bootstrap)
- [Bootstrap zLifecycle](#bootstrap-zlifecycle)
- [Register Teams](#register-teams)
- [Build Terraform Docker Image](build-terraform-docker-image)

# Overview

zLifecycle is a product to manage lifecycle for infrastructure across various cloud providers as well as on-prem.

For more details & diagrams look at: https://app.diagrams.net/#G1gXeFRlERpqjXpeSjxRPLP6YZMRyFG5SN

## Setup New Customer

* Create new Github service account (example: zLifecycle with zLifecycle@compuzest.com email)
* Generate Personal Token & ssh key to be used by secret created (Check LastPass secret note: "zLifecycle - k8s secrets")
* Add new github service account to the customer github org and give perms to following repos
    * terraform-environment - write access
    * helm-charts - read access
    * zLifecycle-CompuZest-config - read access

## Initial Bootstrap

### Terraform Shared State

Run `tfstate` terraform to provision S3 bucket and Dynamo DB table that will be used for Terraform Shared State.

```bash
cd tfstate
terraform init
terraform apply
```

## Bootstrap zLifecycle

To bootstrap zLifecycle run following script and following instructions:

Note: When it asks to create secret go to `infra-deploy-platform/k8s-addons/argo-workflow` folder 
and create secrets using scripts in LastPass

```bash
cd bootstrap
./bootstrap_zLifecycle.sh
```

## Register Teams

You need to manually register teams currently using following script

```bash
cd ../../zLifecycle-CompuZest-config
kubectl apply -f teams/account-team.yaml # Replace yaml file with team name for the team you want to register
```

## Build Terraform Docker Image

```bash
cd terraform
docker build -t shahadarsh/terraform:latest .
docker push shahadarsh/terraform:latest
```
