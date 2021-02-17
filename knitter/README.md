# Index
- [Overview](#overview)
- [Runbooks](#runbooks)
- [Initial Bootstrap](./runbook/setup/initial-bootstrap.md)
- [Setup New Customer](./runbook/setup/new-customer-setup.md)
- [Bootstrap zLifecycle](./runbook/setup/bootstrap-zlifecycle.md)
- [Register Teams](#register-teams)
- [Build Terraform Docker Image](build-terraform-docker-image)

# Overview

zLifecycle is a product to manage lifecycle for infrastructure across various cloud providers as well as on-prem.

For more details & diagrams look at: https://app.diagrams.net/#G1gXeFRlERpqjXpeSjxRPLP6YZMRyFG5SN

## Runbooks

#### What goes in a runbook vs a README?
If a piece of documentation involves step-by-step procedures, executing commands, or directly references code (e.g. use of a variable defined in code), consider creating a runbook for it in the `runbook` directory. See the [base template](./runbook/template.md) as a guide to get started.

## Creating a new environment
1. Configure bootstrap scripts with new terraform workspace, etc. based on environment name
1. Create new customers in that environment, see (onboarding a new customer)[#onboarding-a-new-customer]

## Onboarding a new customer
1. Create a new config repo for that customer
1. Create a new empty IL repo for that customer
1. Create a new directory in the `zlifecycle-il-operator` service for the new customer

## Register Teams

You need to manually register teams currently using following script

```bash
cd ../../compuzest-zlifecycle-config
kubectl apply -f teams/account-team.yaml # Replace yaml file with team name for the team you want to register
```

## Build Terraform Docker Image

```bash
cd terraform-image
docker build -t 413422438110.dkr.ecr.us-east-1.amazonaws.com/zlifecycle-terraform:latest .
docker push 413422438110.dkr.ecr.us-east-1.amazonaws.com/zlifecycle-terraform:latest
```

## Connect to zLifecycle AWS environment
- Configure proper credentials in `~/.aws/credentials` as a profile names `compuzest-shared`
```
aws eks --region us-east-1 update-kubeconfig --name 0-sandbox-eks
```
