# Bootstrap zLifecycle 

## Overview

Steps to Bootstrap zLifecycle from scratch

## When to use this runbook
When you want to bootstrap  zLifecycle on AWS or locally

## Initial Steps Overview

1. [Bootstrap zLifecycle](#bootstrap-zlifecycle)
1. [Bootstrap Gotchas](#bootstrap-gotchas)

## Detailed Steps

#### Bootstrap zLifecycle

To bootstrap zLifecycle in a given environment (e.g. demo, dev-a, dev-b):
1. Download the zlifecycle GitHub service account SSH key pair (from LastPass) to `zlifecycle-provisioner/k8s-addons/argo-workflow` folder on your machine and name the files `zlifecycle` and `zlifecycle.pub`. 
If you already have those files locally no need to do it again unless the key pair changed.
2. Create a `tfvars` file for your environment in `zlifecycle-provisioner/k8s-addons/tfvars` based on the example file. Non `.example` files will be git ignored. Add required values, such as the ArgoCD slack token.
3. Run the bootstrap script:

```bash
cd zlifecycle/bootstrap
./bootstrap_zLifecycle.sh
```
4. When the script stops to ask if secretes have been created, go to `zlifecycle-provisioner/k8s-addons/argo-workflow` folder
and create secrets using scripts in LastPass. This will ensure the GitHub key created in step 1 is used. Then enter `Y` and allow bootstrap to continue


#### Bootstrap Gotchas
1. Make sure you have run `brew bundle` in the `company` repo to ensure you have all the dependencies in `osx/Brewfile`. Missing dependencies, for example `aws-iam-authenticator` can cause weird errors
1. `aws-eks` failing on `null_resource.wait_for_cluster`, try `terraform destroy`ing that resource and re-applying
