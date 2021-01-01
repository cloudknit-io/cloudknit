# Index
- [Overview](#overview)
- [Initial Bootstrap](#initial-bootstrap)
- [Bootstrap zLifecycle](#bootstrap-zlifecycle)

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
