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

1. Manually trigger github action from the UI [here](https://github.com/CompuZest/zlifecycle/actions/workflows/main.yaml), passing `image_tag` value that indicates the feature you'd like to test, e.g. `update-argo`
2. When the image is built, test it out in a demo or sandbox kubernetes environment.
- First pull the latest `main` from your local clone of this repo and then locally update the `terraform-sync-template.yaml` file to point to your new image by changing the tag [here](https://github.com/CompuZest/zlifecycle/blob/4c536ee9223434d996449e1aa53345332d1a0ef9/argo-templates/terraform-sync-template.yaml#L117).
- Then, update your local kubernetes context to the cluster you'd like to text with (using `kubectx`) and run the following from this repo's root folder.
`kubectl apply -f ./argo-templates/terraform-sync-template.yaml`
3. Once tested, re-trigger the github action from step one but this time with `image_tag` of `latest`
4. Reset the file from step 2 to point to the `latest` tag (or just `git checkout` your local changes) and re-apply with step 3


## Connect to zLifecycle AWS environment
- Configure proper credentials in `~/.aws/credentials` as a profile names `compuzest-shared`
```
aws eks --region us-east-1 update-kubeconfig --name sandbox-eks
```

## SSH key overview

1. A customer provide SSH _read_ access to their company config repo and team _config repos_.
- For now there is one key per customer and it is `${company-name}-ssh` stored in the `argocd` namespace
2. Zlifeycle operator needs _write_ access to zlifecycle il repos
- This kubernetes secret storing the key is called `zlifecycle-operator-ssh` and is stored in the operator's namespace `zlifecycle-il-operator-system`
- The kubernetes secret is parsed by operator code and the key is `sshPrivateKey`
3. Zlifeycle provisioners (Argo workflows and go provisioner service) need _read_ access to zlifecycle il repos for all customers and to `helm-charts`
- This key is called `zlifecycle-provisioner-ssh` and is stored in the operator's namespace `argocd`
- The kubernetes secret is used to mount files so it stores keys `id_rsa`