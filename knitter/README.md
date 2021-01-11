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
* Add new github service account to the customer github org and give perms to following repos
    * compuzest-zlifecycle-il - write access
    * helm-charts - read access
    * compuZest-zlifecycle-config - read access
* Generate Personal Token & ssh key for the Github service account to be used by secret created (Check LastPass secret note: "zLifecycle - k8s secrets")

## Initial Bootstrap 

### Terraform Shared State

zLifecycle environments (e.g. demo, dev) are managed by terraform workspaces. These terraform workspaces share a parent state directory maintained in terraform (`zlifecycle-tfstate`) that needs to be initialized before environments can be created. This bootstrap script is for this use case, where no zlifecycle environments exist yet.

Run `tfstate` terraform to provision S3 bucket and Dynamo DB table that will be used for Terraform Shared State.

```bash
cd tfstate
terraform init
terraform apply
```

## Bootstrap zLifecycle

To bootstrap zLifecycle in a given environment (e.g. demo, dev-a, dev-b):
1. Add the zlifecycle GitHub service account SSH key pair to `zlifecycle-provisioner/k8s-addons/argo-workflow` in files named `zlifecycle` and `zlifecycle.pub`.
2. Create a `tfvars` file for your environment in `zlifecycle-provisioner/k8s-addons/tfvars` based on the example file. Non `.example` files will be git ignored. Add required values, such as the ArgoCD slack token.
3. Run following script and following instructions with the following note:

Note: When it asks to create secret go to `zlifecycle-provisioner/k8s-addons/argo-workflow` folder
and create secrets using scripts in LastPass. This will ensure the GitHub key created in step 1 is used.

```bash
cd zlifecycle/bootstrap
./bootstrap_zLifecycle.sh
```

## Register Teams

You need to manually register teams currently using following script

```bash
cd ../../compuzest-zlifecycle-config
kubectl apply -f teams/account-team.yaml # Replace yaml file with team name for the team you want to register
```

## Build Terraform Docker Image

```bash
cd terraform
docker build -t shahadarsh/terraform:latest .
docker push shahadarsh/terraform:latest
```

## Connect to zLifecycle AWS environment
- Configure proper credentials in `~/.aws/credentials` as a profile names `compuzest-shared`
```
aws eks --region us-east-1 update-kubeconfig --name 0-sandbox-eks
```