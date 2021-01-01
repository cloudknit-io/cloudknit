# Index
- [Overview](#overview)
- [Initial Bootstrap](#initial-bootstrap)
- [Bootstrap zLifecycle](#bootstrap-zlifecycle)
- [Register Teams](#register-teams)

# Overview
Repo with everything needed to Bootstrap zLifecycle Platform

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
cd scripts
./bootstrap_zLifecycle.sh
```

## Register Teams

You need to manually register teams currently using following script

```bash
cd ../../zLifecycle-teams
kubectl apply -f teams/account-team.yaml # Replace yaml file with team name for the team you want to register
```
